package services

import (
	"context"
	"fmt"
	"k8-project/utils"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
)

// CreateService creates an internal service in the given namespace.
func CreateService(clientSet kubernetes.Clientset, namespace string, exerciseName string, containerPort int32) *apiv1.Service {
	fmt.Println("Creating service client")
	serviceClient := clientSet.CoreV1().Services(namespace)

	service := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      exerciseName,
			Namespace: namespace,
			Labels: map[string]string{
				"app": exerciseName,
			},
		},
		Spec: apiv1.ServiceSpec{
			Ports: []apiv1.ServicePort{
				{
					Port:       containerPort,
					TargetPort: intstr.FromInt(int(containerPort)), // intstr.FromInt(32000)
				},
			},
			Selector: map[string]string{
				"app": exerciseName,
			},
			ClusterIP: "",
		},
	}
	createdService, err := serviceClient.Create(context.TODO(), service, metav1.CreateOptions{})
	utils.ErrHandler(err)
	fmt.Printf("Created service client with name %s\n", createdService.Name)
	return createdService
}

// CreateExposeService creates a service in the given namespace that will be exposed on a
// port assigned by the system.
func CreateExposeService(clientSet kubernetes.Clientset, nameSpace string, exerciseName string, containerPort int32) *apiv1.Service {
	fmt.Println("Creating expose service client")
	serviceClient := clientSet.CoreV1().Services(nameSpace)

	exposeService := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      exerciseName + "-expose",
			Namespace: nameSpace,
			Labels: map[string]string{
				"app": exerciseName,
			},
		},
		Spec: apiv1.ServiceSpec{
			Type: apiv1.ServiceTypeNodePort,
			Ports: []apiv1.ServicePort{
				{
					Port:       containerPort,
					Protocol:   apiv1.ProtocolTCP,
					TargetPort: intstr.FromInt(int(containerPort)),
				},
			},
			Selector: map[string]string{
				"app": exerciseName,
			},
		},
	}
	resultExposeService, err := serviceClient.Create(context.TODO(), exposeService, metav1.CreateOptions{})
	utils.ErrHandler(err)
	fmt.Printf("Created expose service client with name %s\n", resultExposeService.Name)
	return resultExposeService
}

// DeleteService deletes a service in the given namespace.
func DeleteService(clientSet kubernetes.Clientset, namespace string, serviceName string) bool {
	fmt.Printf("Deleting service %s \n", serviceName)
	serviceClient := clientSet.CoreV1().Services(namespace)
	deletePolicy := metav1.DeletePropagationForeground
	err := serviceClient.Delete(context.TODO(), serviceName, metav1.DeleteOptions{PropagationPolicy: &deletePolicy})
	if err != nil {
		fmt.Println("Service could not be deleted")
		return false
	} else {
		fmt.Println("Service successfully deleted")
		return true
	}
}
