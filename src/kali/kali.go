package kali

import (
	"context"
	"k8s-project/deployments"
	"k8s-project/services"
	"k8s-project/utils"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"time"

	"k8s.io/client-go/kubernetes"
)

const kaliRDPPort = 13389
const kaliImageName = "kali"

func StartKali(clientSet kubernetes.Clientset, namespace string) {
	StartKaliImage(clientSet, namespace)
}

func StartKaliImage(clientSet kubernetes.Clientset, namespace string) (string, int32) {
	utils.InfoLogger.Println("Starting Kali")
	podLabels := make(map[string]string)
	podLabels[utils.KaliPodLabelKey] = utils.KaliPodLabelValue
	podLabels[utils.NetworkPolicyLabelKey] = utils.NetworkPolicyLabelValue
	ports := []int32{kaliRDPPort}
	deployments.CreateDeployment(clientSet, namespace, "kali", kaliImageName, ports, podLabels)
	// TODO sksal jeg også vente på deployment ready? Laver den en pod

	service := services.CreateService(clientSet, namespace, "kali", ports)
	waitForServiceReady(clientSet, service)
	ip := service.Spec.ClusterIP
	port := service.Spec.Ports[0].Port // TODO is it in the form something:something? Then this might not work
	return ip, port
}

func waitForServiceReady(clientSet kubernetes.Clientset, service *v1.Service) {
	for {
		updatedService, err := clientSet.CoreV1().Services(service.Namespace).Get(context.TODO(), service.Name, metav1.GetOptions{})
		if err != nil {
			log.Fatalf("Error getting service: %v", err)
		}

		if updatedService.Spec.ClusterIP != "" {
			break
		}

		time.Sleep(2 * time.Second)
	}
}
