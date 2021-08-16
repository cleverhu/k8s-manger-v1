package deployment

import (
	"k8s-manger-v1/core"
	"k8s-manger-v1/lib"
	v1 "k8s.io/api/apps/v1"
)

func GetPodsByDep(dep v1.Deployment) []*Pod {
	rsLabelsMap, err := core.RSMap.GetRsLabelsByDeployment(&dep)
	lib.CheckError(err)
	pods, err := core.PodMap.ListByRsLabelsAndNS(dep.Namespace, rsLabelsMap)
	lib.CheckError(err)
	ret := make([]*Pod, 0)
	for _, pod := range pods {
		ret = append(ret, &Pod{
			Name:       pod.Name,
			NameSpace:  pod.Namespace,
			Images:     GetImagesByPod(pod.Spec.Containers),
			NodeName:   pod.Spec.NodeName,
			CreateTime: TimeFormat(pod.CreationTimestamp.Time),
			IPs:        []string{pod.Status.PodIP, pod.Status.HostIP},
		})
	}
	return ret
}

func Detail(namespace string, name string) *Deployment {

	//ctx := context.Background()
	//getOpts := metav1.GetOptions{}
	//deploy, err := lib.K8sClient.AppsV1().Deployments(namespace).Get(ctx, name, getOpts)
	deploy, err := core.DeploymentMap.Get(namespace, name)
	lib.CheckError(err)
	return &Deployment{
		Name:       name,
		NameSpace:  namespace,
		Images:     GetImages(*deploy),
		CreateTime: TimeFormat(deploy.CreationTimestamp.Time),
		Pods:       GetPodsByDep(*deploy),
		Replicas:   [3]int32{deploy.Status.Replicas, deploy.Status.AvailableReplicas, deploy.Status.UnavailableReplicas},
	}
}
