package api_endpoints

import (
	"github.com/gin-gonic/gin"
	"k8s-project/kali"
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

	var body User // TODO der skal opdateres her så den får navn og password med i body
	if err := ctx.BindJSON(&body); err != nil {
		message := "bad request"
		ctx.JSON(400, ErrorResponse{Message: message})
	} else {
		// Start Kali
		kali.StartKali(*c.ClientSet, body.Name)
		// Find Kali info

		// Får auth token
		_, _ = c.Guacamole.GetAuthToken() // TODO handle error

		// laver user

		// laver connection

		// assigner connection til user

		// info besked tilbage om at logge ind

		message := "You can now rdp into your Kali."
		kaliresp := Kali{Name: body.Name, Message: message}
		ctx.JSON(200, kaliresp)
	}
}
