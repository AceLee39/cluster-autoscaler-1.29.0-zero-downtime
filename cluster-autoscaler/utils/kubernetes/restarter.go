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

package kubernetes

import (
	"golang.org/x/net/context"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	client "k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
	"time"
)

// restart builds a configmap lister for the passed namespace (including all).
func Restart(kubeClient client.Interface,
	namespace string, name string) {
	klog.V(4).Info("customized going to restart : %v %v", namespace, name)
	deploy, err := kubeClient.AppsV1().Deployments(namespace).Get(context.TODO(), name, v1.GetOptions{})
	if err != nil {
		klog.Errorf("customized get deploy by namespace, name: %v", err)
	}
	if deploy.Spec.Template.ObjectMeta.Annotations == nil {
		deploy.Spec.Template.ObjectMeta.Annotations = make(map[string]string)
	}
	deploy.Spec.Template.ObjectMeta.Annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().Format(time.RFC3339)
	_, err = kubeClient.AppsV1().Deployments(namespace).Update(context.TODO(), deploy, v1.UpdateOptions{})
	if err != nil {
		klog.Errorf("customized deploy namespace, name: %v", err)
	}

	//data := fmt.Sprintf(`{"spec": {"template": {"metadata": {"annotations": {"kubectl.kubernetes.io/restartedAt": "%s"}}}}}`, time.Now().Format("20060102150405"))
	//deploymentsClient := kubeClient.AppsV1().Deployments(namespace)
	//deployment, err := deploymentsClient.Patch(context.TODO(), deployment_name, k8stypes.StrategicMergePatchType, []byte(data), v1.PatchOptions{})
	//if err != nil {
	//	klog.Errorf("customized complete deploymentName: %v", err)
	//}

	//
	//deploymentName := "deployment/abc"
	//streams, _, _, _ := genericiooptions.NewTestIOStreams()
	//cmdutil.NewMatchVersionFlags()
	//cmdutil.NewFactory()
	//r := &rollout.RestartOptions{
	//	PrintFlags: genericclioptions.NewPrintFlags("restarted").WithTypeSetter(scheme.Scheme),
	//	Resources:  []string{deploymentName},
	//	IOStreams:  streams,
	//}
	//err := r.Complete(f, nil, []string{deploymentName})
	//if err != nil {
	//	klog.Errorf("customized complete deploymentName: %v", err)
	//}
	//err = r.RunRestart()
	//if err != nil {
	//	klog.Errorf("customized RunRestart deploymentName: %v", err)
	//}
	return
}
