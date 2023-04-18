package api_endpoints

import (
	"github.com/gin-gonic/gin"
	"k8-project/wireguard"
)

type Wireguard struct {
	Key       string `json:"key"`
	Namespace string `json:"namespace"`
}

type ConfigFile struct {
	File string `json:"file"`
}

// PostWireguard godoc
// @Summary Sends a public key and starts Wireguard
// @Produce json
// @Param publicKey body Wireguard true "Wireguard"
// @Success 200 {object} ConfigFile
// @Router /wireguard/ [post]
func (c Controller) PostWireguard(ctx *gin.Context) {
	var body Wireguard
	if err := ctx.BindJSON(&body); err != nil {
		message := "bad request"
		ctx.JSON(400, ErrorResponse{Message: message})
	}
	str := wireguard.PostWireguard(*c.ClientSet, body.Namespace, body.Key, c.Endpoint, c.Subnet)

	ctx.JSON(200, ConfigFile{File: str})
}
