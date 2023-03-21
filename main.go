package main

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"k8-project/apis"
	_ "k8-project/docs"
)

func main() {

	// Creates a router without any middleware by default
	r := gin.New()

	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default gin.DefaultWriter = os.Stdout
	r.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	namespace := r.Group("/namespace/")
	{
		namespace.GET("/:name", apis.GetNamespace)
		namespace.POST("/", apis.PostNamespace)
		namespace.DELETE("/", apis.DeleteNamespace)
	}

	challenge := r.Group("/challenge/")
	{
		challenge.POST("/", apis.PostChallenge)
		challenge.DELETE("/", apis.DeleteChallenge)
	}

	r.Run(":5000")

}
