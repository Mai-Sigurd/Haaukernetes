package api_endpoints

import (
	"github.com/gin-gonic/gin"
	"k8s-project/kali"
	"k8s-project/utils"
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

		err = c.Guacamole.UpdateAuthToken()
		if err != nil {
			utils.ErrLogger(err)
		}

		err = c.Guacamole.CreateUser(body.Name, body.Password)
		if err != nil {
			utils.ErrLogger(err)
		}

		connIdentifier, err := c.Guacamole.CreateConnection(ip, port, body.Name)
		if err != nil {
			utils.ErrLogger(err)
		}

		err = c.Guacamole.AssignConnection(connIdentifier, body.Name)
		if err != nil {
			utils.ErrLogger(err)
		}

		message := "You can now RDP into your Kali by visiting the Guacamole interface at: " + c.Guacamole.BaseUrl
		utils.InfoLogger.Printf("Message to user " + body.Name + ": " + message)
		kaliResp := Kali{Name: body.Name, Message: message}
		ctx.JSON(200, kaliResp)
	}
}
