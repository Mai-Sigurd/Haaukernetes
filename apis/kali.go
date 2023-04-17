package apis

import (
	"fmt"
	"k8-project/deployments"
	"k8-project/services"

	"github.com/gin-gonic/gin"
	"k8s.io/client-go/kubernetes"
)

type Kali struct {
	// Namespace name
	// in: string
	Namespace string `json:"namespace"`

	// Ipaddress ip
	// in: string
	Ip string `json:"ip"`

	// Message m
	// in: string
	Message string `json:"message"`
}

// GetKali godoc
// @Summary Retrieves kali ip based on namespace name
// @Produce json
// @Param name path string true "Namespace name"
// @Success 200 {object} Kali
// @Router /kali/{namespace} [get]
func (c Controller) GetKali(ctx *gin.Context) {
	// TODO get the kali ip - is deprecated, and will become guac based instead
	name := ctx.Param("name")
	message := "You can now vnc into your Kali. If on Mac first do `minikube service kali-vnc-expose -n <namespace>`"
	kali := Kali{Namespace: name, Ip: "ip addreess", Message: message}
	ctx.JSON(200, kali)
}

// PostKali godoc
// @Summary Creates Kali based on given namespace name
// @Produce json
// @Param name path string true "Namespace name"
// @Success 200 {object} Kali
// @Router /kali/{namespace} [post]
func (c Controller) PostKali(ctx *gin.Context) {

	name := ctx.Param("namespace")
	StartKali(*c.ClientSet, name)
	message := "You can now vnc into your Kali. If on Mac first do `minikube service kali-vnc-expose -n <namespace>`"
	kali := Kali{Namespace: name, Ip: "ip addreess", Message: message}
	ctx.JSON(200, kali)
}

// TODO maybe port 5900 with new image?! if no work, check this
func StartKali(clientSet kubernetes.Clientset, namespace string) {
	fmt.Println("Starting Kali")
	podLabels := make(map[string]string)
	podLabels["app"] = "kali-vnc"
	ports := []int32{5901}
	deployments.CreateDeployment(clientSet, namespace, "kali-vnc", ports, podLabels)
	services.CreateService(clientSet, namespace, "kali-vnc", ports)
	services.CreateExposeService(clientSet, namespace, "kali-vnc", 5901) //TODO: deprecated, no nodeport for kali?
}
