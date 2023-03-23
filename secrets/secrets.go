package secrets

import (
	"context"
	"fmt"
	"k8-project/utils"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func CreateWireGuardSecret(clientSet kubernetes.Clientset, teamName string, privatekey string) {
	data := make(map[string][]byte)
	data["privatekey"] = []byte(privatekey)
	secret := configureSecret("wg-secret", teamName, v1.SecretTypeOpaque, data)
	CreateSecret(clientSet, teamName, secret)
}

func CreateSecret(clientSet kubernetes.Clientset, teamName string, secret v1.Secret) {
	secretsClient := clientSet.CoreV1().Secrets(teamName)
	result, err := secretsClient.Create(context.TODO(), &secret, metav1.CreateOptions{})
	utils.ErrHandler(err)
	fmt.Printf("Created secret %q.\n", result.GetObjectMeta().GetName())
}

func configureSecret(name string, namespace string, secretType v1.SecretType, data map[string][]byte) v1.Secret {
	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			// Labels: map[string]string{
			// 	"app": name,
			// },
			Namespace: namespace,
		},
		Type: secretType,
		Data: data,
	}

	return *secret
}
