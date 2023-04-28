package secrets

import (
	"context"
	"k8s-project/utils"
	"log"
	"os"

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

func CreateImageRepositorySecret(clientSet kubernetes.Clientset, teamName string) {
	secretPath := os.Getenv("DO_SECRET_PATH") //running without docker requires 'export DO_SECRET_PATH="$HOME/do_secret"'
	dockerconfigjson, err := os.ReadFile(secretPath)
	utils.ErrLogger(err)
	data := make(map[string][]byte)
	data[".dockerconfigjson"] = dockerconfigjson
	secret := configureSecret("regcred", teamName, v1.SecretTypeDockerConfigJson, data)
	CreateSecret(clientSet, teamName, secret)
}

func CreateSecret(clientSet kubernetes.Clientset, teamName string, secret v1.Secret) {
	secretsClient := clientSet.CoreV1().Secrets(teamName)
	result, err := secretsClient.Create(context.TODO(), &secret, metav1.CreateOptions{})
	utils.ErrLogger(err)
	log.Printf("Created secret %q.\n", result.GetObjectMeta().GetName())
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
