//https://medium.com/programming-kubernetes/building-stuff-with-the-kubernetes-api-part-4-using-go-b1d0e3c1c899
//https://github.com/kubernetes/client-go/blob/master/examples/create-update-delete-deployment/main.go

//HUSK AT MINIKUBE SKAL KØRE :)))

//næste trin er at få exposed pod helt ud til browseren
//TechWorldWithNanas video der forklarer services er ret god. Særligt headless service kunne
//virke brugbart. Derudover, så kan man måske med en "ClusterIP" nøjes med 2x yaml fil pr. pod? (og ikke 3 som i mine eksempler)

package main

import (
	"bufio"
	"fmt"
	"k8-project/deployments"
	"k8-project/namespaces"
	"k8-project/services"
	"k8-project/utils"
	"k8s.io/client-go/kubernetes"
	"os"
	"path/filepath"
	//ovenstående er for at bringe v1.DeploymentInterface typen ind til brug som argument i func
	//-> var selv nødt til at finde den på docs, autoimport virkede ikke
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var exerciseToPorts = map[string]int32{"logon": 80, "heartbleed": 443}

func main() {
	home := homedir.HomeDir()
	kubeConfigPath := filepath.Join(home, ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	utils.ErrHandler(err)
	clientSet, err := kubernetes.NewForConfig(config)
	utils.ErrHandler(err)

	fmt.Println("--------------------")
	fmt.Println("Ohøj i skuret! Velkommen til haaukins")
	fmt.Println("--------------------")

	scanner := bufio.NewScanner(os.Stdin)
	teamName := ""
	for teamName == "" {
		fmt.Println("Skriv dit alias")
		scanner.Scan()
		teamName = scanner.Text()
		namespaces.CreateNamespace(*clientSet, teamName)
	}

	for {
		fmt.Println("")
		fmt.Println("Du har nu følgende valgmuligheder")
		fmt.Println("Skriv 'exercise' to turn on an exercise")
		fmt.Println("Skriv 'kali' to launch VM with selected exercises via vnc")
		scanner.Scan()
		input := scanner.Text()

		switch input {
		case "exercise":
			fmt.Println("Skriv navnet on the exercise to turn on")
			scanner.Scan()
			exerciseName := scanner.Text()
			if port, ok := exerciseToPorts[exerciseName]; ok {
				deployments.CreateDeployment(*clientSet, teamName, exerciseName, port)
				services.CreateService(*clientSet, teamName, exerciseName, port)
			} else {
				fmt.Println("Invalid exercise")
			}
		case "kali":
			fmt.Println("KALIIIII")
			deployments.CreateDeployment(*clientSet, teamName, "kali-vnc", 5901)
			services.CreateService(*clientSet, teamName, "kali-vnc", 5901)
			services.CreateExposeService(*clientSet, teamName, "kali-vnc", 5901)
		default:
			fmt.Println("Invalid input")
		}
	}
}
