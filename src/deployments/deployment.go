package deployments

import (
	"context"
	"fmt"
	utils "k8s-project/utils"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func CreatePrebuiltDeployment(clientSet kubernetes.Clientset, namespace string, deployment *appsv1.Deployment) error {
	deploymentsClient := clientSet.AppsV1().Deployments(namespace)
	result, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})

	if err != nil {
		return err
	}
	utils.InfoLogger.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())
	return nil
}

// CreateDeployment configures a deployment and then creates a deployment from that configuration
// in the given namespace.
// -----
// NOTE ON PARAMETERS
// "name" and "imageName" are primarily separated for allowing multiple pods/deployments running the same image, as this is used in the tests.
func CreateDeployment(clientSet kubernetes.Clientset, namespace string, name string, imageName string, containerPorts []int32, podLabels map[string]string) error {
	deployment := configureDeployment(namespace, name, imageName, containerPorts, podLabels)
	deploymentsClient := clientSet.AppsV1().Deployments(namespace)
	result, err := deploymentsClient.Create(context.TODO(), &deployment, metav1.CreateOptions{})

	if err != nil {
		return err
	}

	utils.InfoLogger.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())
	return nil
}

func portArray(ports []int32) []apiv1.ContainerPort {
	result := make([]apiv1.ContainerPort, len(ports))
	for i := 0; i < len(ports); i++ {
		result[i] = apiv1.ContainerPort{
			Name:          fmt.Sprintf("port%d", i),
			Protocol:      apiv1.ProtocolTCP,
			ContainerPort: ports[i],
		}
	}
	return result
}

// configureDeployment makes a deployment configuration for a pod and replicaset
// -----
// NOTE ON PARAMETERS
// "name" and "imageName" are primarily separated for allowing multiple pods/deployments running the same image, as this is used in the tests.
func configureDeployment(nameSpace string, name string, imageName string, containerPorts []int32, podLabels map[string]string) appsv1.Deployment {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				"app": name,
			},
			Namespace: nameSpace, //NAMESPACE SKAL EKSISTERE OG MATCHE NAMESPACE I DEPLOYMENTKLIENTEN
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: utils.Int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": name,
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: podLabels,
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  name,
							Image: utils.ImageRepoUrl + imageName,
							Ports: portArray(containerPorts),
						},
					},
					ImagePullSecrets: []apiv1.LocalObjectReference{
						{
							Name: "regcred",
						},
					},
				},
			},
		},
	}
	return *deployment
}

// CheckIfDeploymentExists checks if a deployment exists in the given namespace.
func CheckIfDeploymentExists(clientSet kubernetes.Clientset, namespace string, deploymentName string) bool {
	deploymentsClient := clientSet.AppsV1().Deployments(namespace)
	list, err := deploymentsClient.List(context.TODO(), metav1.ListOptions{})
	utils.ErrLogger(err)
	for _, d := range list.Items {
		if d.Name == deploymentName {
			return true
		}
	}
	return false
}

// DeleteDeployment deletes a deployment in the given namespace.
func DeleteDeployment(clientSet kubernetes.Clientset, namespace string, deploymentName string) bool {
	utils.InfoLogger.Printf("Deleting deployment %s \n", deploymentName)
	deploymentsClient := clientSet.AppsV1().Deployments(namespace)
	deletePolicy := metav1.DeletePropagationForeground
	err := deploymentsClient.Delete(context.TODO(), deploymentName, metav1.DeleteOptions{PropagationPolicy: &deletePolicy})
	if err != nil {
		utils.ErrorLogger.Println(err)
		utils.InfoLogger.Println("Deployment could not be deleted")
		return false
	} else {
		utils.InfoLogger.Println("Deployment successfully deleted")
		return true
	}
}
