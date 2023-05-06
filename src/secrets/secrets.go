package secrets

import (
	"context"
	"log"
	"os"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func CreateWireGuardSecret(clientSet kubernetes.Clientset, namespace string, privatekey string) error {
	data := make(map[string][]byte)
	data["privatekey"] = []byte(privatekey)
	secret := configureSecret("wg-secret", namespace, v1.SecretTypeOpaque, data)
	return CreateSecret(clientSet, namespace, secret)
}

func CreateImageRepositorySecret(clientSet kubernetes.Clientset, namespace string) error {
	secretPath := os.Getenv("DO_SECRET_PATH") //running without docker requires 'export DO_SECRET_PATH="$HOME/do_secret"'
	dockerconfigjson, err := os.ReadFile(secretPath)

	if err != nil {
		return err
	}

	data := make(map[string][]byte)
	data[".dockerconfigjson"] = dockerconfigjson
	secret := configureSecret("regcred", namespace, v1.SecretTypeDockerConfigJson, data)
	return CreateSecret(clientSet, namespace, secret)
}

func CreateSecret(clientSet kubernetes.Clientset, namespace string, secret v1.Secret) error {
	secretsClient := clientSet.CoreV1().Secrets(namespace)
	result, err := secretsClient.Create(context.TODO(), &secret, metav1.CreateOptions{})

	if err != nil {
		return err
	}

	log.Printf("Created secret %q.\n", result.GetObjectMeta().GetName())
	return nil
}

func configureSecret(name string, namespace string, secretType v1.SecretType, data map[string][]byte) v1.Secret {
	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Type: secretType,
		Data: data,
	}

	return *secret
}
