package apis

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type Namespace struct {
	// Namespace name
	// in: string
	Name string `json:"name"`
}

// GetNamespace GetUser godoc
// @Summary Retrieves namespace based on given name
// @Produce json
// @Param name path string true "Namespace name"
// @Success 200 {object} Namespace
// @Router /namespace/{name} [get]
func GetNamespace(ctx *gin.Context) {
	name := ctx.Param("name")
	fmt.Print(name)
	namespace := Namespace{Name: "hello"}
	ctx.JSON(200, namespace)
}
