package deploy

import (
	"context"
	"k8s-manger-v1/lib"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ListAll(namespace string) []*Deployment {

	if namespace == "" {
		namespace = "default"
	}
	deps := make([]*Deployment, 0)
	ctx := context.Background()
	listOpts := metav1.ListOptions{}
	list, _ := lib.K8sClient.AppsV1().Deployments(namespace).List(ctx, listOpts)
	for _, item := range list.Items {

		dep := &Deployment{
			NameSpace: namespace,
			Name:      item.Name,
			Replicas:  [3]int32{item.Status.Replicas, item.Status.AvailableReplicas, item.Status.UnavailableReplicas},
			Images:    GetImages(item),
		}
		deps = append(deps, dep)
	}

	return deps
}
