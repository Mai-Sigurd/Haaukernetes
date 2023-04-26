package namespaces

import (
	"context"
	"errors"
	"k8s-project/netpol"
	"k8s-project/secrets"
	"k8s-project/utils"
	"regexp"
	"strings"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const haaukins = "Haaukins"
const k8sNamespaceRegex = "[a-z0-9]([-a-z0-9]*[a-z0-9])?"

func CreateNamespace(clientSet kubernetes.Clientset, name string) error {
	nameIsNotOk := !regexp.MustCompile(k8sNamespaceRegex).MatchString(name)
	if nameIsNotOk {
		return errors.New("invalid namespace name: a lowercase RFC 1123 label must consist of lower case alphanumeric characters or '-', and must start and end with an alphanumeric character (e.g. 'my-name',  or '123-abc'")
	}
	namespace := &apiv1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: strings.ToLower(name),
			Labels: map[string]string{
				haaukins: "",
			},
		},
	}
	newNamespace, err := clientSet.CoreV1().Namespaces().Create(context.TODO(), namespace, metav1.CreateOptions{})
	utils.ErrLogger(err)
	utils.InfoLogger.Printf("Created namespace with name %s\n", newNamespace.Name)
	return nil
}

func PostNamespace(clientSet kubernetes.Clientset, name string) error {
	err := CreateNamespace(clientSet, name)
	if err != nil {
		return err
	}
	netpol.CreateEgressPolicy(clientSet, name)
	netpol.CreateChallengeIngressPolicy(clientSet, name)
	secrets.CreateImageRepositorySecret(clientSet, name)
	return nil
}

func GetNamespaces(clientSet kubernetes.Clientset) ([]string, error) {
	namespaceList, err := clientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{LabelSelector: haaukins})
	if err != nil {
		return nil, err
	}
	var result []string
	for _, n := range namespaceList.Items {
		result = append(result, n.Name)
	}
	return result, nil
}

func NamespaceExists(clientSet kubernetes.Clientset, name string) bool {
	_, err := clientSet.CoreV1().Namespaces().Get(context.TODO(), name, metav1.GetOptions{})
	return err == nil
}

func GetNamespacePods(clientSet kubernetes.Clientset, name string) ([]string, error) {
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

func DeleteNamespace(clientSet kubernetes.Clientset, namespace string) error {
	err := clientSet.CoreV1().Namespaces().Delete(context.TODO(), namespace, *metav1.NewDeleteOptions(0))
	return err
}
