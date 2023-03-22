package secrets

import (
	"context"
	"fmt"
	"k8-project/utils"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func CreateSecret(clientSet kubernetes.Clientset, teamName string) {
	secret := configureSeceret(teamName) //PLACEHOLDER, use
	secretsClient := clientSet.CoreV1().Secrets(teamName)
	result, err := secretsClient.Create(context.TODO(), &secret, metav1.CreateOptions{})
	utils.ErrHandler(err)
	fmt.Printf("Created secret %q.\n", result.GetObjectMeta().GetName())
}

func configureSeceret(namespace string) v1.Secret {
	secret := &v1.Secret{}

	//do stuff

	return *secret
}
