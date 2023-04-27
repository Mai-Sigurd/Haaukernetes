package api_endpoints

import (
	"k8s-project/kali"

	"github.com/gin-gonic/gin"
)

type Kali struct {
	// Name
	// in: string
	Name string `json:"name"`

	// Message m
	// in: string
	Message string `json:"message"`
}

// PostKali godoc
// @Summary Creates Kali based on given user
// @Produce json
// @Param user body User true "User"
// @Success 200 {object} Kali
// @Router /kali/ [post]
func (c Controller) PostKali(ctx *gin.Context) {

	var body User
	if err := ctx.BindJSON(&body); err != nil {
		message := "bad request"
		ctx.JSON(400, ErrorResponse{Message: message})
	} else {
		kali.StartKali(*c.ClientSet, body.Name)
		message := "You can now rdp into your Kali."
		kaliresp := Kali{Name: body.Name, Message: message}
		ctx.JSON(200, kaliresp)
	}
}
