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

type GuacUser struct {
	// Username
	// in: string
	Name string `json:"name"`
	// Password
	// in: string
	Password string `json:"password"`
}

// PostKali godoc
// @Summary Creates Kali based on given user
// @Produce json
// @Param user body GuacUser true "Guacamole User"
// @Success 200 {object} Kali
// @Router /kali/ [post]
func (c Controller) PostKali(ctx *gin.Context) {

	var body GuacUser
	if err := ctx.BindJSON(&body); err != nil {
		message := "bad request"
		ctx.JSON(400, ErrorResponse{Message: message})
	} else {

		ip, port := kali.StartKali(*c.ClientSet, body.Name)
		_, _ = c.Guacamole.GetAuthToken()                                   // TODO handle error
		_ = c.Guacamole.CreateUser(body.Name, body.Password)                // TODO handle error
		connIdentifier, _ := c.Guacamole.CreateConnection(ip, string(port)) // TODO handle error
		_, _ = c.Guacamole.AssignConnection(connIdentifier, body.Name)      // TODO handle error

		// info besked tilbage om at logge ind

		message := "You can now rdp into your Kali."
		kaliresp := Kali{Name: body.Name, Message: message}
		ctx.JSON(200, kaliresp)
	}
}
