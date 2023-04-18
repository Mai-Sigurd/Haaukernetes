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
type Namespaces struct {
	// Namespaces names
	// in: array
	Names []string `json:"names"`
}

type Pods struct {
	// Pods names
	// in: array
	Names []string `json:"names"`
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

// GetNamespaces godoc
// @Summary Retrieves all namespaces
// @Produce json
// @Success 200 {object} Namespaces
// @Router /namespaces [get]
func (c Controller) GetNamespaces(ctx *gin.Context) {
	result, err := GetNamespacesKubernetes(*c.ClientSet)
	if err != nil {
		ctx.JSON(400, ErrorResponse{Message: err.Error()})
	} else {

		ctx.JSON(200, Namespaces{Names: result})
	}
}

func GetNamespacesKubernetes(clientSet kubernetes.Clientset) ([]string, error) {
	namespaceList, err := clientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	var result []string
	for _, n := range namespaceList.Items {
		result = append(result, n.Name)
	}
	return result, nil
}

// GetNamespacePods godoc
// @Summary Retrieves all pods in a namespace
// @Produce json
// @Param name path string true "Namespace name"
// @Success 200 {object} Pods
// @Router /namespace/pods/{name} [get]
func (c Controller) GetNamespacePods(ctx *gin.Context) {
	name := ctx.Param("name")
	result, err := GetNamespacePodsKubernetes(*c.ClientSet, name)
	if err != nil {
		ctx.JSON(400, ErrorResponse{Message: err.Error()})
	}
	ctx.JSON(200, Pods{Names: result})
}

func GetNamespacePodsKubernetes(clientSet kubernetes.Clientset, name string) ([]string, error) {
	podList, err := clientSet.CoreV1().Pods(name).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	var result []string
	for _, n := range podList.Items {
		result = append(result, n.Name)
	}
	return result, nil

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
