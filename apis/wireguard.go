package apis

import (
	"k8-project/netpol"
	"k8-project/wireguard"
	"k8s.io/client-go/kubernetes"

	"github.com/gin-gonic/gin"
)

type Wireguard struct {
	Key       string `json:"key"`
	Namespace string `json:"namespace"`
}

type ConfigFile struct {
	File string `json:"file"`
}

// StartWireguard godoc
// @Summary Sends a public key and starts Wireguard
// @Produce json
// @Param publicKey body Wireguard true "Wireguard"
// @Success 200 {object} ConfigFile
// @Router /wireguard/ [post]
func (c Controller) StartWireguard(ctx *gin.Context) {
	var body Wireguard
	if err := ctx.BindJSON(&body); err != nil {
		message := "bad request"
		ctx.JSON(400, ErrorResponse{Message: message})
	}
	file := StartWireguardKubernetes(*c.ClientSet, body.Namespace, body.Key, c.Endpoint, c.Subnet)
	ctx.JSON(200, file)
}

func StartWireguardKubernetes(clientSet kubernetes.Clientset, namespace string, key string, endpoint string, subnet string) ConfigFile {
	config := wireguard.StartWireguard(clientSet, namespace, key, endpoint, subnet)
	netpol.AddWireguardToChallengeIngressPolicy(clientSet, namespace)
	return ConfigFile{File: config}
}
