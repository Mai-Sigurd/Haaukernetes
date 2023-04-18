package challenge

import (
	"k8-project/deployments"
	"k8-project/services"
	"k8s.io/client-go/kubernetes"
)

func CreateChallenge(clientSet kubernetes.Clientset, namespace string, challengeName string, imageName string, ports []int32) {
	podLabels := make(map[string]string)
	podLabels["app"] = challengeName
	podLabels["type"] = "challenge"
	deployments.CreateDeployment(clientSet, namespace, challengeName, imageName, ports, podLabels)
	services.CreateService(clientSet, namespace, challengeName, ports)
}

func DeleteChallenge(clientSet kubernetes.Clientset, namespace string, challengeName string) (bool, bool) {
	deploymentDeleteStatus := deployments.DeleteDeployment(clientSet, namespace, challengeName)
	serviceDeleteStatus := services.DeleteService(clientSet, namespace, challengeName)
	return deploymentDeleteStatus, serviceDeleteStatus
}

func Challenge_exists(clientSet kubernetes.Clientset, namespace string, challengeName string) bool {
	return deployments.CheckIfDeploymentExists(clientSet, namespace, challengeName)
}