package namespaces

import (
	"context"
	"errors"
	"k8-project/utils"
	"log"
	"regexp"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const k8sNamespaceRegex = "[a-z0-9]([-a-z0-9]*[a-z0-9])?"

func CreateNamespace(clientSet kubernetes.Clientset, name string) error {
	nameIsNotOk := !regexp.MustCompile(k8sNamespaceRegex).MatchString(name)
	if nameIsNotOk {
		return errors.New("invalid namespace name: a lowercase RFC 1123 label must consist of lower case alphanumeric characters or '-', and must start and end with an alphanumeric character (e.g. 'my-name',  or '123-abc'")
	}
	namespace := &apiv1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
	newNamespace, err := clientSet.CoreV1().Namespaces().Create(context.TODO(), namespace, metav1.CreateOptions{})
	utils.ErrHandler(err)
	log.Printf("Created namespace with name %s\n", newNamespace.Name)
	return nil
}

func NamespaceExists(clientSet kubernetes.Clientset, name string) bool {
	_, err := clientSet.CoreV1().Namespaces().Get(context.TODO(), name, metav1.GetOptions{})
	return err == nil
}

func GetAllNamespaces(clientSet kubernetes.Clientset) *apiv1.NamespaceList {
	list, _ := clientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	return list
}

func DeleteNamespace(clientSet kubernetes.Clientset, namespace string) {
	err := clientSet.CoreV1().Namespaces().Delete(context.TODO(), namespace, *metav1.NewDeleteOptions(0))
	utils.ErrHandler(err)
}
