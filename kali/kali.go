package kali

import (
	"fmt"
	"k8-project/deployments"
	"k8-project/services"
	"strconv"

	"k8s.io/client-go/kubernetes"
)

const kaliRDPPort = 13389
const kaliImageName = "kali"

func StartKali(clientSet kubernetes.Clientset, namespace string) string {
	fmt.Println("Starting Kali")
	podLabels := make(map[string]string)
	podLabels["app"] = "kali"
	ports := []int32{kaliRDPPort}
	deployments.CreateDeployment(clientSet, namespace, "kali", kaliImageName, ports, podLabels)
	service := services.CreateService(clientSet, namespace, "kali", ports)
	kaliSocketAddress := service.Spec.ClusterIP + ":" + strconv.Itoa(kaliRDPPort)
	return kaliSocketAddress
}
