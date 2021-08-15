package main

import (
	"github.com/gin-gonic/gin"
	"k8s-manger-v1/core"
	"k8s-manger-v1/deploy"
	"k8s-manger-v1/lib"

	"net/http"
)

func main() {
	r := gin.New()

	r.Use(func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.AbortWithStatusJSON(400, gin.H{"error": err})
			}
		}()
		c.Next()
	})

	r.Static("/static", "./static")
	r.LoadHTMLGlob("html/**/*")

	deploy.RegHandlers(r)

	r.GET("/deployments", func(c *gin.Context) {
		c.HTML(http.StatusOK, "deploy_list.html",
			lib.DataBuilder().
				SetTitle("deployment列表").
				SetData("DepList", deploy.ListAll("default")))
	})

	r.GET("/deployments/:name", func(c *gin.Context) {
		c.HTML(http.StatusOK, "deploy_detail.html",
			lib.DataBuilder().
				SetTitle("deployment详细信息-"+c.Param("name")).
				SetData("DepDetail", deploy.Detail("default", c.Param("name"))))
	})

	core.InitDeployment()
	r.Run(":80")
}
