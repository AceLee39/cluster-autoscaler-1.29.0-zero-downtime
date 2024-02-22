/*
Copyright 2016 The Kubernetes Authors.

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
//customized
package kubernetes

import (
	"golang.org/x/net/context"
	apiv1 "k8s.io/api/core/v1"
	kube_errors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	client "k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
	"time"
)

// restart builds a configmap lister for the passed namespace (including all).
func Restart(kubeClient client.Interface, namespace string, name string) {
	klog.V(4).Infof("customized going to restart : %v %v", namespace, name)
	deploy, err := kubeClient.AppsV1().Deployments(namespace).Get(context.TODO(), name, v1.GetOptions{})
	if err != nil {
		klog.Errorf("customized get deploy by namespace, name: %v", err)
		return
	}
	if deploy.Spec.Template.ObjectMeta.Annotations == nil {
		deploy.Spec.Template.ObjectMeta.Annotations = make(map[string]string)
	}
	deploy.Spec.Template.ObjectMeta.Annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().Format(time.RFC3339)
	_, err = kubeClient.AppsV1().Deployments(namespace).Update(context.TODO(), deploy, v1.UpdateOptions{})
	if err != nil {
		klog.Errorf("customized deploy namespace=%s, name=%s: %v", namespace, name, err)
		return
	}
	return
}

func WaitPodsToDisappear(kubeClient client.Interface, node *apiv1.Node, pods []*apiv1.Pod) error {
	var allGone bool
	maxTermination := 300
	for start := time.Now(); time.Now().Sub(start) < time.Duration(maxTermination)*time.Second; time.Sleep(5 * time.Second) {
		allGone = true
		for _, pod := range pods {
			klog.V(1).Infof("customized waiting pod disappear %s/%s", pod.Namespace, pod.Name)
			podReturned, err := kubeClient.CoreV1().Pods(pod.Namespace).Get(context.TODO(), pod.Name, v1.GetOptions{})
			if err == nil && (podReturned == nil || podReturned.Spec.NodeName == node.Name) {
				klog.V(1).Infof("customized pod node.Name didn't update %s/%s", pod.Namespace, pod.Name)
				allGone = false
				//break
			}
			if err != nil && !kube_errors.IsNotFound(err) {
				klog.Errorf("customized Failed to check pod %s/%s: %v", pod.Namespace, pod.Name, err)
				allGone = false
				//break
			}
		}
		if allGone {
			return nil
		}
	}

	for _, pod := range pods {
		podReturned, err := kubeClient.CoreV1().Pods(pod.Namespace).Get(context.TODO(), pod.Name, v1.GetOptions{})
		if err == nil && (podReturned == nil || podReturned.Spec.NodeName == node.Name) {
			klog.Errorf("customized rollout restart pod nodeName didn't update %s/%s %s: %v", pod.Namespace, pod.Name, podReturned.Spec.NodeName, err)
		} else if err != nil && !kube_errors.IsNotFound(err) {
			klog.Errorf("customized Failed to rollout restart pod %s/%s: %v", pod.Namespace, pod.Name, err)
		} else {
			klog.V(1).Infof("customized Successfully to rollout restart pod %s/%s: %v", pod.Namespace, pod.Name, err)
		}
	}

	return nil
}
