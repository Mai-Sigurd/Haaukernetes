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

	// Wait for kali pod and service to be ready
	kaliPod, _ := utils.FindPod(clientSet, namespace, "kali") // TODO HANDLE ERROR
	_ = utils.WaitForPodReady(kaliPod)                        // TODO HANDLE ERROR
	service := services.CreateService(clientSet, namespace, "kali", ports)
	_ = utils.WaitForServiceReady(clientSet, service) // TODO HANDLE ERROR
	ip := service.Spec.ClusterIP
	port := service.Spec.Ports[0].Port // TODO is it in the form something:something? Then this might not work
	return ip, port
}
