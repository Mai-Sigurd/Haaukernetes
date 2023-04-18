package kali

import (
	"fmt"
	"k8-project/deployments"
	"k8-project/services"

	"k8s.io/client-go/kubernetes"
)

const kaliRDPPort = 13389
const kaliImageName = "kali"

func PostKaliKubernetes(clientSet kubernetes.Clientset, namespace string) {
	fmt.Println("Starting Kali")
	podLabels := make(map[string]string)
	podLabels["app"] = "kali"
	ports := []int32{kaliRDPPort}
	deployments.CreateDeployment(clientSet, namespace, "kali", kaliImageName, ports, podLabels)
	services.CreateService(clientSet, namespace, "kali", ports)
}
