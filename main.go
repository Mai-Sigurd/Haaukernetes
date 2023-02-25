//https://medium.com/programming-kubernetes/building-stuff-with-the-kubernetes-api-part-4-using-go-b1d0e3c1c899
//https://github.com/kubernetes/client-go/blob/master/examples/create-update-delete-deployment/main.go

//HUSK AT MINIKUBE SKAL KØRE :)))

package main

import (
	"context"
	"fmt"
	"log"
	"path/filepath"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/apps/v1"

	//ovenstående er for at bringe v1.DeploymentInterface typen ind til brug som argument i func
	//-> var selv nødt til at finde den på docs, autoimport virkede ikke
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	//
	// Uncomment to load all auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth"
	//
	// Or uncomment to load specific auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
)

func main() {
	home := homedir.HomeDir()
	kubeconfig := filepath.Join(home, ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	err_handler(err)
	clientset, err := kubernetes.NewForConfig(config)
	err_handler(err)

	deploymentsClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)
	//deploymentsClient := clientset.AppsV1().Deployments("user-a")
	list_deployments(deploymentsClient)

	create_deployment(deploymentsClient, logon())

	//name := "user-a-logon"
	//delete_deployment(deploymentsClient, name)

	//list_deployments(deploymentsClient)

	create_namespace(*clientset, *namespace_test())

}

func namespace_test() *apiv1.Namespace {
	namespace := &apiv1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test",
		},
	}
	return namespace

}

func create_namespace(clientset kubernetes.Clientset, namespace apiv1.Namespace) {
	new_namespace, err := clientset.CoreV1().Namespaces().Create(context.TODO(), &namespace, metav1.CreateOptions{})
	err_handler(err)
	fmt.Printf("Created namespace with name %s\n", new_namespace.Name)
}

//TODO:
//- dette er blot en demo. Værdier skal parametriseres
func logon() appsv1.Deployment {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "user-a-logon",
			Labels: map[string]string{
				"app": "haaukins",
			},
			//Namespace: "user-a", //NAMESPACE SKAL EKSISTERE OG MATCHE NAMESPACE I DEPLOYMENTKLIENTEN
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "user-a-logon",
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "user-a-logon",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "logon",
							Image: "logon",
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 80,
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

func create_deployment(deploymentsClient v1.DeploymentInterface, deployment appsv1.Deployment) {
	fmt.Printf("Creating deployment %s\n", deployment.ObjectMeta.Name)
	result, err := deploymentsClient.Create(context.TODO(), &deployment, metav1.CreateOptions{})
	err_handler(err)

	fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())

}

func list_deployments(deploymentsClient v1.DeploymentInterface) {
	list, err := deploymentsClient.List(context.TODO(), metav1.ListOptions{})
	err_handler(err)

	fmt.Println("Listing deployments in default namespace")
	fmt.Printf("%d deployments exist", len(list.Items))
	for _, d := range list.Items {
		fmt.Printf(" * %s (%d replicas)\n", d.Name, *d.Spec.Replicas)
	}
}

func delete_deployment(deploymentsClient v1.DeploymentInterface, name string) {
	fmt.Printf("Ohai, deleting deployment %s \n", name)
	deletePolicy := metav1.DeletePropagationForeground
	if err := deploymentsClient.Delete(context.TODO(), name, metav1.DeleteOptions{PropagationPolicy: &deletePolicy}); err != nil {
		panic(err)
	}
	fmt.Println("Deployment deleted")
}

func err_handler(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

//spaghet fra eksempel-filen
func int32Ptr(i int32) *int32 { return &i }
