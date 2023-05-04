package main

import (
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"k8s-project/connections/browser/guacamole"
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

	kubeConfigPath := os.Getenv("KUBECONFIG")
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	utils.ErrLogger(err)
	clientSet, err := kubernetes.NewForConfig(config)
	utils.ErrLogger(err)

	guac, err := setupGuacamole(*clientSet)
	if err != nil {
		utils.ErrLogger(err)
		return
	}
	controller := api_endpoints.Controller{ClientSet: clientSet, Guacamole: guac}

	// Creates a router without any middleware by default
	r := gin.New()

	// Global middleware
	r.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r = createRouterGroups(r, controller)
	err = r.Run(port)
	if err != nil {
		utils.ErrLogger(err)
		return
	}
}

func setupGuacamole(clientSet kubernetes.Clientset) (guacamole.Guacamole, error) {
	guacBaseAddress, err := guacamole.GetGuacamoleBaseAddress(clientSet)
	if err != nil {
		return guacamole.Guacamole{}, err
	}

	// Default guacamole credentials used to initially change the admin password
	guac := guacamole.Guacamole{
		Username: "guacadmin",
		Password: "guacadmin",
		BaseUrl:  guacBaseAddress,
	}

	err = guac.UpdateDefaultGuacAdminPassword(clientSet, "guacadmin")
	if err != nil {
		return guacamole.Guacamole{}, err
	}

	return guac, nil
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
