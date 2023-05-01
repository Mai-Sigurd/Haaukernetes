package utils

import (
	"context"
	"errors"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"strings"
	"time"
)

func FindDeployment(clientSet kubernetes.Clientset, namespace string, deploymentName string) (*appsv1.Deployment, error) {
	return clientSet.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
}

func WaitForDeployment(clientSet kubernetes.Clientset, namespace string, deploymentName string) (string, error) {
	for {
		deployment, _ := FindDeployment(clientSet, namespace, deploymentName) // TODO HANLE ERROR
		if deployment.Status.ReadyReplicas == *deployment.Spec.Replicas {
			return "", nil
		}
		time.Sleep(2 * time.Second)
	}
}

func FindPod(clientSet kubernetes.Clientset, namespace string, podName string) (v1.Pod, error) {
	pods, _ := clientSet.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{}) // TODO handle error
	for _, pod := range pods.Items {
		if strings.Contains(pod.Name, podName) {
			return pod, nil
		}
	}
	return v1.Pod{}, errors.New("could not find pod")
}

func WaitForPodReady(clientSet kubernetes.Clientset, namespace string, podName string) error {
	for {
		pod, _ := FindPod(clientSet, namespace, podName) // TODO HANDLE ERROR
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

func WaitForServiceReady(clientSet kubernetes.Clientset, namespace string, serviceName string) error {
	for {
		updatedService, _ := FindService(clientSet, namespace, serviceName) // TODO Handle error
		if updatedService.Spec.ClusterIP != "" {
			InfoLogger.Println("Service " + serviceName + " is ready")
			return nil
		}
		InfoLogger.Println("Service " + serviceName + " is not ready")
		time.Sleep(2 * time.Second)
	}
}
