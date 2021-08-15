package deploy

import (
	"k8s-manger-v1/core"
)

func ListAll(namespace string) []*Deployment {
	//if namespace == "" {
	//	namespace = "default"
	//}
	ret := make([]*Deployment, 0)

	deps, _ := core.DepMap.ListByNS(namespace)

	for _, dep := range deps {
		tmp := &Deployment{
			NameSpace: dep.Namespace,
			Name:      dep.Name,
			Replicas:  [3]int32{dep.Status.Replicas, dep.Status.AvailableReplicas, dep.Status.UnavailableReplicas},
			Images:    GetImages(*dep),
		}
		ret = append(ret, tmp)
	}

	return ret
}
