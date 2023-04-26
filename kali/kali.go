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
	StartKaliImage(clientSet, namespace)
}

func StartKaliImage(clientSet kubernetes.Clientset, namespace string) {
	utils.InfoLogger.Println("Starting Kali")
	podLabels := make(map[string]string)
	podLabels[utils.KaliPodLabelKey] = utils.KaliPodLabelValue
	ports := []int32{kaliRDPPort}
	deployments.CreateDeployment(clientSet, namespace, "kali", kaliImageName, ports, podLabels)
	services.CreateService(clientSet, namespace, "kali", ports)
}
