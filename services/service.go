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

// serviceport: https://stackoverflow.com/questions/74655705/how-to-create-a-service-port-in-client-go
func CreateServices(clientSet kubernetes.Clientset, nameSpace string, exerciseName string, containerPort int32) {
	fmt.Println("Creating service and expose service clients")
	serviceClient := clientSet.CoreV1().Services(nameSpace)

	service := getService(nameSpace, exerciseName, containerPort)
	resultS, err := serviceClient.Create(context.TODO(), service, metav1.CreateOptions{})
	utils.ErrHandler(err)
	fmt.Printf("Created service client with name %s\n", resultS.Name)

	exposeService := getExposeService(nameSpace, exerciseName, containerPort)
	resultE, err := serviceClient.Create(context.TODO(), exposeService, metav1.CreateOptions{})
	utils.ErrHandler(err)
	fmt.Printf("Created expose service client with name %s\n", resultE.Name)
}

func getService(nameSpace string, exerciseName string, containerPort int32) *apiv1.Service {
	service := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      exerciseName,
			Namespace: nameSpace,
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
				"app": nameSpace,
			},
			ClusterIP: "",
		},
	}
	return service
}

func getExposeService(nameSpace string, exerciseName string, containerPort int32) *apiv1.Service {
	//create nodeport (expose to outside world)
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
					//NodePort:   32000, // Den kan ikke være det samme for alle fordi den kan kun allokeres en gang
					Port:       containerPort,
					Protocol:   apiv1.ProtocolTCP,
					TargetPort: intstr.FromInt(int(containerPort)),
				},
			},
			Selector: map[string]string{
				"app": nameSpace,
			},
		},
	}
	return exposeService
}
