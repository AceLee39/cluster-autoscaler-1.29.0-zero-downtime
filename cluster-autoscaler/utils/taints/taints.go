/*
Copyright 2020 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package taints

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/autoscaler/cluster-autoscaler/config"
	"k8s.io/autoscaler/cluster-autoscaler/utils/kubernetes"
	kube_client "k8s.io/client-go/kubernetes"
	kube_record "k8s.io/client-go/tools/record"
	cloudproviderapi "k8s.io/cloud-provider/api"

	klog "k8s.io/klog/v2"
)

const (
	// ToBeDeletedTaint is a taint used to make the node unschedulable.
	ToBeDeletedTaint = "ToBeDeletedByClusterAutoscaler"
	// DeletionCandidateTaint is a taint used to mark unneeded node as preferably unschedulable.
	DeletionCandidateTaint = "DeletionCandidateOfClusterAutoscaler"

	// IgnoreTaintPrefix any taint starting with it will be filtered out from autoscaler template node.
	IgnoreTaintPrefix = "ignore-taint.cluster-autoscaler.kubernetes.io/"

	// StartupTaintPrefix (Same as IgnoreTaintPrefix) any taint starting with it will be filtered out from autoscaler template node.
	StartupTaintPrefix = "startup-taint.cluster-autoscaler.kubernetes.io/"

	// StatusTaintPrefix any taint starting with it will be filtered out from autoscaler template node but unlike IgnoreTaintPrefix & StartupTaintPrefix it should not be trated as unready.
	StatusTaintPrefix = "status-taint.cluster-autoscaler.kubernetes.io/"

	gkeNodeTerminationHandlerTaint = "cloud.google.com/impending-node-termination"

	// AWS: Indicates that a node has volumes stuck in attaching state and hence it is not fit for scheduling more pods
	awsNodeWithImpairedVolumesTaint = "NodeWithImpairedVolumes"

	// statusNodeTaintReportedType is the value used when reporting node taint count defined as status taint in given taintConfig.
	statusNodeTaintReportedType = "status-taint"

	// startupNodeTaintReportedType is the value used when reporting node taint count defined as startup taint in given taintConfig.
	startupNodeTaintReportedType = "startup-taint"

	// unlistedNodeTaintReportedType is the value used when reporting node taint count in case taint key is other than defined in explicitlyReportedNodeTaints and taintConfig.
	unlistedNodeTaintReportedType = "other"
)

var (
	// NodeConditionTaints lists taint keys used as node conditions
	NodeConditionTaints = TaintKeySet{
		apiv1.TaintNodeNotReady:                     true,
		apiv1.TaintNodeUnreachable:                  true,
		apiv1.TaintNodeUnschedulable:                true,
		apiv1.TaintNodeMemoryPressure:               true,
		apiv1.TaintNodeDiskPressure:                 true,
		apiv1.TaintNodeNetworkUnavailable:           true,
		apiv1.TaintNodePIDPressure:                  true,
		cloudproviderapi.TaintExternalCloudProvider: true,
		cloudproviderapi.TaintNodeShutdown:          true,
		gkeNodeTerminationHandlerTaint:              true,
		awsNodeWithImpairedVolumesTaint:             true,
	}

	// Mutable only in unit tests
	maxRetryDeadline      time.Duration = 5 * time.Second
	conflictRetryInterval time.Duration = 750 * time.Millisecond
)

// TaintKeySet is a set of taint key
type TaintKeySet map[string]bool

// TaintConfig is a config of taints that require special handling
type TaintConfig struct {
	startupTaints            TaintKeySet
	statusTaints             TaintKeySet
	startupTaintPrefixes     []string
	statusTaintPrefixes      []string
	explicitlyReportedTaints TaintKeySet
}

// NewTaintConfig returns the taint config extracted from options
func NewTaintConfig(opts config.AutoscalingOptions) TaintConfig {
	startupTaints := make(TaintKeySet)
	for _, taintKey := range opts.StartupTaints {
		klog.V(4).Infof("Startup taint %s on all NodeGroups", taintKey)
		startupTaints[taintKey] = true
	}

	statusTaints := make(TaintKeySet)
	for _, taintKey := range opts.StatusTaints {
		klog.V(4).Infof("Status taint %s on all NodeGroups", taintKey)
		statusTaints[taintKey] = true
	}

	explicitlyReportedTaints := TaintKeySet{
		ToBeDeletedTaint:       true,
		DeletionCandidateTaint: true,
	}

	for k, v := range NodeConditionTaints {
		explicitlyReportedTaints[k] = v
	}

	return TaintConfig{
		startupTaints:            startupTaints,
		statusTaints:             statusTaints,
		startupTaintPrefixes:     []string{IgnoreTaintPrefix, StartupTaintPrefix},
		statusTaintPrefixes:      []string{StatusTaintPrefix},
		explicitlyReportedTaints: explicitlyReportedTaints,
	}
}

// IsStartupTaint checks whether given taint is a startup taint.
func (tc TaintConfig) IsStartupTaint(taint string) bool {
	if _, ok := tc.startupTaints[taint]; ok {
		return true
	}
	return matchesAnyPrefix(tc.startupTaintPrefixes, taint)
}

// IsStatusTaint checks whether given taint is a status taint.
func (tc TaintConfig) IsStatusTaint(taint string) bool {
	if _, ok := tc.statusTaints[taint]; ok {
		return true
	}
	return matchesAnyPrefix(tc.statusTaintPrefixes, taint)
}

func (tc TaintConfig) isExplicitlyReportedTaint(taint string) bool {
	_, ok := tc.explicitlyReportedTaints[taint]
	return ok
}

// getKeyShortName converts taint key to short name for logging
func getKeyShortName(key string) string {
	switch key {
	case ToBeDeletedTaint:
		return "ToBeDeletedTaint"
	case DeletionCandidateTaint:
		return "DeletionCandidateTaint"
	default:
		return key
	}
}

// MarkToBeDeleted sets a taint that makes the node unschedulable.
func MarkToBeDeleted(node *apiv1.Node, client kube_client.Interface, cordonNode bool) error {
	taint := apiv1.Taint{
		Key:    ToBeDeletedTaint,
		Value:  fmt.Sprint(time.Now().Unix()),
		Effect: apiv1.TaintEffectNoSchedule,
	}
	return AddTaint(node, client, taint, cordonNode)
}

// MarkDeletionCandidate sets a soft taint that makes the node preferably unschedulable.
func MarkDeletionCandidate(node *apiv1.Node, client kube_client.Interface) error {
	taint := apiv1.Taint{
		Key:    DeletionCandidateTaint,
		Value:  fmt.Sprint(time.Now().Unix()),
		Effect: apiv1.TaintEffectPreferNoSchedule,
	}
	return AddTaint(node, client, taint, false)
}

// AddTaint sets the specified taint on the node.
func AddTaint(node *apiv1.Node, client kube_client.Interface, taint apiv1.Taint, cordonNode bool) error {
	retryDeadline := time.Now().Add(maxRetryDeadline)
	freshNode := node.DeepCopy()
	var err error
	refresh := false
	for {
		if refresh {
			// Get the newest version of the node.
			freshNode, err = client.CoreV1().Nodes().Get(context.TODO(), node.Name, metav1.GetOptions{})
			if err != nil || freshNode == nil {
				klog.Warningf("Error while adding %v taint on node %v: %v", getKeyShortName(taint.Key), node.Name, err)
				return fmt.Errorf("failed to get node %v: %v", node.Name, err)
			}
		}

		if !addTaintToSpec(freshNode, taint, cordonNode) {
			if !refresh {
				// Make sure we have the latest version before skipping update.
				refresh = true
				continue
			}
			return nil
		}
		_, err = client.CoreV1().Nodes().Update(context.TODO(), freshNode, metav1.UpdateOptions{})
		if err != nil && errors.IsConflict(err) && time.Now().Before(retryDeadline) {
			refresh = true
			time.Sleep(conflictRetryInterval)
			continue
		}

		if err != nil {
			klog.Warningf("Error while adding %v taint on node %v: %v", getKeyShortName(taint.Key), node.Name, err)
			return err
		}
		klog.V(1).Infof("Successfully added %v on node %v", getKeyShortName(taint.Key), node.Name)
		return nil
	}
}

func addTaintToSpec(node *apiv1.Node, taint apiv1.Taint, cordonNode bool) bool {
	for _, t := range node.Spec.Taints {
		if t.Key == taint.Key {
			klog.V(2).Infof("%v already present on node %v, t: %v", taint.Key, node.Name, t)
			return false
		}
	}
	node.Spec.Taints = append(node.Spec.Taints, taint)
	if cordonNode {
		klog.V(1).Infof("Marking node %v to be cordoned by Cluster Autoscaler", node.Name)
		node.Spec.Unschedulable = true
	}
	return true
}

// HasToBeDeletedTaint returns true if ToBeDeleted taint is applied on the node.
func HasToBeDeletedTaint(node *apiv1.Node) bool {
	return HasTaint(node, ToBeDeletedTaint)
}

// HasDeletionCandidateTaint returns true if DeletionCandidate taint is applied on the node.
func HasDeletionCandidateTaint(node *apiv1.Node) bool {
	return HasTaint(node, DeletionCandidateTaint)
}

// HasTaint returns true if the specified taint is applied on the node.
func HasTaint(node *apiv1.Node, taintKey string) bool {
	for _, taint := range node.Spec.Taints {
		if taint.Key == taintKey {
			return true
		}
	}
	return false
}

// GetToBeDeletedTime returns the date when the node was marked by CA as for delete.
func GetToBeDeletedTime(node *apiv1.Node) (*time.Time, error) {
	return GetTaintTime(node, ToBeDeletedTaint)
}

// GetDeletionCandidateTime returns the date when the node was marked by CA as for delete.
func GetDeletionCandidateTime(node *apiv1.Node) (*time.Time, error) {
	return GetTaintTime(node, DeletionCandidateTaint)
}

// GetTaintTime returns the date when the node was marked by CA with the specified taint.
func GetTaintTime(node *apiv1.Node, taintKey string) (*time.Time, error) {
	for _, taint := range node.Spec.Taints {
		if taint.Key == taintKey {
			resultTimestamp, err := strconv.ParseInt(taint.Value, 10, 64)
			if err != nil {
				return nil, err
			}
			result := time.Unix(resultTimestamp, 0)
			return &result, nil
		}
	}
	return nil, nil
}

// CleanToBeDeleted cleans CA's NoSchedule taint from a node.
func CleanToBeDeleted(node *apiv1.Node, client kube_client.Interface, cordonNode bool) (bool, error) {
	return CleanTaint(node, client, ToBeDeletedTaint, cordonNode)
}

// CleanDeletionCandidate cleans CA's soft NoSchedule taint from a node.
func CleanDeletionCandidate(node *apiv1.Node, client kube_client.Interface) (bool, error) {
	return CleanTaint(node, client, DeletionCandidateTaint, false)
}

// CleanTaint cleans the specified taint from a node.
func CleanTaint(node *apiv1.Node, client kube_client.Interface, taintKey string, cordonNode bool) (bool, error) {
	retryDeadline := time.Now().Add(maxRetryDeadline)
	freshNode := node.DeepCopy()
	var err error
	refresh := false
	for {
		if refresh {
			// Get the newest version of the node.
			freshNode, err = client.CoreV1().Nodes().Get(context.TODO(), node.Name, metav1.GetOptions{})
			if err != nil || freshNode == nil {
				klog.Warningf("Error while adding %v taint on node %v: %v", getKeyShortName(taintKey), node.Name, err)
				return false, fmt.Errorf("failed to get node %v: %v", node.Name, err)
			}
		}
		newTaints := make([]apiv1.Taint, 0)
		for _, taint := range freshNode.Spec.Taints {
			if taint.Key == taintKey {
				klog.V(1).Infof("Releasing taint %+v on node %v", taint, node.Name)
			} else {
				newTaints = append(newTaints, taint)
			}
		}
		if len(newTaints) == len(freshNode.Spec.Taints) {
			if !refresh {
				// Make sure we have the latest version before skipping update.
				refresh = true
				continue
			}
			return false, nil
		}

		freshNode.Spec.Taints = newTaints
		if cordonNode {
			klog.V(1).Infof("Marking node %v to be uncordoned by Cluster Autoscaler", freshNode.Name)
			freshNode.Spec.Unschedulable = false
		}
		_, err = client.CoreV1().Nodes().Update(context.TODO(), freshNode, metav1.UpdateOptions{})

		if err != nil && errors.IsConflict(err) && time.Now().Before(retryDeadline) {
			refresh = true
			time.Sleep(conflictRetryInterval)
			continue
		}

		if err != nil {
			klog.Warningf("Error while releasing %v taint on node %v: %v", getKeyShortName(taintKey), node.Name, err)
			return false, err
		}
		klog.V(1).Infof("Successfully released %v on node %v", getKeyShortName(taintKey), node.Name)
		return true, nil
	}
}

// CleanAllToBeDeleted cleans ToBeDeleted taints from given nodes.
func CleanAllToBeDeleted(nodes []*apiv1.Node, client kube_client.Interface, recorder kube_record.EventRecorder, cordonNode bool) {
	CleanAllTaints(nodes, client, recorder, ToBeDeletedTaint, cordonNode)
}

// CleanAllDeletionCandidates cleans DeletionCandidate taints from given nodes.
func CleanAllDeletionCandidates(nodes []*apiv1.Node, client kube_client.Interface, recorder kube_record.EventRecorder) {
	CleanAllTaints(nodes, client, recorder, DeletionCandidateTaint, false)
}

// CleanAllTaints cleans all specified taints from given nodes.
func CleanAllTaints(nodes []*apiv1.Node, client kube_client.Interface, recorder kube_record.EventRecorder, taintKey string, cordonNode bool) {
	for _, node := range nodes {
		if !HasTaint(node, taintKey) {
			continue
		}
		cleaned, err := CleanTaint(node, client, taintKey, cordonNode)
		if err != nil {
			recorder.Eventf(node, apiv1.EventTypeWarning, "ClusterAutoscalerCleanup",
				"failed to clean %v on node %v: %v", getKeyShortName(taintKey), node.Name, err)
		} else if cleaned {
			recorder.Eventf(node, apiv1.EventTypeNormal, "ClusterAutoscalerCleanup",
				"removed %v taint from node %v", getKeyShortName(taintKey), node.Name)
		}
	}
}

func matchesAnyPrefix(prefixes []string, key string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(key, prefix) {
			return true
		}
	}
	return false
}

// SanitizeTaints returns filtered taints
func SanitizeTaints(taints []apiv1.Taint, taintConfig TaintConfig) []apiv1.Taint {
	var newTaints []apiv1.Taint
	for _, taint := range taints {
		switch taint.Key {
		case ToBeDeletedTaint:
			klog.V(4).Infof("Removing autoscaler taint when creating template from node")
			continue
		case DeletionCandidateTaint:
			klog.V(4).Infof("Removing autoscaler soft taint when creating template from node")
			continue
		}

		// ignore conditional taints as they represent a transient node state.
		if exists := NodeConditionTaints[taint.Key]; exists {
			klog.V(4).Infof("Removing node condition taint %s, when creating template from node", taint.Key)
			continue
		}

		if taintConfig.IsStartupTaint(taint.Key) || taintConfig.IsStatusTaint(taint.Key) {
			klog.V(4).Infof("Removing taint %s, when creating template from node", taint.Key)
			continue
		}

		newTaints = append(newTaints, taint)
	}
	return newTaints
}

// FilterOutNodesWithStartupTaints override the condition status of the given nodes to mark them as NotReady when they have
// filtered taints.
func FilterOutNodesWithStartupTaints(taintConfig TaintConfig, allNodes, readyNodes []*apiv1.Node) ([]*apiv1.Node, []*apiv1.Node) {
	newAllNodes := make([]*apiv1.Node, 0)
	newReadyNodes := make([]*apiv1.Node, 0)
	nodesWithStartupTaints := make(map[string]*apiv1.Node)
	for _, node := range readyNodes {
		if len(node.Spec.Taints) == 0 {
			newReadyNodes = append(newReadyNodes, node)
			continue
		}
		ready := true
		for _, t := range node.Spec.Taints {
			if taintConfig.IsStartupTaint(t.Key) {
				ready = false
				nodesWithStartupTaints[node.Name] = kubernetes.GetUnreadyNodeCopy(node, kubernetes.StartupNodes)
				klog.V(3).Infof("Overriding status of node %v, which seems to have startup taint %q", node.Name, t.Key)
				break
			}
		}
		if ready {
			newReadyNodes = append(newReadyNodes, node)
		}
	}
	// Override any node with ignored taint with its "unready" copy
	for _, node := range allNodes {
		if newNode, found := nodesWithStartupTaints[node.Name]; found {
			newAllNodes = append(newAllNodes, newNode)
		} else {
			newAllNodes = append(newAllNodes, node)
		}
	}
	return newAllNodes, newReadyNodes
}

// CountNodeTaints counts used node taints.
func CountNodeTaints(nodes []*apiv1.Node, taintConfig TaintConfig) map[string]int {
	foundTaintsCount := make(map[string]int)
	for _, node := range nodes {
		for _, taint := range node.Spec.Taints {
			key := getTaintTypeToReport(taint.Key, taintConfig)
			foundTaintsCount[key] += 1
		}
	}
	return foundTaintsCount
}

func getTaintTypeToReport(key string, taintConfig TaintConfig) string {
	// Track deprecated taints.
	if strings.HasPrefix(key, IgnoreTaintPrefix) {
		return IgnoreTaintPrefix
	}

	if taintConfig.isExplicitlyReportedTaint(key) {
		return key
	}
	if taintConfig.IsStartupTaint(key) {
		return startupNodeTaintReportedType
	}
	if taintConfig.IsStatusTaint(key) {
		return statusNodeTaintReportedType
	}
	return unlistedNodeTaintReportedType
}
