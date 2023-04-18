package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/homedir"

	"k8-project/apis"
	_ "k8-project/docs"
	"k8-project/utils"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"k8s.io/client-go/tools/clientcmd"
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

	home := homedir.HomeDir()
	kubeConfigPath := filepath.Join(home, ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	utils.ErrHandler(err)
	clientSet, err := kubernetes.NewForConfig(config)
	utils.ErrHandler(err)

	settings := utils.ReadYaml("settings.yaml")
	controller := apis.Controller{ClientSet: clientSet, Endpoint: settings.Endpoint, Subnet: settings.Subnet}

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

	wireguard := r.Group("/wireguard/")
	{
		wireguard.POST("/", controller.StartWireguard)
	}

	//TODO guac api?

	r.Run(port)
}
