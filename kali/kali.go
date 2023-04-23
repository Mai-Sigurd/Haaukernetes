package kali

import (
	"k8-project/deployments"
	"k8-project/services"
	"k8-project/utils"

	"k8s.io/client-go/kubernetes"
)

const kaliRDPPort = 13389
const kaliImageName = "kali"

func StartKali(clientSet kubernetes.Clientset, namespace string) {
	StartKaliImage(clientSet, namespace, kaliImageName)
}

func StartKaliImage(clientSet kubernetes.Clientset, namespace string, kaliImage string) {
	utils.InfoLogger.Println("Starting Kali")
	podLabels := make(map[string]string)
	podLabels["app"] = "kali"
	ports := []int32{kaliRDPPort}
	deployments.CreateDeployment(clientSet, namespace, "kali", kaliImageName, ports, podLabels)
	services.CreateService(clientSet, namespace, "kali", ports)
}
