package apis

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type Kali struct {
	// Namespace name
	// in: string
	Name string `json:"name"`

	// Ipaddress ip
	// in: string
	Ip string `json:"ip"`
}

// GetKali godoc
// @Summary Retrieves kali ip based on namespace name
// @Produce json
// @Param name path string true "Namespace name"
// @Success 200 {object} Kali
// @Router /kali/{name} [get]
func GetKali(ctx *gin.Context) {
	//TODO get the kali ip
	name := ctx.Param("name")
	fmt.Print(name)
	kali := Kali{Name: name, Ip: "ip addreess"}
	ctx.JSON(200, kali)
}

// PostKali godoc
// @Summary Creates Kali based on given namespace name
// @Produce json
// @Param name path string true "Namespace name"
// @Success 200 {object} Kali
// @Router /kali/{name} [post]
func PostKali(ctx *gin.Context) {
	//TODO
	name := ctx.Param("name")
	fmt.Print(name)
	kali := Kali{Name: name, Ip: "ip addreess"}
	ctx.JSON(200, kali)
}
