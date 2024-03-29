package challenge

import (
	"k8s-project/deployments"
	"k8s-project/services"
	"k8s-project/utils"

	"k8s.io/client-go/kubernetes"
)

func CreateChallenge(clientSet kubernetes.Clientset, namespace string, challengeName string, imageName string, ports []int32) error {
	podLabels := make(map[string]string)
	podLabels["app"] = challengeName
	podLabels[utils.ChallengePodLabelKey] = utils.ChallengePodLabelValue

	err := deployments.CreateDeployment(clientSet, namespace, challengeName, imageName, ports, podLabels)
	if err != nil {
		utils.ErrLogger(err)
		return err
	}

	_, err = services.CreateService(clientSet, namespace, challengeName, ports)
	if err != nil {
		utils.ErrLogger(err)
		return err
	}

	return nil
}

func DeleteChallenge(clientSet kubernetes.Clientset, namespace string, challengeName string) (bool, bool) {
	deploymentDeleteStatus := deployments.DeleteDeployment(clientSet, namespace, challengeName)
	serviceDeleteStatus := services.DeleteService(clientSet, namespace, challengeName)
	return deploymentDeleteStatus, serviceDeleteStatus
}

func ChallengeExists(clientSet kubernetes.Clientset, namespace string, challengeName string) bool {
	return deployments.CheckIfDeploymentExists(clientSet, namespace, challengeName)
}
