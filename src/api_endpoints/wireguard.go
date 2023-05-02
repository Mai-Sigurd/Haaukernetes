package api_endpoints

import (
	"k8s-project/utils"
	"k8s-project/wireguard"

	"github.com/gin-gonic/gin"
)

type Wireguard struct {
	Key  string `json:"key"`
	User string `json:"namespace"`
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
	str, err := wireguard.StartWireguard(*c.ClientSet, body.User, body.Key, utils.WireguardEndpoint, utils.WireguardSubnet)
	if err != nil {
		ctx.JSON(400, ErrorResponse{Message: err.Error()})
	}
	ctx.JSON(200, ConfigFile{File: str})
}
