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

// GetNamespace godoc
// @Summary Retrieves namespace based on given name
// @Produce json
// @Param name path string true "Namespace name"
// @Success 200 {object} Namespace
// @Router /namespace/{name} [get]
func GetNamespace(ctx *gin.Context) {
	//TODO
	name := ctx.Param("name")
	fmt.Print(name)
	namespace := Namespace{Name: "hello"}
	ctx.JSON(200, namespace)
}

// PostNamespace godoc
// @Summary Creates namespace based on given name
// @Produce json
// @Param namespace body Namespace true "Namespace"
// @Success 200 {object} Namespace
// @Router /namespace/ [post]
func PostNamespace(ctx *gin.Context) {
	//TODO
	var namespaceBody Namespace
	if err := ctx.BindJSON(&namespaceBody); err != nil {
		// DO SOMETHING WITH THE ERROR
	}
	fmt.Println(namespaceBody)
}

// DeleteNamespace godoc
// @Summary Deletes namespace based on given name
// @Produce json
// @Param name path string true "Namespace name"
// @Success 200
// @Router /namespace/{name} [delete]
func DeleteNamespace(ctx *gin.Context) {
	//TODO
	name := ctx.Param("name")
	fmt.Print(name)
	namespace := Namespace{Name: "hello"}
	ctx.JSON(200, namespace)
}
