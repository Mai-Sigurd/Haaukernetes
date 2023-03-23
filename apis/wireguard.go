package apis

import (
	"github.com/gin-gonic/gin"
)

type PublicKey struct {
	Key string `json:"key"`
}
type ConfigFile struct {
	File string `json:"file"`
}

// PostPublicKey godoc
// @Summary Sends a public key to Wireguard
// @Produce json
// @Param publicKey body PublicKey true "PublicKey"
// @Success 200 {object} ConfigFile
// @Router /wireguard/ [post]
func (c Controller) PostPublicKey(ctx *gin.Context) {
	var body PublicKey
	if err := ctx.BindJSON(&body); err != nil {
		message := "bad request"
		ctx.JSON(400, ErrorResponse{Message: message})
	}
	// TODO acutal wireguard functionoliaty
	message := "function not done yet"
	ctx.JSON(400, ErrorResponse{Message: message})
}
