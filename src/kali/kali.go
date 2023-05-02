package kali

import (
	"k8s-project/deployments"
	"k8s-project/services"
	"k8s-project/utils"
	"k8s.io/client-go/kubernetes"
)

const kaliRDPPort = 13389

func StartKali(clientSet kubernetes.Clientset, namespace string, imageName string) (string, int32, error) {
	utils.InfoLogger.Println("Starting Kali")
	podLabels := make(map[string]string)
	podLabels[utils.KaliPodLabelKey] = utils.KaliPodLabelValue
	podLabels[utils.NetworkPolicyLabelKey] = utils.NetworkPolicyLabelValue
	ports := []int32{kaliRDPPort}

	err := deployments.CreateDeployment(clientSet, namespace, "kali", imageName, ports, podLabels)
	if err != nil {
		return "", 0, err
	}

	err = utils.WaitForDeployment(clientSet, namespace, "kali")
	if err != nil {
		return "", 0, err
	}

	err = utils.WaitForPodReady(clientSet, namespace, "kali")
	if err != nil {
		return "", 0, err
	}

	service, err := services.CreateService(clientSet, namespace, "kali", ports)
	if err != nil {
		return "", 0, err
	}

	err = utils.WaitForServiceReady(clientSet, namespace, "kali")
	if err != nil {
		return "", 0, err
	}

	ip := service.Spec.ClusterIP
	port := service.Spec.Ports[0].Port
	return ip, port, nil
}
