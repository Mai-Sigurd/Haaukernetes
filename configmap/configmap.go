package configmap

import (
	"context"
	"fmt"
	"k8-project/utils"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func CreateConfigMap(clientSet kubernetes.Clientset, teamName string) {
	configmap := configureConfigMap(teamName) //PLACEHOLDER, use
	configMapClient := clientSet.CoreV1().ConfigMaps(teamName)
	result, err := configMapClient.Create(context.TODO(), &configmap, metav1.CreateOptions{})
	utils.ErrHandler(err)
	fmt.Printf("Created configmap %q.\n", result.GetObjectMeta().GetName())
}

func configureConfigMap(namespace string) v1.ConfigMap {
	configmap := &v1.ConfigMap{}

	//do stuff

	return *configmap
}
