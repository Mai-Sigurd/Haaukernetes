package api_endpoints

import (
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
	str := wireguard.StartWireguard(*c.ClientSet, body.User, body.Key, c.Endpoint, c.Subnet)

	ctx.JSON(200, ConfigFile{File: str})
}
