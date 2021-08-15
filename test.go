package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"k8s-manger-v1/lib"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"log"
	"sync"
	"time"
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

type DepHandler struct {
}

func (this *DepHandler) OnAdd(obj interface{}) {
	dep := obj.(*v1.Deployment)
	DepMap.Add(dep)
}
func (this *DepHandler) OnUpdate(oldObj interface{}, newObj interface{}) {

	if dep, ok := oldObj.(*v1.Deployment); ok {
		fmt.Println(dep.Name)
	}
}
func (this *DepHandler) OnDelete(obj interface{}) {

}

var DepMap = &DeploymentMap{}

func main() {
	//_, c := cache.NewInformer(
	//	cache.NewListWatchFromClient(lib.K8sClient.AppsV1().RESTClient(), "deployments", "default", fields.Everything()),
	//	&v1.Deployment{},
	//	0,
	//	&DepHandler{})
	//c.Run(wait.NeverStop)
	factory := informers.NewSharedInformerFactory(lib.K8sClient, 0)
	depInformer := factory.Apps().V1().Deployments().Informer()
	depInformer.AddEventHandler(&DepHandler{})
	factory.Start(wait.NeverStop)
	c, cancel := context.WithTimeout(context.Background(), time.Second*3)
	select {
	case <-c.Done():
		cancel()
		log.Fatal("超时了")
	default:
		r := gin.New()
		r.GET("/", func(c *gin.Context) {
			result := make(map[string][]string, 0)
			DepMap.Data.Range(func(key, value interface{}) bool {
				for _, dep := range value.([]*v1.Deployment) {
					result[dep.Namespace] = append(result[dep.Namespace], dep.Name)
				}
				return true
			})
			c.JSON(200, result)
		})
		r.Run(":8080")
	}
}
