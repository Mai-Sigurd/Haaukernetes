package main

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"k8-project/apis"
	_ "k8-project/docs"
	"k8-project/utils"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
)

func main() {
	home := homedir.HomeDir()
	kubeConfigPath := filepath.Join(home, ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	utils.ErrHandler(err)
	clientSet, err := kubernetes.NewForConfig(config)
	utils.ErrHandler(err)
	controller := apis.Controller{ClientSet: clientSet}

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
		namespace.GET("/:name", controller.GetNamespace)
		namespace.POST("/", controller.PostNamespace)
		namespace.DELETE("/", controller.DeleteNamespace)
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

	//TODO wireguard api?

	//TODO guac api?

	r.Run(":5012")

}
