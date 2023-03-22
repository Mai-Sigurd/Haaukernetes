package secrets

import (
	"context"
	"fmt"
	"k8-project/utils"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

//vil det være federe hvis CreateSecret bare tager en secret som er præconf? så skal den ikke tage 1000 args...
//Men det stemmer jo ikke overens med måden det er sat op de andre steder..
//det er lækkert nok at struct bygges i configure ift. at alle funktionerne ikke bliver kæmpestore, men omvendt
//så kan configuresecret jo med tiden ende med at tage 7000 args fordi der er forskellige usecases... det er
//lidt et tradeoff uden et klart svar..
func CreateWireGuardSecret(clientSet kubernetes.Clientset, teamName string, privatekey string) {
	data := make(map[string][]byte)
	data["privateKey"] = []byte(privatekey)
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

//denne følger "konventionen" fulgt i de andre filer...
//men det skal nok laves om / suppleres af den anden måde at lave create, hvor struct til create blot gives som param
// func dep_CreateSecret(clientSet kubernetes.Clientset, teamName string) {
// 	secret := configureSecret(teamName) //PLACEHOLDER, use
// 	secretsClient := clientSet.CoreV1().Secrets(teamName)
// 	result, err := secretsClient.Create(context.TODO(), &secret, metav1.CreateOptions{})
// 	utils.ErrHandler(err)
// 	fmt.Printf("Created secret %q.\n", result.GetObjectMeta().GetName())
// }
