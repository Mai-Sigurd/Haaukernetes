package apis

import (
	"k8-project/netpol"
	"k8-project/wireguard"

	"github.com/gin-gonic/gin"
)

type Wireguard struct {
	Key       string `json:"key"`
	Namespace string `json:"namespace"`
}

type ConfigFile struct {
	File string `json:"file"`
}

// PostPublicKey godoc
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
	config := wireguard.StartWireguard(*c.ClientSet, body.Namespace, body.Key)
	netpol.AddWireguardToChallengeIngressPolicy(*c.ClientSet, body.Namespace)
	ctx.JSON(200, ConfigFile{File: config})
}
