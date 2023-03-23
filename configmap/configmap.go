package configmap

import (
	"context"
	"fmt"
	"k8-project/utils"
	"k8-project/wireguardconfigs"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

//should possibly live somewhere else :))
func CreateWireGuardConfigMap(clientSet kubernetes.Clientset, teamName string, serverPrivateKey string, clientPublicKey string) {
	data := make(map[string]string)
	data["wg0.conf"] = wireguardconfigs.GetServerConfig(serverPrivateKey, clientPublicKey)
	fmt.Println(data["wg0.conf"])
	configMap := configureConfigMap("wg-configmap", teamName, data)
	CreateConfigMap(clientSet, teamName, configMap)
}

func CreateConfigMap(clientSet kubernetes.Clientset, teamName string, configMap v1.ConfigMap) {
	configMapClient := clientSet.CoreV1().ConfigMaps(teamName)
	result, err := configMapClient.Create(context.TODO(), &configMap, metav1.CreateOptions{})
	utils.ErrHandler(err)
	fmt.Printf("Created configmap %q.\n", result.GetObjectMeta().GetName())
}

func configureConfigMap(name string, namespace string, data map[string]string) v1.ConfigMap {
	configmap := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Data: data,
	}

	return *configmap
}
