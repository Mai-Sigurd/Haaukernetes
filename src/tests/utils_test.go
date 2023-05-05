package tests

import (
	"fmt"
	"k8s-project/challenge"
	"k8s-project/connections/browser/kali"
	"k8s-project/connections/vpn/wireguard"
	"k8s-project/namespaces"
	"k8s-project/utils"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"strings"
)

var ports = map[string][]int32{"heartbleed": {443}, "for-fun-and-profit": {22}, "hide-and-seek": {13371}, "program-behaviour": {20, 21, 12020, 12021, 12022, 12023, 12024, 12025}, "reverseapk": {80}}

func getClientSet() *kubernetes.Clientset {
	kubeConfigPath := os.Getenv("KUBECONFIG") //running without docker requires 'export KUBECONFIG="$HOME/.kube/config"'
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	utils.ErrLogger(err)
	clientSet, err := kubernetes.NewForConfig(config)
	utils.ErrLogger(err)
	return clientSet
}

// The test uses a random public key
func setUpKubernetesResourcesWithWireguard(clientSet kubernetes.Clientset, namespace string, endpoint string, subnet string) error {
	err := namespaces.PostNamespace(clientSet, namespace)
	utils.TestLogger.Println("Starting wireguard")
	if err != nil {
		return err
	}
	_, err = wireguard.StartWireguard(clientSet, namespace, "2A/Rj6X3+YxP6lXOv2BgbRQfpCn5z6Ob8scKhxiCRyM=", endpoint, subnet)
	if err != nil {
		return err
	}
	return nil
}

func setUpKubernetesResourcesWithKali(clientSet kubernetes.Clientset, namespace string) error {
	err := namespaces.PostNamespace(clientSet, namespace)
	if err != nil {
		return err
	}
	_, _, err = kali.StartKali(clientSet, namespace, "kali-firefox-test")
	if err != nil {
		return err
	}
	return nil
}

func setUpKubernetesResourcesWithWireguardAndChannel(clientSet kubernetes.Clientset, namespace string, endpoint string, subnet string, channel chan string) error {
	err := namespaces.PostNamespace(clientSet, namespace)
	if err != nil {
		channel <- "error"
		return err
	}
	_, err = wireguard.StartWireguard(clientSet, namespace, "2A/Rj6X3+YxP6lXOv2BgbRQfpCn5z6Ob8scKhxiCRyM=", endpoint, subnet)
	if err != nil {
		channel <- "error"
		return err
	}
	return nil
}

func setUpKubernetesResourcesWithKaliAndChannel(clientSet kubernetes.Clientset, namespace string, channel chan string) error {
	err := namespaces.PostNamespace(clientSet, namespace)
	if err != nil {
		channel <- err.Error()
		return err
	}
	_, _, err = kali.StartKali(clientSet, namespace, "kali-firefox-test")
	if err != nil {
		channel <- err.Error()
		return err
	}
	return nil
}

func startChallenge(name string, imageName string, clientSet kubernetes.Clientset, namespace string, challengePorts []int32) error {
	err := challenge.CreateChallenge(clientSet, namespace, name, imageName, challengePorts)
	if err != nil {
		return err
	}
	return nil
}

func startAllChallenges(clientSet kubernetes.Clientset, namespace string) error {
	utils.TestLogger.Printf("Start 5 challenges")
	for key := range ports {
		err := startChallenge(key, key, clientSet, namespace, ports[key])
		if err != nil {
			return err
		}
	}
	return nil
}

func startAllChallengesWithDuplicates(clientSet kubernetes.Clientset, namespace string) {
	utils.TestLogger.Printf("Starting 5x5 challenges")
	for key := range ports {
		for i := 1; i < 7; i++ {
			challengePorts := ports[key]
			startChallenge(fmt.Sprintf(key+"%d", i), key, clientSet, namespace, challengePorts)
		}
	}
}

func findPodIp(pods *v1.PodList) string {
	for i := range pods.Items {
		if strings.Contains(pods.Items[i].Name, "logon") {
			return pods.Items[i].Status.PodIP
		}
	}
	return "IP of wireguard pod not found"
}
