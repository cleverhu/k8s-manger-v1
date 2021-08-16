package core

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"log"
	"sync"
)

type PodMapStruct struct {
	Data sync.Map //key ns value []*corev1.Pod
}

func (this *PodMapStruct) Add(pod *corev1.Pod) {
	key := pod.Namespace
	if value, ok := this.Data.Load(key); ok {
		value = append(value.([]*corev1.Pod), pod)
		this.Data.Store(key, value)
	} else {
		this.Data.Store(key, []*corev1.Pod{pod})
	}
}

func (this *PodMapStruct) Delete(pod *corev1.Pod) {
	key := pod.Namespace
	if value, ok := this.Data.Load(key); ok {
		for index, p := range value.([]*corev1.Pod) {
			if p.Name == pod.Name {
				value = append(value.([]*corev1.Pod)[0:index], value.([]*corev1.Pod)[index+1:]...)
				this.Data.Store(key, value)
				return
			}
		}
	}
}

func (this *PodMapStruct) Update(pod *corev1.Pod) error {
	key := pod.Namespace
	if value, ok := this.Data.Load(key); ok {

		for index, p := range value.([]*corev1.Pod) {
			if p.Name == pod.Name {
				value.([]*corev1.Pod)[index] = pod
				this.Data.Store(key, value)
				return nil
			}
		}
	}

	return fmt.Errorf("pod-%s not found", pod.Name)
}

func (this *PodMapStruct) ListByNS(ns string) ([]*corev1.Pod, error) {

	if ns != "" {
		if list, ok := this.Data.Load(ns); ok {
			return list.([]*corev1.Pod), nil
		}
	}
	return nil, fmt.Errorf("pods not found")
}

func (this *PodMapStruct) ListByLabelsAndNS(ns, labels string) ([]*corev1.Pod, error) {
	//fmt.Println("|||")
	//defer func() {
	//	fmt.Println("|||")
	//}()
	pods, err := this.ListByNS(ns)
	if err != nil {
		return nil, err
	}
	ret := make([]*corev1.Pod, 0)
	for _, pod := range pods {
		//fmt.Println("pod",GetLabels(pod.Labels) ,"|",labels,GetLabels(pod.Labels) == labels)
		if GetLabels(pod.Labels) == labels {
			ret = append(ret, pod)
		}
	}

	return ret, nil
}

type PodHandler struct {
}

func (this *PodHandler) OnAdd(obj interface{}) {
	PodMap.Add(obj.(*corev1.Pod))
}
func (this *PodHandler) OnUpdate(oldObj interface{}, newObj interface{}) {
	err := PodMap.Update(newObj.(*corev1.Pod))
	if err != nil {
		log.Println(err)
	}
}
func (this *PodHandler) OnDelete(obj interface{}) {
	PodMap.Delete(obj.(*corev1.Pod))
}

var PodMap *PodMapStruct

func init() {
	PodMap = &PodMapStruct{}
}
