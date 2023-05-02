package services

import (
	"context"
	"fmt"
	"k8s-project/utils"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
)

func CreatePrebuiltService(clientSet kubernetes.Clientset, teamName string, service apiv1.Service) (*apiv1.Service, error) {
	serviceClient := clientSet.CoreV1().Services(teamName)
	result, err := serviceClient.Create(context.TODO(), &service, metav1.CreateOptions{})

	if err != nil {
		return nil, err
	}

	utils.InfoLogger.Printf("Created service %q.\n", result.GetObjectMeta().GetName())
	return result, nil
}

func portArray(ports []int32) []apiv1.ServicePort {
	result := make([]apiv1.ServicePort, len(ports))
	for i := 0; i < len(ports); i++ {
		result[i] = apiv1.ServicePort{
			Name:       fmt.Sprintf("port%d", i),
			Port:       ports[i],
			TargetPort: intstr.FromInt(int(ports[i])), // intstr.FromInt(32000)
		}
	}
	return result
}

// CreateService creates an internal service in the given namespace.
func CreateService(clientSet kubernetes.Clientset, namespace string, challengeName string, containerPorts []int32) (*apiv1.Service, error) {
	serviceClient := clientSet.CoreV1().Services(namespace)

	service := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      challengeName,
			Namespace: namespace,
			Labels: map[string]string{
				"app": challengeName,
			},
		},
		Spec: apiv1.ServiceSpec{
			Ports: portArray(containerPorts),
			Selector: map[string]string{
				"app": challengeName,
			},
			ClusterIP: "",
		},
	}
	createdService, err := serviceClient.Create(context.TODO(), service, metav1.CreateOptions{})

	if err != nil {
		return nil, err
	}

	utils.InfoLogger.Printf("Created service client with name %s\n", createdService.Name)
	return createdService, nil
}

// DeleteService deletes a service in the given namespace.
func DeleteService(clientSet kubernetes.Clientset, namespace string, serviceName string) bool {
	utils.InfoLogger.Printf("Deleting service %s \n", serviceName)
	serviceClient := clientSet.CoreV1().Services(namespace)
	deletePolicy := metav1.DeletePropagationForeground
	err := serviceClient.Delete(context.TODO(), serviceName, metav1.DeleteOptions{PropagationPolicy: &deletePolicy})
	if err != nil {
		utils.ErrLogger(err)
		utils.InfoLogger.Println("Service could not be deleted")
		return false
	} else {
		utils.InfoLogger.Println("Service successfully deleted")
		return true
	}
}
