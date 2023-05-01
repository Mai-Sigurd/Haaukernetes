package kali

import (
	"context"
	"errors"
	"k8s-project/deployments"
	"k8s-project/services"
	"k8s-project/utils"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"strings"
	"time"

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
	pods, _ := clientSet.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{}) // TODO handle error
	kaliPod, _ := findKaliPod(pods)                                                          // TODO HANDLE ERROR
	_ = waitForPodReady(kaliPod)                                                             // TODO HANDLE ERROR
	service := services.CreateService(clientSet, namespace, "kali", ports)
	waitForServiceReady(clientSet, service)
	ip := service.Spec.ClusterIP
	port := service.Spec.Ports[0].Port // TODO is it in the form something:something? Then this might not work
	return ip, port
}

func findKaliPod(pods *v1.PodList) (v1.Pod, error) {
	for _, pod := range pods.Items {
		if strings.Contains(pod.Name, "kali") {
			return pod, nil
		}
	}
	return v1.Pod{}, errors.New("could not find a Kali pod")
}

func waitForPodReady(pod v1.Pod) error {
	for {
		if pod.Status.Phase == "Running" && pod.Status.Conditions[0].Status == "True" {
			return nil
		}
		time.Sleep(2 * time.Second)
	}
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
	utils.InfoLogger.Println("Kali service ready")
}
