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
	"log"
	"os"
	"path/filepath"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
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

	browser(*clientset)
}

func gui_main() {
	home := homedir.HomeDir()
	kubeconfig := filepath.Join(home, ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	err_handler(err)
	clientset, err := kubernetes.NewForConfig(config)
	err_handler(err)

	deploymentsClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)
	_ = deploymentsClient //compiler

	fmt.Println(benis())

	for {
		fmt.Println("--------------------")
		fmt.Println("Ohøj i skuret! Velkommen til haaukins")
		fmt.Println("--------------------")
		fmt.Println("Du har nu følgende valgmuligheder")
		fmt.Println("Skriv 'l' for at se deployments")
		fmt.Println("Skriv 'c' for at oprette deployments")
		fmt.Println("Skriv 'd' for at slette deployments")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		input := scanner.Text()
		switch input {
		case "l":
			list_deployments(deploymentsClient)
		case "c":
			create_deployment(deploymentsClient, logon())
		case "d":
			delete_deployment(deploymentsClient, "user-a-logon")
		}
	}

}

func old_main() {
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

//baseret på mine indledende eksperimenter med at tilgå noget fra browser
//expose til browser demo
//serviceport: https://stackoverflow.com/questions/74655705/how-to-create-a-service-port-in-client-go
func browser(clientset kubernetes.Clientset) {
	//create deployment
	deploymentsClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)
	create_deployment(deploymentsClient, logon_browser())

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

//tilpasset version af logon()-funktion, til brug for browser-eksperimenter
func logon_browser() appsv1.Deployment {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "haaukins-deployment",
			Labels: map[string]string{
				"app": "haaukins",
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
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

//demo namespace
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

func create_deployment(deploymentsClient v1.DeploymentInterface, deployment appsv1.Deployment) {
	fmt.Printf("Creating deployment %s\n", deployment.ObjectMeta.Name)
	result, err := deploymentsClient.Create(context.TODO(), &deployment, metav1.CreateOptions{})
	err_handler(err)

	fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())

}

//obs: fordi vi lister deployments, så kan det umiddelbart ikke ses
//hvis den er slettet og pods derfor er ved at terminate...
func list_deployments(deploymentsClient v1.DeploymentInterface) {
	list, err := deploymentsClient.List(context.TODO(), metav1.ListOptions{})
	err_handler(err)

	fmt.Println("Listing deployments in default namespace")
	fmt.Printf("%d deployments exist\n", len(list.Items))
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

//https://patorjk.com/software/taag/#p=display&f=Basic&t=BENIS
func benis() string {
	s :=
		`
	_____________________ _______  .___  _________
	\______   \_   _____/ \      \ |   |/   _____/
	 |    |  _/|    __)_  /   |   \|   |\_____  \ 
	 |    |   \|        \/    |    \   |/        \
	 |______  /_______  /\____|__  /___/_______  /
	        \/        \/         \/            \/ 
	`
	return s
}

func benis2() string {

	s1 := ".----------------.  .----------------.  .-----------------. .----------------.  .----------------."
	s2 := "| .--------------. || .--------------. || .--------------. || .--------------. || .--------------. |"
	s3 := "| |   ______     | || |  _________   | || | ____  _____  | || |     _____    | || |    _______   | |"
	s4 := `| |  |_   _ \    | || | |_   ___  |  | || ||_   \|_   _| | || |    |_   _|   | || |   /  ___  |  | |`
	s5 := `| |    | |_) |   | || |   | |_  \_|  | || |  |   \ | |   | || |      | |     | || |  |  (__ \_|  | |`
	s6 := "| |    |  __'.   | || |   |  _|  _   | || |  | |\\ \\| |   | || |      | |     | || |   '.___`-.   | |"
	s7 := "| |   _| |__) |  | || |  _| |___/ |  | || | _| |_\\   |_  | || |     _| |_    | || |  |`\\____) |  | |"
	s8 := "| |  |_______/   | || | |_________|  | || ||_____|\\____| | || |    |_____|   | || |  |_______.'  | |"
	s9 := "| |              | || |              | || |              | || |              | || |              | |"
	s10 := "| '--------------' || '--------------' || '--------------' || '--------------' || '--------------' |"
	s11 := " '----------------'  '----------------'  '----------------'  '----------------'  '----------------' "

	return s1 + s2 + s3 + s4 + s5 + s6 + s7 + s8 + s9 + s10 + s11

}

// func benis2() string {
// 	s := fmt.Sprintf(
// 	`
// 	.----------------.  .----------------.  .-----------------. .----------------.  .----------------.
// 	| .--------------. || .--------------. || .--------------. || .--------------. || .--------------. |
// 	| |   ______     | || |  _________   | || | ____  _____  | || |     _____    | || |    _______   | |
// 	| |  |_   _ \    | || | |_   ___  |  | || ||_   \|_   _| | || |    |_   _|   | || |   /  ___  |  | |
// 	| |    | |_) |   | || |   | |_  \_|  | || |  |   \ | |   | || |      | |     | || |  |  (__ \_|  | |
// 	| |    |  __'.   | || |   |  _|  _   | || |  | |\ \| |   | || |      | |     | || |   '.___`-.   | |
// 	| |   _| |__) |  | || |  _| |___/ |  | || | _| |_\   |_  | || |     _| |_    | || |  |`\____) |  | |
// 	| |  |_______/   | || | |_________|  | || ||_____|\____| | || |    |_____|   | || |  |_______.'  | |
// 	| |              | || |              | || |              | || |              | || |              | |
// 	| '--------------' || '--------------' || '--------------' || '--------------' || '--------------' |
// 	 '----------------'  '----------------'  '----------------'  '----------------'  '----------------'
// 	`
// }

func haaukins() string {
	s :=
		`
	 ___ ___    _____      _____   ____ ___ ____  __.___ _______    _________ ________     _______   
	/   |   \  /  _  \    /  _  \ |    |   \    |/ _|   |\      \  /   _____/ \_____  \    \   _  \  
   /    ~    \/  /_\  \  /  /_\  \|    |   /      < |   |/   |   \ \_____  \   /  ____/    /  /_\  \ 
   \    Y    /    |    \/    |    \    |  /|    |  \|   /    |    \/        \ /       \    \  \_/   \
	\___|_  /\____|__  /\____|__  /______/ |____|__ \___\____|__  /_______  / \_______ \ /\ \_____  /
		  \/         \/         \/                 \/           \/        \/          \/ \/       \/`
	return s
}
