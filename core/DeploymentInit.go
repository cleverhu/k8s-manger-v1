package core

import (
	"errors"
	"fmt"
	"k8s-manger-v1/lib"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"log"
	"sync"
)

type DeploymentMap struct {
	Data sync.Map //key ns value []*v1.Deployment
}

func (this *DeploymentMap) Add(deploy *v1.Deployment) {
	key := deploy.Namespace
	if value, ok := this.Data.Load(key); ok {
		value = append(value.([]*v1.Deployment), deploy)
		this.Data.Store(key, value)
	} else {
		this.Data.Store(key, []*v1.Deployment{deploy})
	}
}

func (this *DeploymentMap) Delete(deploy *v1.Deployment) {
	key := deploy.Namespace
	if value, ok := this.Data.Load(key); ok {
		for index, dep := range value.([]*v1.Deployment) {
			if dep.Name == deploy.Name {
				value = append(value.([]*v1.Deployment)[0:index], value.([]*v1.Deployment)[index+1:]...)
				DepMap.Data.Store(key, value)
				return
			}
		}
	}
}

func (this *DeploymentMap) Update(deploy *v1.Deployment) error {
	key := deploy.Namespace
	if value, ok := this.Data.Load(key); ok {

		for index, dep := range value.([]*v1.Deployment) {
			if dep.Name == deploy.Name {
				value.([]*v1.Deployment)[index] = deploy
				return nil
			}
		}
	}

	return fmt.Errorf("deployment-%s not found", deploy.Name)
}

func (this *DeploymentMap) ListByNS(ns string) ([]*v1.Deployment, error) {
	if list, ok := this.Data.Load(ns); ok {
		return list.([]*v1.Deployment), nil
	}
	return nil, errors.New("record not found")
}

type DepHandler struct {
}

func (this *DepHandler) OnAdd(obj interface{}) {
	dep := obj.(*v1.Deployment)
	DepMap.Add(dep)
}
func (this *DepHandler) OnUpdate(oldObj interface{}, newObj interface{}) {
	err := DepMap.Update(newObj.(*v1.Deployment))
	if err != nil {
		log.Println(err)
	}
}
func (this *DepHandler) OnDelete(obj interface{}) {
	DepMap.Delete(obj.(*v1.Deployment))
}

var DepMap *DeploymentMap

func init() {
	DepMap = &DeploymentMap{}

}

func InitDeployment() {
	factory := informers.NewSharedInformerFactory(lib.K8sClient, 0)
	depInformer := factory.Apps().V1().Deployments().Informer()
	depInformer.AddEventHandler(&DepHandler{})
	factory.Start(wait.NeverStop)
}
