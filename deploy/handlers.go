package deploy

import (
	"github.com/gin-gonic/gin"
	"k8s-manger-v1/lib"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func RegHandlers(r *gin.Engine) {
	r.POST("/update/deployment/scale", incrReplicas)
}

func incrReplicas(c *gin.Context) {
	req := &struct {
		NameSpace  string `json:"ns"`
		Deployment string `json:"deployment"`
		Dec        bool   `json:"dec"`
	}{}
	lib.CheckError(c.ShouldBindJSON(req))

	deploy, err := lib.K8sClient.AppsV1().Deployments(req.NameSpace).Get(c, req.Deployment, metav1.GetOptions{})
	lib.CheckError(err)
	replicas := deploy.Spec.Replicas
	if req.Dec {
		*replicas--
	} else {
		*replicas++
	}
	_, err = lib.K8sClient.AppsV1().Deployments(req.NameSpace).Update(c, deploy, metav1.UpdateOptions{})
	lib.CheckError(err)
	lib.Success("Ok", c)
}
