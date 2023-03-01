package deployments

import (
	"context"
	"fmt"
	utils "k8-project/utils"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/apps/v1"
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

func CreateDeployment(clientSet kubernetes.Clientset, teamName string, exerciseName string, containerPort int32) {
	deployment := configureDeployment(teamName, exerciseName, containerPort)
	fmt.Printf("Creating deployment %s\n", deployment.ObjectMeta.Name)
	deploymentsClient := clientSet.AppsV1().Deployments(teamName)
	result, err := deploymentsClient.Create(context.TODO(), &deployment, metav1.CreateOptions{})
	utils.ErrHandler(err)
	fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())
}

// name = "logon"
// containerPort = 80
// appLabel = "haaukins"
// nameSpace = "user-a"
func configureDeployment(nameSpace string, name string, containerPort int32) appsv1.Deployment {
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

// Deployment er en pod
// obs: fordi vi lister deployments, så kan det umiddelbart ikke ses
// hvis den er slettet og pods derfor er ved at terminate...
func ListDeployments(deploymentsClient v1.DeploymentInterface) {
	list, err := deploymentsClient.List(context.TODO(), metav1.ListOptions{})
	utils.ErrHandler(err)

	fmt.Println("Listing deployments in default namespace")
	fmt.Printf("%d deployments exist\n", len(list.Items))
	for _, d := range list.Items {
		fmt.Printf(" * %s (%d replicas)\n", d.Name, *d.Spec.Replicas)
	}
}

func DeleteDeployment(deploymentsClient v1.DeploymentInterface, name string) {
	fmt.Printf("Ohai, deleting deployment %s \n", name)
	deletePolicy := metav1.DeletePropagationForeground
	if err := deploymentsClient.Delete(context.TODO(), name, metav1.DeleteOptions{PropagationPolicy: &deletePolicy}); err != nil {
		panic(err)
	}
	fmt.Println("Deployment deleted")
}
