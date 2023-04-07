package namespaces

import (
	"context"
	"k8-project/utils"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"log"
)

func CreateNamespace(clientSet kubernetes.Clientset, name string) {
	namespace := &apiv1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
	newNamespace, err := clientSet.CoreV1().Namespaces().Create(context.TODO(), namespace, metav1.CreateOptions{})
	utils.ErrHandler(err)
	log.Printf("Created namespace with name %s\n", newNamespace.Name)
}

func NamespaceExists(clientSet kubernetes.Clientset, name string) bool {
	_, err := clientSet.CoreV1().Namespaces().Get(context.TODO(), name, metav1.GetOptions{})
	return err == nil
}

func GetAllNamespaces(clientSet kubernetes.Clientset) *apiv1.NamespaceList {
	list, _ := clientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	return list
}
