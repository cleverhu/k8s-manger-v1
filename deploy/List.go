package deploy

import (
	"context"
	"k8s-manger-v1/lib"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ListAll(namespace string) []*Deployment {
	deps := []*Deployment{}
	ctx := context.Background()
	opt := metav1.ListOptions{}
	list, _ := lib.K8sClient.AppsV1().Deployments(namespace).List(ctx, opt)
	for _, item := range list.Items {
		dep := &Deployment{Name: item.Name}
		deps = append(deps, dep)
	}
	return deps
}
