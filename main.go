package main

import (
	"bufio"
	"fmt"
	"k8-project/utils"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"k8-project/api_endpoints"
	_ "k8-project/docs"
)

func main() {
	currentTime := time.Now()
	f, err := os.OpenFile(currentTime.Format("2006.01.02 15:04:05"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	fmt.Println("Write the port you want the web app to run on")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	port := scanner.Text()
	port = ":" + port

	home := homedir.HomeDir()
	kubeConfigPath := filepath.Join(home, ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	utils.ErrHandler(err)
	clientSet, err := kubernetes.NewForConfig(config)
	utils.ErrHandler(err)
	controller := api_endpoints.Controller{ClientSet: clientSet}

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

	//TODO guac api?

	r.Run(port)

}
