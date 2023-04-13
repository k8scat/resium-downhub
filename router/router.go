package router

import (
	"github.com/gin-gonic/gin"
	"github.com/k8scat/downhub/controllers"
	"github.com/k8scat/downhub/middlewares"
)

func Run() {
	router := gin.Default()
	router.Use(middlewares.Auth())
	parse := router.Group("/parse")
	{
		parse.POST("/wenku", controllers.Wenku)
	}

	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
