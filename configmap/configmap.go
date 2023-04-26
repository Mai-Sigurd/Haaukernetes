package configmap

import (
	"context"
	"k8s-project/utils"
	"k8s-project/wireguardconfigs"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// should possibly live somewhere else :)) TODO what dis mean
func CreateWireGuardConfigMap(clientSet kubernetes.Clientset, teamName string, serverPrivateKey string, clientPublicKey string) {
	data := make(map[string]string)
	data["wg0.conf"] = wireguardconfigs.GetServerConfig(serverPrivateKey, clientPublicKey)
	configMap := configureConfigMap("wg-configmap", teamName, data)
	CreateConfigMap(clientSet, teamName, configMap)
}

func CreateConfigMap(clientSet kubernetes.Clientset, teamName string, configMap v1.ConfigMap) {
	configMapClient := clientSet.CoreV1().ConfigMaps(teamName)
	result, err := configMapClient.Create(context.TODO(), &configMap, metav1.CreateOptions{})
	utils.ErrLogger(err) // TODO is this a fatal error? or can we continue
	utils.InfoLogger.Printf("Created configmap %q.\n", result.GetObjectMeta().GetName())
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
