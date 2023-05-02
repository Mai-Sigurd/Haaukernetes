package main

import (
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"k8s-project/guacamole"
	"os"

	"k8s.io/client-go/kubernetes"

	"k8s-project/api_endpoints"
	_ "k8s-project/docs"
	"k8s-project/utils"

	"github.com/gin-gonic/gin"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	utils.SetLog()

	port := ":" + utils.APIPort

	kubeConfigPath := os.Getenv("KUBECONFIG") //running without docker requires 'export KUBECONFIG="$HOME/.kube/config"'
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	utils.ErrLogger(err)
	clientSet, err := kubernetes.NewForConfig(config)
	utils.ErrLogger(err)

	guac := setupGuacamole(*clientSet)
	controller := api_endpoints.Controller{ClientSet: clientSet, Guacamole: guac}

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

func setupGuacamole(clientSet kubernetes.Clientset) guacamole.Guacamole {
	guacUser, guacPassword, err := guacamole.GetGuacamoleSecret(clientSet)
	if err != nil {
		utils.ErrLogger(err)
	}

	guacBaseAddress, err := guacamole.GetGuacamoleBaseAddress(clientSet)
	if err != nil {
		utils.ErrLogger(err)
	}

	guac := guacamole.Guacamole{
		Username: guacUser,
		Password: guacPassword,
		BaseUrl:  guacBaseAddress,
	}
	err = guac.UpdateAdminPasswordInGuac("guacadmin")
	if err != nil {
		utils.ErrLogger(err)
	}
	return guac
}

func createRouterGroups(r *gin.Engine, controller api_endpoints.Controller) *gin.Engine {
	user := r.Group("/user/")
	{
		user.GET("", controller.GetUser)
		user.POST("", controller.PostUser)
		user.DELETE("", controller.DeleteUser)
	}

	namespaces := r.Group("/users/")
	{
		namespaces.GET("", controller.GetUsers)
	}

	challenge := r.Group("/challenge/")
	{
		challenge.POST("", controller.PostChallenge)
		challenge.DELETE("", controller.DeleteChallenge)
	}

	kali := r.Group("/kali/")
	{
		kali.POST("", controller.PostKali)
	}

	wireguard := r.Group("/wireguard/")
	{
		wireguard.POST("", controller.PostWireguard)
	}
	return r
}

// utils.InfoLogger.Printf
//
