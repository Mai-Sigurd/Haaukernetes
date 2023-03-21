package apis

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"k8-project/namespaces"
	"k8-project/netpol"
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
		existError := fmt.Errorf("\nSorry namespace %s does not exist \n ", name)
		ctx.Error(existError)
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
		ctx.Error(err)
	} else if namespaces.NamespaceExists(*c.ClientSet, body.Name) {
		existError := fmt.Errorf("\nSorry namespace %s already exists \n ", body.Name)
		ctx.Error(existError)
	} else {
		namespaces.CreateNamespace(*c.ClientSet, body.Name)
		netpol.CreateKaliEgressPolicy(*c.ClientSet, body.Name)
		netpol.CreateChallengeIngressPolicy(*c.ClientSet, body.Name)
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
	// TODO
	name := ctx.Param("name")
	err := c.ClientSet.CoreV1().Namespaces().Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		ctx.Error(err)
	} else {
		ctx.JSON(200, "Deleted")
	}

}
