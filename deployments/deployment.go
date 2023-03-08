package deployments

import (
	"context"
	"fmt"
	utils "k8-project/utils"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	//ovenstående er for at bringe v1.DeploymentInterface typen ind til brug som argument i func
	//-> var selv nødt til at finde den på docs, autoimport virkede ikke
	//
	// Uncomment to load all auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth"
	//
	// Or uncomment to load specific auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
)

func CreateKaliDeployment(clientSet kubernetes.Clientset, teamName string, podName string, containerPort int32) {
	deployment := configureKaliDeployment(teamName, podName, containerPort)
	fmt.Printf("Creating kali deployment %s\n", deployment.ObjectMeta.Name)
	deploymentsClient := clientSet.AppsV1().Deployments(teamName)
	result, err := deploymentsClient.Create(context.TODO(), &deployment, metav1.CreateOptions{})
	utils.ErrHandler(err)
	fmt.Printf("Created kali deployment %q.\n", result.GetObjectMeta().GetName())

}

// CreateExerciseDeployment configures a deployment and then creates a deployment from that configuration
// in the given namespace.
func CreateExerciseDeployment(clientSet kubernetes.Clientset, teamName string, exerciseName string, containerPort int32) {
	deployment := configureExerciseDeployment(teamName, exerciseName, containerPort)
	fmt.Printf("Creating deployment %s\n", deployment.ObjectMeta.Name)
	deploymentsClient := clientSet.AppsV1().Deployments(teamName)
	result, err := deploymentsClient.Create(context.TODO(), &deployment, metav1.CreateOptions{})
	utils.ErrHandler(err)
	fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())
}

func configureKaliDeployment(nameSpace string, name string, containerPort int32) appsv1.Deployment {
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
					Labels: map[string]string{
						"app": name,
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:            name,
							Image:           name,
							ImagePullPolicy: apiv1.PullNever,
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: containerPort,
								},
							},
						},
					},
				},
			},
		},
	}
	return *deployment
}

// configureExerciseDeployment makes a deployment configuration for a pod and replicaset
func configureExerciseDeployment(nameSpace string, name string, containerPort int32) appsv1.Deployment {
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
					Labels: map[string]string{
						"app":  name,
						"type": "exercise",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:            name,
							Image:           name,
							ImagePullPolicy: apiv1.PullNever,
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: containerPort,
								},
							},
						},
					},
				},
			},
		},
	}
	return *deployment
}

// ListDeployments lists the existing deployments in the given namespace.
// This also includes terminating deployments.
func ListDeployments(clientSet kubernetes.Clientset, namespace string) {
	deploymentsClient := clientSet.AppsV1().Deployments(namespace)
	list, err := deploymentsClient.List(context.TODO(), metav1.ListOptions{})
	utils.ErrHandler(err)

	fmt.Println("Listing deployments in default namespace")
	fmt.Printf("%d deployments exist\n", len(list.Items))
	for _, d := range list.Items {
		fmt.Printf(" * %s (%d replicas)\n", d.Name, *d.Spec.Replicas)
	}
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
