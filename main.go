//https://medium.com/programming-kubernetes/building-stuff-with-the-kubernetes-api-part-4-using-go-b1d0e3c1c899
//https://github.com/kubernetes/client-go/blob/master/examples/create-update-delete-deployment/main.go

// HUSK AT MINIKUBE SKAL KØRE :)))
package main

import (
	"bufio"
	"fmt"
	"k8-project/deployments"
	"k8-project/namespaces"
	"k8-project/netpol"
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

var challengeToPort = map[string]int32{"logon": 80, "heartbleed": 443, "for-fun-and-profit": 22, "always-be-crypting": 1337}

func main() {
	home := homedir.HomeDir()
	kubeConfigPath := filepath.Join(home, ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	utils.ErrHandler(err)
	clientSet, err := kubernetes.NewForConfig(config)
	utils.ErrHandler(err)

	fmt.Println("--------------------")
	fmt.Println("Hello from cyberspace! Welcome to haaukins")
	fmt.Println("--------------------")

	scanner := bufio.NewScanner(os.Stdin)
	teamName := ""
	for teamName == "" {
		fmt.Println("Write your alias")
		scanner.Scan()
		teamName = scanner.Text()
		if namespaces.NamespaceExists(*clientSet, teamName) {
			teamName = ""
		} else {
			namespaces.CreateNamespace(*clientSet, teamName)
			netpol.CreateKaliEgressPolicy(*clientSet, teamName)
			netpol.CreateChallengeIngressPolicy(*clientSet, teamName)
		}
	}

	for {
		fmt.Println("")
		fmt.Println("You have the following choices:")
		fmt.Println("Write 'on' to turn on a challenge")
		fmt.Println("Write 'off' to turn off a challenge")
		fmt.Println("Write 'kali' to launch VM with selected challenges via vnc")
		scanner.Scan()
		input := scanner.Text()

		switch input {
		case "on":
			fmt.Println("Write the name of the challenge to turn on")
			scanner.Scan()
			challengeName := scanner.Text()
			if port, ok := challengeToPort[challengeName]; ok {
				podLabels := make(map[string]string)
				podLabels["app"] = challengeName
				podLabels["type"] = "challenge"
				deployments.CreateDeployment(*clientSet, teamName, challengeName, port, podLabels)
				services.CreateService(*clientSet, teamName, challengeName, port)
			} else {
				fmt.Printf("Challenge %s does not exist", challengeName)
			}
		case "off":
			fmt.Println("Write the name of the challenge you want to turn off")
			scanner.Scan()
			challengeName := scanner.Text()
			if _, ok := challengeToPort[challengeName]; ok {
				deleteChallenge(*clientSet, teamName, challengeName)
			} else {
				fmt.Printf("Challenge %s does not exist", challengeName)
			}
		case "kali":
			startKali(*clientSet, teamName)
		default:
			fmt.Println("Invalid input")
		}
	}
}

func deleteChallenge(clientSet kubernetes.Clientset, teamName string, challengeName string) {
	if !deployments.CheckIfDeploymentExists(clientSet, teamName, challengeName) {
		fmt.Printf("Challenge %s is not turned on \n", challengeName)
	} else {
		deploymentDeleteStatus := deployments.DeleteDeployment(clientSet, teamName, challengeName)
		serviceDeleteStatus := services.DeleteService(clientSet, teamName, challengeName)
		if deploymentDeleteStatus && serviceDeleteStatus {
			fmt.Printf("Challenge %s turned off\n", challengeName)
		} else {
			fmt.Printf("Challenge %s could not be turned off\n", challengeName)
		}
	}
}

func startKali(clientSet kubernetes.Clientset, teamName string) {
	fmt.Println("Starting Kali")
	podLabels := make(map[string]string)
	podLabels["app"] = "kali-vnc"
	deployments.CreateDeployment(clientSet, teamName, "kali-vnc", 5901, podLabels)
	services.CreateService(clientSet, teamName, "kali-vnc", 5901)
	services.CreateExposeService(clientSet, teamName, "kali-vnc", 5901)
	fmt.Println("You can now vnc into your Kali. If on Mac first do `minikube service kali-vnc-expose -n <teamName>`")
	fmt.Println("If on Mac first do `minikube service kali-vnc-expose -n <teamName>` and use that url with vnc")
}
