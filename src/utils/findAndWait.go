package utils

import (
	"context"
	"errors"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"strings"
	"time"
)

func FindPod(clientSet kubernetes.Clientset, namespace string, podName string) (v1.Pod, error) {
	pods, _ := clientSet.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{}) // TODO handle error
	for _, pod := range pods.Items {
		if strings.Contains(pod.Name, podName) {
			return pod, nil
		}
	}
	return v1.Pod{}, errors.New("could not find pod")
}

func WaitForPodReady(pod v1.Pod) error { // TODO går det at podden ikke bliver "hentet på ny" hver gang????
	for {
		if pod.Status.Phase == "Running" && pod.Status.Conditions[0].Status == "True" {
			InfoLogger.Println("Pod " + pod.Name + " is ready")
			return nil
		}
		InfoLogger.Println("Pod " + pod.Name + " is not ready")
		time.Sleep(2 * time.Second)
	}
}

func FindService(clientSet kubernetes.Clientset, namespace string, serviceName string) (*v1.Service, error) {
	return clientSet.CoreV1().Services(namespace).Get(context.TODO(), serviceName, metav1.GetOptions{}) // TODO Handle error
}

func WaitForServiceReady(clientSet kubernetes.Clientset, service *v1.Service) error {
	for {
		updatedService, _ := clientSet.CoreV1().Services(service.Namespace).Get(context.TODO(), service.Name, metav1.GetOptions{}) // TODO Handle error

		if updatedService.Spec.ClusterIP != "" {
			InfoLogger.Println("Service " + service.Name + " is ready")
			return nil
		}
		InfoLogger.Println("Service " + service.Name + " is not ready")
		time.Sleep(2 * time.Second)
	}
}
