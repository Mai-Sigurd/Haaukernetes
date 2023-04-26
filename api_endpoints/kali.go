package api_endpoints

import (
	"k8s-project/kali"

	"github.com/gin-gonic/gin"
)

type Kali struct {
	// Namespace name
	// in: string
	Namespace string `json:"namespace"`

	// Message m
	// in: string
	Message string `json:"message"`
}

// GetKali godoc
// @Summary Retrieves kali ip based on namespace name
// @Produce json
// @Param name path string true "User name"
// @Success 200 {object} Kali
// @Router /kali/{namespace} [get]
func (c Controller) GetKali(ctx *gin.Context) {
	name := ctx.Param("name")
	message := "You can now rdp into your Kali."
	kali := Kali{Namespace: name, Message: message}
	ctx.JSON(200, kali)
}

// PostKali godoc
// @Summary Creates Kali based on given namespace name
// @Produce json
// @Param name path string true "User name"
// @Success 200 {object} Kali
// @Router /kali/{namespace} [post]
func (c Controller) PostKali(ctx *gin.Context) {

	name := ctx.Param("namespace")
	kali.StartKali(*c.ClientSet, name)
	message := "You can now rdp into your Kali."
	kali := Kali{Namespace: name, Message: message}
	ctx.JSON(200, kali)
}
