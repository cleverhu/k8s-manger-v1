package core

import (
	"errors"
	"fmt"
	"k8s-manger-v1/lib"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"log"
	"sort"
	"sync"
)

type DeploymentMapStruct struct {
	Data sync.Map //key ns value []*v1.Deployment
}

func (this *DeploymentMapStruct) Add(deploy *v1.Deployment) {
	key := deploy.Namespace
	if value, ok := this.Data.Load(key); ok {
		value = append(value.([]*v1.Deployment), deploy)
		this.Data.Store(key, value)
	} else {
		this.Data.Store(key, []*v1.Deployment{deploy})
	}
}

func (this *DeploymentMapStruct) Delete(deploy *v1.Deployment) {
	key := deploy.Namespace
	if value, ok := this.Data.Load(key); ok {
		for index, dep := range value.([]*v1.Deployment) {
			if dep.Name == deploy.Name {
				value = append(value.([]*v1.Deployment)[0:index], value.([]*v1.Deployment)[index+1:]...)
				this.Data.Store(key, value)
				return
			}
		}
	}
}

func (this *DeploymentMapStruct) Update(deploy *v1.Deployment) error {
	key := deploy.Namespace
	if value, ok := this.Data.Load(key); ok {
		for index, dep := range value.([]*v1.Deployment) {
			if dep.Name == deploy.Name {
				value.([]*v1.Deployment)[index] = deploy
				this.Data.Store(key, value)
				return nil
			}
		}
	}

	return fmt.Errorf("deployment-%s not found", deploy.Name)
}

func (this *DeploymentMapStruct) ListByNS(ns string) ([]*v1.Deployment, error) {
	if ns != "" {
		if list, ok := this.Data.Load(ns); ok {
			return list.([]*v1.Deployment), nil
		}
	} else {
		ret := make([]*v1.Deployment, 0)
		DeploymentMap.Data.Range(func(key, value interface{}) bool {
			for _, dep := range value.([]*v1.Deployment) {
				ret = append(ret, dep)
			}
			sort.Slice(ret, func(i, j int) bool {
				if ret[i].Namespace == ret[j].Namespace {
					return ret[i].Name < ret[j].Name
				} else {
					return ret[i].Namespace < ret[j].Namespace
				}
			})
			return true
		})
		return ret, nil
	}

	return nil, errors.New("deployments record not found")
}

func (this *DeploymentMapStruct) Get(ns, name string) (*v1.Deployment, error) {
	deps, err := this.ListByNS(ns)
	if err != nil {
		return nil, err
	}
	for _, dep := range deps {
		if dep.Name == name {
			return dep, nil
		}
	}
	return nil, errors.New("deployment record not found")
}

type DepHandler struct {
}

func (this *DepHandler) OnAdd(obj interface{}) {
	dep := obj.(*v1.Deployment)
	DeploymentMap.Add(dep)
}
func (this *DepHandler) OnUpdate(oldObj interface{}, newObj interface{}) {
	err := DeploymentMap.Update(newObj.(*v1.Deployment))
	if err != nil {
		log.Println(err)
	}
}
func (this *DepHandler) OnDelete(obj interface{}) {
	DeploymentMap.Delete(obj.(*v1.Deployment))
}

var DeploymentMap *DeploymentMapStruct

func init() {
	DeploymentMap = &DeploymentMapStruct{}

}

func InitDeployment() {
	factory := informers.NewSharedInformerFactory(lib.K8sClient, 0)
	depInformer := factory.Apps().V1().Deployments().Informer()
	depInformer.AddEventHandler(&DepHandler{})

	podInformer := factory.Core().V1().Pods().Informer()
	podInformer.AddEventHandler(&PodHandler{})

	rsInformer := factory.Apps().V1().ReplicaSets().Informer()
	rsInformer.AddEventHandler(&RSHandler{})

	factory.Start(wait.NeverStop)
}
