package api_endpoints

import (
	"github.com/gin-gonic/gin"
	"k8s-project/connections/browser"
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
		message, err := browser.SetupGuacAndKali(c.Guacamole, *c.ClientSet, body.Name, body.Password)

		if err != nil {
			ctx.JSON(400, ErrorResponse{Message: err.Error()})
		}

		kaliResp := Kali{Name: body.Name, Message: message}
		ctx.JSON(200, kaliResp)
	}
}
