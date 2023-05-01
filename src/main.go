package main

import (
	"fmt"
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
	fmt.Println(port)

	kubeConfigPath := os.Getenv("KUBECONFIG") //running without docker requires 'export KUBECONFIG="$HOME/.kube/config"'
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	utils.ErrLogger(err)
	clientSet, err := kubernetes.NewForConfig(config)
	utils.ErrLogger(err)

	guacUser, guacPassword, _ := guacamole.GetGuacamoleSecret(*clientSet) // TODO HANDLE ERROR
	guacBaseAddress := guacamole.GetGuacamoleBaseAddress(*clientSet)
	guac := guacamole.Guacamole{
		Username: guacUser,
		Password: guacPassword,
		BaseUrl:  guacBaseAddress,
	}

	fmt.Println("username: " + guac.Username)
	fmt.Println("pass: " + guac.Password)
	fmt.Println("url: " + guac.BaseUrl)

	/*controller := api_endpoints.Controller{ClientSet: clientSet, Guacamole: guac}

	// Creates a router without any middleware by default
	r := gin.New()

	// Global middleware
	r.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r = createRouterGroups(r, controller)
	r.Run(port)*/
}

func createRouterGroups(r *gin.Engine, controller api_endpoints.Controller) *gin.Engine {
	user := r.Group("/user/")
	{
		user.GET("", controller.GetUser)
		user.GET("/challenges/", controller.GetUserChallenges)
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
