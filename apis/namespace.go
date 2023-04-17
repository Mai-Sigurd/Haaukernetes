package apis

import (
	"context"
	"k8-project/namespaces"
	"k8-project/netpol"
	"k8-project/secrets"

	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
func (c Controller) GetNamespace(ctx *gin.Context) {
	name := ctx.Param("name")
	if !namespaces.NamespaceExists(*c.ClientSet, name) {
		message := "Sorry namespace " + name + " does not exist"
		ctx.JSON(400, ErrorResponse{Message: message})
	} else {
		ctx.JSON(200, Namespace{name})
	}
}

// PostNamespace godoc
// @Summary Creates namespace based on given name
// @Produce json
// @Param namespace body Namespace true "Namespace"
// @Success 200 {object} Namespace
// @Router /namespace/ [post]
func (c Controller) PostNamespace(ctx *gin.Context) {
	var body Namespace
	if err := ctx.BindJSON(&body); err != nil {
		message := "bad request"
		ctx.JSON(400, ErrorResponse{Message: message})
	} else if namespaces.NamespaceExists(*c.ClientSet, body.Name) {
		message := "Sorry namespace " + body.Name + " already exists"
		ctx.JSON(400, ErrorResponse{Message: message})
	} else {
		err := namespaces.CreateNamespace(*c.ClientSet, body.Name)
		if err != nil {
			ctx.JSON(400, ErrorResponse{Message: err.Error()})
		}
		netpol.CreateEgressPolicy(*c.ClientSet, body.Name)
		netpol.CreateChallengeIngressPolicy(*c.ClientSet, body.Name)
		secrets.CreateImageRepositorySecret(*c.ClientSet, body.Name)
		ctx.JSON(200, body)
	}
}

// DeleteNamespace godoc
// @Summary Deletes namespace based on given name
// @Produce json
// @Param name path string true "Namespace name"
// @Success 200
// @Router /namespace/{name} [delete]
func (c Controller) DeleteNamespace(ctx *gin.Context) {
	name := ctx.Param("name")
	err := c.ClientSet.CoreV1().Namespaces().Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		ctx.JSON(400, ErrorResponse{Message: err.Error()})
	} else {
		ctx.JSON(200, "Namespace "+name+" Deleted")
	}

}
