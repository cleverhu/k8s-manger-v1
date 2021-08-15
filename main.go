package main

import (
	"github.com/gin-gonic/gin"
	"k8s-manger-v1/core"
	"k8s-manger-v1/deploy"
	"k8s-manger-v1/lib"

	"net/http"
)

func main() {
	//gin.SetMode(gin.DebugMode)
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
				SetData("DepList", deploy.ListAll("")))
	})

	r.POST("/deployments", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Ok", "result": deploy.ListAll("")})
	})

	r.GET("/deployments/:namespace/:name", func(c *gin.Context) {
		c.HTML(http.StatusOK, "deploy_detail.html",
			lib.DataBuilder().
				SetTitle("deployment详细信息-"+c.Param("name")).
				SetData("DepDetail", deploy.Detail(c.Param("namespace"), c.Param("name"))))
	})

	core.InitDeployment()

	r.Run(":80")
}
