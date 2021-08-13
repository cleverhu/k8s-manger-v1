package deploy

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"k8s-manger-v1/lib"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

func Create() {
	ctx := context.Background()
	opt := metav1.ListOptions{}
	list, _ := lib.K8sClient.AppsV1().Deployments("default").List(ctx, opt)
	for _, item := range list.Items {
		fmt.Println(item.Name, item.Spec.Template.Spec.Containers[0].Image)
	}
	ngxDep := &v1.Deployment{}
	b, _ := ioutil.ReadFile("ymls/myngx.yml")
	ngxJson, _ := yaml.ToJSON(b)
	json.Unmarshal(ngxJson, ngxDep)
	createOpt := metav1.CreateOptions{}
	lib.K8sClient.AppsV1().Deployments("default").Create(ctx, ngxDep, createOpt)
}
