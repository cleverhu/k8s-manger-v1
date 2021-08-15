package main

import (
	"context"
	"fmt"
	"k8s-manger-v1/lib"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	deploy, _ := lib.K8sClient.AppsV1().Deployments("default").Get(context.Background(), "ngx1", metav1.GetOptions{})
	selector, _ := metav1.LabelSelectorAsSelector(deploy.Spec.Selector)
	rs, _ := lib.K8sClient.AppsV1().ReplicaSets("default").List(context.Background(), metav1.ListOptions{LabelSelector: selector.String()})
	labelName := ""
	for _, item := range rs.Items {
		if item.Annotations["deployment.kubernetes.io/revision"] != deploy.Annotations["deployment.kubernetes.io/revision"] {
			continue
		}
		for _, v := range item.OwnerReferences {
			if v.Name == deploy.Name {
				s, _ := metav1.LabelSelectorAsSelector(item.Spec.Selector)
				labelName = s.String()
			}
		}
	}
	pods, _ := lib.K8sClient.CoreV1().Pods("default").List(context.Background(), metav1.ListOptions{LabelSelector: labelName})
	for _, pod := range pods.Items {
		fmt.Println(pod.Name)
	}
}
