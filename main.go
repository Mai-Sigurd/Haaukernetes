package main

import (
	"os"

	"k8s.io/client-go/kubernetes"

	"k8-project/api_endpoints"
	_ "k8-project/docs"
	"k8-project/utils"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	utils.SetLog()

	port := ":33333" //hardcoded because getting user input in docker is not convenient

	kubeConfigPath := os.Getenv("KUBECONFIG") //running without docker requires 'export KUBECONFIG="$HOME/.kube/config"'
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	utils.ErrLogger(err)
	clientSet, err := kubernetes.NewForConfig(config)
	utils.ErrLogger(err)

	settings := utils.ReadYaml("settings.yaml")
	controller := api_endpoints.Controller{ClientSet: clientSet, Endpoint: settings.Endpoint, Subnet: settings.Subnet}

	// Creates a router without any middleware by default
	r := gin.New()

	// Global middleware
	r.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r = createRouterGroups(r, controller)
	r.Run(port)
}

func createRouterGroups(r *gin.Engine, controller api_endpoints.Controller) *gin.Engine {
	namespace := r.Group("/namespace/")
	{
		namespace.GET("/:name", controller.GetNamespace)
		namespace.GET("/pods/:name", controller.GetNamespacePods)
		namespace.POST("/", controller.PostNamespace)
		namespace.DELETE("/", controller.DeleteNamespace)
	}

	namespaces := r.Group("/namespaces/")
	{
		namespaces.GET("", controller.GetNamespaces)
	}

	challenge := r.Group("/challenge/")
	{
		challenge.POST("/", controller.PostChallenge)
		challenge.DELETE("/", controller.DeleteChallenge)
	}

	kali := r.Group("/kali/")
	{
		kali.POST("/:namespace", controller.PostKali)
		kali.GET("/:namespace", controller.GetKali)
	}

	wireguard := r.Group("/wireguard/")
	{
		wireguard.POST("/", controller.PostWireguard)
	}
	return r
}
