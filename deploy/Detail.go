package deploy

import (
	"context"
	"k8s-manger-v1/lib"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Detail(namespace string, name string) *Deployment {

	ctx := context.Background()
	getOpts := metav1.GetOptions{}
	deploy, _ := lib.K8sClient.AppsV1().Deployments(namespace).Get(ctx, name, getOpts)
	return &Deployment{
		Name:       name,
		NameSpace:  namespace,
		Images:     GetImages(*deploy),
		CreateTime: deploy.CreationTimestamp.Format("2006-01-02 15:04:05"),
	}
}
