package deploy

import (
	"context"
	"k8s-manger-v1/lib"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

func GetPodsByDep(namespace string, dep v1.Deployment) []*Pod {
	ctx := context.Background()
	listOpts := metav1.ListOptions{
		LabelSelector: GetLabels(dep.Spec.Selector.MatchLabels),
	}
	pods, err := lib.K8sClient.CoreV1().Pods(namespace).List(ctx, listOpts)
	lib.CheckError(err)
	ret := make([]*Pod, 0)
	for _, pod := range pods.Items {
		//fmt.Println(pod.Name, dep.Name)
		if strings.HasPrefix(pod.Name, dep.Name) {
			//fmt.Println(pod.Status.ContainerStatuses)
			ret = append(ret, &Pod{
				Name:       pod.Name,
				Images:     GetImagesByPod(pod.Spec.Containers),
				NodeName:   pod.Spec.NodeName,
				CreateTime: TimeFormat(pod.CreationTimestamp.Time),
				IP:         pod.Status.PodIP,
			})
		}
	}
	return ret
}

func Detail(namespace string, name string) *Deployment {

	ctx := context.Background()
	getOpts := metav1.GetOptions{}
	deploy, err := lib.K8sClient.AppsV1().Deployments(namespace).Get(ctx, name, getOpts)
	lib.CheckError(err)
	return &Deployment{
		Name:       name,
		NameSpace:  namespace,
		Images:     GetImages(*deploy),
		CreateTime: TimeFormat(deploy.CreationTimestamp.Time),
		Pods:       GetPodsByDep(namespace, *deploy),
		Replicas:   [3]int32{deploy.Status.Replicas, deploy.Status.AvailableReplicas, deploy.Status.UnavailableReplicas},
	}
}
