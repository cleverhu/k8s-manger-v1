package core

import (
	"errors"
	"fmt"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"sync"
)

type RSMapStruct struct {
	Data sync.Map //key ns value []*corev1.Pod
}

func (this *RSMapStruct) Add(rs *v1.ReplicaSet) {
	key := rs.Namespace
	if value, ok := this.Data.Load(key); ok {
		value = append(value.([]*v1.ReplicaSet), rs)
		this.Data.Store(key, value)
	} else {
		this.Data.Store(key, []*v1.ReplicaSet{rs})
	}
}

func (this *RSMapStruct) Delete(rs *v1.ReplicaSet) {
	key := rs.Namespace
	if value, ok := this.Data.Load(key); ok {
		for index, r := range value.([]*v1.ReplicaSet) {
			if r.Name == rs.Name {
				value = append(value.([]*v1.ReplicaSet)[0:index], value.([]*v1.ReplicaSet)[index+1:]...)
				this.Data.Store(key, value)
				return
			}
		}
	}
}

func (this *RSMapStruct) Update(rs *v1.ReplicaSet) error {
	key := rs.Namespace
	if value, ok := this.Data.Load(key); ok {
		for index, r := range value.([]*v1.ReplicaSet) {
			if r.Name == rs.Name {
				value.([]*v1.ReplicaSet)[index] = rs
				this.Data.Store(key, value)
				return nil
			}
		}
	}

	return fmt.Errorf("rs-%s not found", rs.Name)
}

func (this *RSMapStruct) ListByNS(ns string) ([]*v1.ReplicaSet, error) {
	if list, ok := this.Data.Load(ns); ok {
		return list.([]*v1.ReplicaSet), nil
	}
	return nil, errors.New("rs record not found")
}

func (this *RSMapStruct) GetPodLabelsByDeployment(deploy *v1.Deployment) (string, error) {
	rs, err := this.ListByNS(deploy.Namespace)
	if err != nil {
		return "", err
	}
	podLabel := ""
	for _, item := range rs {
		if item.Annotations["deployment.kubernetes.io/revision"] != deploy.Annotations["deployment.kubernetes.io/revision"] {
			continue
		}
		for _, v := range item.OwnerReferences {
			if v.Name == deploy.Name {
				s, _ := metav1.LabelSelectorAsSelector(item.Spec.Selector)
				podLabel = s.String()
				return podLabel, nil
			}
		}
	}
	return podLabel, nil
}

type RSHandler struct {
}

func (this *RSHandler) OnAdd(obj interface{}) {
	RSMap.Add(obj.(*v1.ReplicaSet))
}
func (this *RSHandler) OnUpdate(oldObj interface{}, newObj interface{}) {
	err := RSMap.Update(newObj.(*v1.ReplicaSet))
	if err != nil {
		log.Println(err)
	}
}
func (this *RSHandler) OnDelete(obj interface{}) {
	RSMap.Delete(obj.(*v1.ReplicaSet))
}

var RSMap *RSMapStruct

func init() {
	RSMap = &RSMapStruct{}
}
