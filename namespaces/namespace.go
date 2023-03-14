package namespaces

import (
	"context"
	"fmt"
	"k8-project/utils"
	"k8s.io/api/core/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "log"
)

func CreateNamespace(clientSet kubernetes.Clientset, name string) {
	namespace := &apiv1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
	newNamespace, err := clientSet.CoreV1().Namespaces().Create(context.TODO(), namespace, metav1.CreateOptions{})
	utils.ErrHandler(err)
	fmt.Printf("Created namespace with name %s\n", newNamespace.Name)
}

func NamespaceExists(clientSet kubernetes.Clientset, name string) bool {
	_, err := clientSet.CoreV1().Namespaces().Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return false
	}
	return true
}

func GetAllNamespaces(clientSet kubernetes.Clientset) *v1.NamespaceList {
	list, _ := clientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	return list
}

// DeleteAllNamespaces
func DeleteAllNamespaces(clientSet kubernetes.Clientset) {
	list := GetAllNamespaces(clientSet)
	for _, d := range list.Items {
		err := clientSet.CoreV1().Namespaces().Delete(context.TODO(), d.Name, metav1.DeleteOptions{})
		utils.ErrHandler(err)
	}

}
