package apis

import (
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
// @Success 200 {object} Namespace
// @Router /namespace/ [get]
func GetNamespace(ctx *gin.Context) {
	namespace := Namespace{Name: "hello"}
	ctx.JSON(200, namespace)
}
