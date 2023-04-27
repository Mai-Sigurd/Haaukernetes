package kali

import (
	"k8s-project/deployments"
	"k8s-project/services"
	"k8s-project/utils"

	"k8s.io/client-go/kubernetes"
)

const kaliRDPPort = 13389
const kaliImageName = "kali"

func StartKali(clientSet kubernetes.Clientset, namespace string) error {
	return StartKaliImage(clientSet, namespace, kaliImageName)
}

func StartKaliImage(clientSet kubernetes.Clientset, namespace string, imageName string) error {
	utils.InfoLogger.Println("Starting Kali")
	podLabels := make(map[string]string)
	podLabels[utils.KaliPodLabelKey] = utils.KaliPodLabelValue
	ports := []int32{kaliRDPPort}
	err := deployments.CreateDeployment(clientSet, namespace, "kali", imageName, ports, podLabels)
	if err != nil {
		return err
	}
	err = services.CreateService(clientSet, namespace, "kali", ports)
	return err
}
