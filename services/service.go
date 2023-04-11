package services

import (
	"context"
	"k8-project/utils"
	"log"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
)

func CreatePrebuiltService(clientSet kubernetes.Clientset, teamName string, service apiv1.Service) *apiv1.Service {
	log.Printf("Creating service %s\n", service.ObjectMeta.Name)
	serviceClient := clientSet.CoreV1().Services(teamName)
	result, err := serviceClient.Create(context.TODO(), &service, metav1.CreateOptions{})
	utils.ErrHandler(err)
	log.Printf("Created service %q.\n", result.GetObjectMeta().GetName())
	return result
}

func portArray(ports []int32) []apiv1.ServicePort {
	var result []apiv1.ServicePort
	for i := 0; i < len(ports); i++ {
		result[i] = apiv1.ServicePort{
			Port:       ports[i],
			TargetPort: intstr.FromInt(int(ports[i])), // intstr.FromInt(32000)
		}
	}
	return result
}

// CreateService creates an internal service in the given namespace.
func CreateService(clientSet kubernetes.Clientset, namespace string, challengeName string, containerPorts []int32) *apiv1.Service {
	log.Println("Creating service client")
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
	utils.ErrHandler(err)
	log.Printf("Created service client with name %s\n", createdService.Name)
	return createdService
}

// CreateExposeService creates a service in the given namespace that will be exposed on a
// port assigned by the system.
func CreateExposeService(clientSet kubernetes.Clientset, nameSpace string, challengeName string, containerPort int32) *apiv1.Service {
	log.Println("Creating expose service client")
	serviceClient := clientSet.CoreV1().Services(nameSpace)

	exposeService := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      challengeName + "-expose",
			Namespace: nameSpace,
			Labels: map[string]string{
				"app": challengeName,
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
				"app": challengeName,
			},
		},
	}
	resultExposeService, err := serviceClient.Create(context.TODO(), exposeService, metav1.CreateOptions{})
	utils.ErrHandler(err)
	log.Printf("Created expose service client with name %s\n", resultExposeService.Name)
	return resultExposeService
}

// DeleteService deletes a service in the given namespace.
func DeleteService(clientSet kubernetes.Clientset, namespace string, serviceName string) bool {
	log.Printf("Deleting service %s \n", serviceName)
	serviceClient := clientSet.CoreV1().Services(namespace)
	deletePolicy := metav1.DeletePropagationForeground
	err := serviceClient.Delete(context.TODO(), serviceName, metav1.DeleteOptions{PropagationPolicy: &deletePolicy})
	if err != nil {
		log.Println("Service could not be deleted")
		return false
	} else {
		log.Println("Service successfully deleted")
		return true
	}
}
