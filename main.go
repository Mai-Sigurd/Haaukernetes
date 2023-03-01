//https://medium.com/programming-kubernetes/building-stuff-with-the-kubernetes-api-part-4-using-go-b1d0e3c1c899
//https://github.com/kubernetes/client-go/blob/master/examples/create-update-delete-deployment/main.go

//HUSK AT MINIKUBE SKAL KØRE :)))

//næste trin er at få exposed pod helt ud til browseren
//TechWorldWithNanas video der forklarer services er ret god. Særligt headless service kunne
//virke brugbart. Derudover, så kan man måske med en "ClusterIP" nøjes med 2x yaml fil pr. pod? (og ikke 3 som i mine eksempler)

package main

import (
	"bufio"
	"context"
	"fmt"
	"k8-project/deployments"
	utils "k8-project/utils"
	"os"
	"path/filepath"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	//ovenstående er for at bringe v1.DeploymentInterface typen ind til brug som argument i func
	//-> var selv nødt til at finde den på docs, autoimport virkede ikke
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	"k8-project/namespaces"
)

var exerciseToPorts = map[string]int32{"logon": 1, "heartbleed": 2}

func main() {
	home := homedir.HomeDir()
	kubeConfigPath := filepath.Join(home, ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	utils.ErrHandler(err)
	clientSet, err := kubernetes.NewForConfig(config)
	utils.ErrHandler(err)

	deploymentsClient := clientSet.AppsV1().Deployments(apiv1.NamespaceDefault)
	teamName := ""

	fmt.Println("--------------------")
	fmt.Println("Ohøj i skuret! Velkommen til haaukins")
	fmt.Println("--------------------")

	for {
		fmt.Println("Du har nu følgende valgmuligheder")
		fmt.Println("Skriv 'team' for at tilmelde dig")
		fmt.Println("Skriv 'exercise' to turn on the exercise")
		fmt.Println("Skriv 'kali' to launch VM with selected exercises via vnc")

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		input := scanner.Text()

		switch input {
		case "team":
			fmt.Println("Skriv dit team navn")
			scanner.Scan()
			teamName = scanner.Text()
			namespaces.CreateNamespace(*clientSet, teamName)
		case "exercise":
			if teamName == "" {
				fmt.Println("Please tilmeld dig først")
				time.Sleep(3 * time.Second)
			} else {
				fmt.Println("Skriv navnet on your exercise")
				scanner.Scan()
				exerciseName := scanner.Text()
				if port, ok := exerciseToPorts[exerciseName]; ok {
					deployment := deployments.ConfigureDeployment(teamName, exerciseName, port, exerciseName)
					deployments.CreateDeployment(deploymentsClient, deployment)
				} else {
					fmt.Println("Invalid exercise")
				}
				time.Sleep(3 * time.Second)
			}
		case "kali":
			fmt.Println("KALIIIII")
		default:
			fmt.Println("Invalid input")
		}
	}
}

// serviceport: https://stackoverflow.com/questions/74655705/how-to-create-a-service-port-in-client-go
func browser(clientset kubernetes.Clientset) {
	//create deployment
	deploymentsClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)
	deployments.CreateDeployment(deploymentsClient, logon_browser())

	//create service
	serviceClient := clientset.CoreV1().Services(apiv1.NamespaceDefault)
	service := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "haaukins",
			Namespace: "default",
			Labels: map[string]string{
				"app": "myapp",
			},
		},
		Spec: apiv1.ServiceSpec{
			Ports: []apiv1.ServicePort{
				{
					Port:       80,
					TargetPort: intstr.FromInt(32000),
				},
			},
			Selector: map[string]string{
				"app": "haaukins",
			},
			ClusterIP: "",
		},
	}
	serviceClient.Create(context.TODO(), service, metav1.CreateOptions{})

	//create nodeport (expose to outside world)
	expose := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "logon-expose",
			Namespace: "default",
			Labels: map[string]string{
				"app": "haaukins",
			},
		},
		Spec: apiv1.ServiceSpec{
			Type: apiv1.ServiceTypeNodePort,
			Ports: []apiv1.ServicePort{
				{
					NodePort:   32000,
					Port:       80,
					Protocol:   apiv1.ProtocolTCP,
					TargetPort: intstr.FromInt(80),
				},
			},
			Selector: map[string]string{
				"app": "haaukins",
			},
		},
	}
	serviceClient.Create(context.TODO(), expose, metav1.CreateOptions{})
}

// tilpasset version af logon()-funktion, til brug for browser-eksperimenter
func logon_browser() appsv1.Deployment {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "haaukins-deployment",
			Labels: map[string]string{
				"app": "haaukins",
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: utils.Int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "haaukins",
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "haaukins",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:            "logon",
							Image:           "logon",
							ImagePullPolicy: apiv1.PullNever,
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
