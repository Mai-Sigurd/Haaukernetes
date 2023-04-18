package deployments

import (
	"context"
	"fmt"
	utils "k8-project/utils"
	"log"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const imageRepoUrl = "registry.digitalocean.com/haaukins-bsc/"

func CreatePrebuiltDeployment(clientSet kubernetes.Clientset, teamName string, deployment *appsv1.Deployment) {
	log.Printf("Creating deployment %s\n", deployment.ObjectMeta.Name)
	deploymentsClient := clientSet.AppsV1().Deployments(teamName)
	result, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	utils.ErrHandler(err)
	log.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())
}

// CreateDeployment configures a deployment and then creates a deployment from that configuration
// in the given namespace.
func CreateDeployment(clientSet kubernetes.Clientset, teamName string, challengeName string, containerPorts []int32, podLabels map[string]string) {
	deployment := configureDeployment(teamName, challengeName, containerPorts, podLabels)
	fmt.Printf("Creating deployment %s\n", deployment.ObjectMeta.Name)
	deploymentsClient := clientSet.AppsV1().Deployments(teamName)
	result, err := deploymentsClient.Create(context.TODO(), &deployment, metav1.CreateOptions{})
	utils.ErrHandler(err)
	fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())
}

func portArray(ports []int32) []apiv1.ContainerPort {
	result := make([]apiv1.ContainerPort, len(ports))
	for i := 0; i < len(ports); i++ {
		result[i] = apiv1.ContainerPort{
			//Name:          "http",
			Name:          fmt.Sprintf("port%d", i), //not optimal but works
			Protocol:      apiv1.ProtocolTCP,
			ContainerPort: ports[i],
		}
	}
	return result
}

// configureDeployment makes a deployment configuration for a pod and replicaset
func configureDeployment(nameSpace string, name string, containerPorts []int32, podLabels map[string]string) appsv1.Deployment {
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
							Image: imageRepoUrl + name,
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

// PrintListDeployments ListDeployments lists the existing deployments in the given namespace.
// This also includes terminating deployments.
func PrintListDeployments(clientSet kubernetes.Clientset, namespace string) {
	list := GetAllDeployments(clientSet, namespace)
	fmt.Println("Listing deployments in default namespace")
	fmt.Printf("%d deployments exist\n", len(list.Items))
	for _, d := range list.Items {
		fmt.Printf(" * %s (%d replicas)\n", d.Name, *d.Spec.Replicas)
	}
}

func GetAllDeployments(clientSet kubernetes.Clientset, namespace string) *appsv1.DeploymentList {
	deploymentsClient := clientSet.AppsV1().Deployments(namespace)
	list, err := deploymentsClient.List(context.TODO(), metav1.ListOptions{})
	utils.ErrHandler(err)
	return list
}

// CheckIfDeploymentExists checks if a deployment exists in the given namespace.
func CheckIfDeploymentExists(clientSet kubernetes.Clientset, namespace string, deploymentName string) bool {
	deploymentsClient := clientSet.AppsV1().Deployments(namespace)
	list, err := deploymentsClient.List(context.TODO(), metav1.ListOptions{})
	utils.ErrHandler(err)
	for _, d := range list.Items {
		if d.Name == deploymentName {
			return true
		}
	}
	return false
}

// DeleteDeployment deletes a deployment in the given namespace.
func DeleteDeployment(clientSet kubernetes.Clientset, namespace string, deploymentName string) bool {
	fmt.Printf("Deleting deployment %s \n", deploymentName)
	deploymentsClient := clientSet.AppsV1().Deployments(namespace)
	deletePolicy := metav1.DeletePropagationForeground
	err := deploymentsClient.Delete(context.TODO(), deploymentName, metav1.DeleteOptions{PropagationPolicy: &deletePolicy})
	if err != nil {
		fmt.Println("Deployment could not be deleted")
		return false
	} else {
		fmt.Println("Deployment successfully deleted")
		return true
	}
}

//Functions for using local images while testing, as to not hit the docker pull limit of ~ 100 images / 6 hours

// CreateDeployment configures a deployment and then creates a deployment from that configuration
// in the given namespace, using a local docker image instead of an imagerepo
func CreateLocalDeployment(clientSet kubernetes.Clientset, teamName string, challengeName string, containerPorts []int32, podLabels map[string]string) {
	deployment := configureLocalDeployment(teamName, challengeName, containerPorts, podLabels)
	fmt.Printf("Creating deployment %s\n", deployment.ObjectMeta.Name)
	deploymentsClient := clientSet.AppsV1().Deployments(teamName)
	result, err := deploymentsClient.Create(context.TODO(), &deployment, metav1.CreateOptions{})
	utils.ErrHandler(err)
	fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())
}

// configureDeployment makes a deployment configuration for a pod and replicaset, but using a local docker image
func configureLocalDeployment(nameSpace string, name string, containerPorts []int32, podLabels map[string]string) appsv1.Deployment {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				"app": name,
			},
			Namespace: nameSpace,
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
							Name:            name,
							Image:           localImageRegistryUrl + name,
							ImagePullPolicy: apiv1.PullIfNotPresent, //PullNever also an option
							Ports:           portArray(containerPorts),
						},
					},
				},
			},
		},
	}
	return *deployment
}
