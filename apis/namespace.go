package apis

import (
	"context"
	"github.com/gin-gonic/gin"
	"k8-project/namespaces"
	"k8-project/netpol"
	"k8-project/secrets"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
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
		errKubernetes := PostNamespaceKubernetes(*c.ClientSet, body.Name)
		if errKubernetes != nil {
			ctx.JSON(400, ErrorResponse{Message: err.Error()})
		}
		ctx.JSON(200, body)
	}
}

func PostNamespaceKubernetes(clientSet kubernetes.Clientset, name string) error {
	err := namespaces.CreateNamespace(clientSet, name)
	if err != nil {
		return err
	}
	netpol.CreateEgressPolicy(clientSet, name)
	netpol.CreateChallengeIngressPolicy(clientSet, name)
	secrets.CreateImageRepositorySecret(clientSet, name)
	return nil
}

// DeleteNamespace godoc
// @Summary Deletes namespace based on given name
// @Produce json
// @Param name path string true "Namespace name"
// @Success 200
// @Router /namespace/{name} [delete]
func (c Controller) DeleteNamespace(ctx *gin.Context) {
	name := ctx.Param("name")
	err := DeleteNamespaceKubernetes(*c.ClientSet, name)
	if err != nil {
		ctx.JSON(400, ErrorResponse{Message: err.Error()})
	} else {
		ctx.JSON(200, "Namespace "+name+" Deleted")
	}
}

func DeleteNamespaceKubernetes(clientSet kubernetes.Clientset, name string) error {
	err := clientSet.CoreV1().Namespaces().Delete(context.TODO(), name, metav1.DeleteOptions{})
	return err
}
