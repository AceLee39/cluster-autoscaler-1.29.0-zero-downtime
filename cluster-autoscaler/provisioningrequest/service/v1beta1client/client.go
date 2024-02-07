/*
Copyright 2023 The Kubernetes Authors.

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

package v1beta1client

import (
	"fmt"
	"time"

	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"

	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/autoscaler/cluster-autoscaler/provisioningrequest/apis/autoscaling.x-k8s.io/v1beta1"
	"k8s.io/autoscaler/cluster-autoscaler/provisioningrequest/client/clientset/versioned"
	"k8s.io/autoscaler/cluster-autoscaler/provisioningrequest/client/informers/externalversions"
	listers "k8s.io/autoscaler/cluster-autoscaler/provisioningrequest/client/listers/autoscaling.x-k8s.io/v1beta1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/rest"

	klog "k8s.io/klog/v2"
)

const (
	provisioningRequestClientCallTimeout = 4 * time.Second
)

// ProvisioningRequestClient represents client for v1beta1 ProvReq CRD.
type ProvisioningRequestClient struct {
	client         versioned.Interface
	provReqLister  listers.ProvisioningRequestLister
	podTemplLister v1.PodTemplateLister
}

// NewProvisioningRequestClient configures and returns a provisioningRequestClient.
func NewProvisioningRequestClient(kubeConfig *rest.Config) (*ProvisioningRequestClient, error) {
	prClient, err := newPRClient(kubeConfig)
	if err != nil {
		return nil, fmt.Errorf("Failed to create Provisioning Request client: %v", err)
	}

	provReqLister, err := newPRsLister(prClient, make(chan struct{}))
	if err != nil {
		return nil, err
	}

	podTemplateClient, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		return nil, fmt.Errorf("Failed to create Pod Template client: %v", err)
	}

	podTemplLister, err := newPodTemplatesLister(podTemplateClient, make(chan struct{}))
	if err != nil {
		return nil, err
	}

	return &ProvisioningRequestClient{
		client:         prClient,
		provReqLister:  provReqLister,
		podTemplLister: podTemplLister,
	}, nil
}

// ProvisioningRequest gets a specific ProvisioningRequest CR.
func (c *ProvisioningRequestClient) ProvisioningRequest(namespace, name string) (*v1beta1.ProvisioningRequest, error) {
	return c.provReqLister.ProvisioningRequests(namespace).Get(name)
}

// ProvisioningRequests gets all ProvisioningRequest CRs.
func (c *ProvisioningRequestClient) ProvisioningRequests() ([]*v1beta1.ProvisioningRequest, error) {
	provisioningRequests, err := c.provReqLister.List(labels.Everything())
	if err != nil {
		return nil, fmt.Errorf("error fetching provisioningRequests: %w", err)
	}
	return provisioningRequests, nil
}

// FetchPodTemplates fetches PodTemplates referenced by the Provisioning Request.
func (c *ProvisioningRequestClient) FetchPodTemplates(pr *v1beta1.ProvisioningRequest) ([]*apiv1.PodTemplate, error) {
	podTemplates := make([]*apiv1.PodTemplate, 0, len(pr.Spec.PodSets))
	for _, podSpec := range pr.Spec.PodSets {
		podTemplate, err := c.podTemplLister.PodTemplates(pr.Namespace).Get(podSpec.PodTemplateRef.Name)
		if errors.IsNotFound(err) {
			klog.Infof("While fetching Pod Template for Provisioning Request %s/%s received not found error", pr.Namespace, pr.Name)
			continue
		} else if err != nil {
			return nil, err
		}
		podTemplates = append(podTemplates, podTemplate)
	}
	return podTemplates, nil
}

// newPRClient creates a new Provisioning Request client from the given config.
func newPRClient(kubeConfig *rest.Config) (*versioned.Clientset, error) {
	return versioned.NewForConfig(kubeConfig)
}

// newPRsLister creates a lister for the Provisioning Requests in the cluster.
func newPRsLister(prClient versioned.Interface, stopChannel <-chan struct{}) (listers.ProvisioningRequestLister, error) {
	factory := externalversions.NewSharedInformerFactory(prClient, 1*time.Hour)
	provReqLister := factory.Autoscaling().V1beta1().ProvisioningRequests().Lister()
	factory.Start(stopChannel)
	informersSynced := factory.WaitForCacheSync(stopChannel)
	for _, synced := range informersSynced {
		if !synced {
			return nil, fmt.Errorf("can't create Provisioning Request lister")
		}
	}
	klog.V(2).Info("Successful initial Provisioning Request sync")
	return provReqLister, nil
}

// newPodTemplatesLister creates a lister for the Pod Templates in the cluster.
func newPodTemplatesLister(client *kubernetes.Clientset, stopChannel <-chan struct{}) (v1.PodTemplateLister, error) {
	factory := informers.NewSharedInformerFactory(client, 1*time.Hour)
	podTemplLister := factory.Core().V1().PodTemplates().Lister()
	factory.Start(stopChannel)
	informersSynced := factory.WaitForCacheSync(stopChannel)
	for _, synced := range informersSynced {
		if !synced {
			return nil, fmt.Errorf("can't create Pod Template lister")
		}
	}
	klog.V(2).Info("Successful initial Pod Template sync")
	return podTemplLister, nil
}
