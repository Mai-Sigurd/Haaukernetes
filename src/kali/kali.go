package kali

import (
	"k8s-project/deployments"
	"k8s-project/services"
	"k8s-project/utils"
	"k8s.io/client-go/kubernetes"
)

const kaliRDPPort = 13389
const kaliImageName = "kali"

func StartKali(clientSet kubernetes.Clientset, namespace string) (string, int32) {
	utils.InfoLogger.Println("Starting Kali")
	podLabels := make(map[string]string)
	podLabels[utils.KaliPodLabelKey] = utils.KaliPodLabelValue
	podLabels[utils.NetworkPolicyLabelKey] = utils.NetworkPolicyLabelValue
	ports := []int32{kaliRDPPort}

	deployments.CreateDeployment(clientSet, namespace, "kali", kaliImageName, ports, podLabels)
	err := utils.WaitForDeployment(clientSet, namespace, "kali")
	if err != nil {
		utils.ErrLogger(err)
	}

	err = utils.WaitForPodReady(clientSet, namespace, "kali")
	if err != nil {
		utils.ErrLogger(err)
	}

	service := services.CreateService(clientSet, namespace, "kali", ports)
	err = utils.WaitForServiceReady(clientSet, namespace, "kali")
	if err != nil {
		utils.ErrLogger(err)
	}

	ip := service.Spec.ClusterIP
	port := service.Spec.Ports[0].Port
	return ip, port
}
