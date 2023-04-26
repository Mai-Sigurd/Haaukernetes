package tests

import (
	"fmt"
	"k8s-project/challenge"
	"k8s-project/kali"
	"k8s-project/namespaces"
	"k8s-project/utils"
	"k8s-project/wireguard"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var ports = map[string][]int32{"logon": {80}, "heartbleed": {443}, "for-fun-and-profit": {22}, "hide-and-seek": {13371}, "program-behaviour": {20, 21, 12020, 12021, 12022, 12023, 12024, 12025}, "reverseapk": {80}}

func getClientSet() *kubernetes.Clientset {
	home := homedir.HomeDir()
	kubeConfigPath := filepath.Join(home, ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	utils.ErrLogger(err)
	clientSet, err := kubernetes.NewForConfig(config)
	utils.ErrLogger(err)
	return clientSet
}

// The test uses a random public key
func setUpKubernetesResourcesWithWireguard(clientSet kubernetes.Clientset, namespace string, endpoint string, subnet string) {
	_ = namespaces.PostNamespace(clientSet, namespace)
	wireguard.StartWireguard(clientSet, namespace, "2A/Rj6X3+YxP6lXOv2BgbRQfpCn5z6Ob8scKhxiCRyM=", endpoint, subnet)
}
func setUpKubernetesResourcesWithKali(clientSet kubernetes.Clientset, namespace string) {
	_ = namespaces.PostNamespace(clientSet, namespace)
	kali.StartKaliImage(clientSet, namespace, "kali-firefox-test")
}

func startChallenge(name string, imageName string, clientSet kubernetes.Clientset, namespace string, challengePorts []int32) {
	challenge.CreateChallenge(clientSet, namespace, name, imageName, challengePorts)
}

func startAllChallenges(clientSet kubernetes.Clientset, namespace string) {
	log.Printf("Start 6 challenges")
	for key := range ports {
		startChallenge(key, key, clientSet, namespace, ports[key])
	}
}

func startAllChallengesWithDuplicates(clientSet kubernetes.Clientset, namespace string) {
	log.Printf("Starting 5x6 challenges")
	for key := range ports {
		for i := 1; i < 6; i++ {
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

// Logs

// TODO delete
func logCPUWithStoredResult(c chan string, results *string) {
	*results += "\n"
	input := ""
	go func() {
		input = <-c
	}()
	for input == "" {
		time.Sleep(500 * time.Millisecond)
		actualCPU, _ := cpu.Percent(500*time.Millisecond, false)
		usage := fmt.Sprintf("%s, %f\n", time.Now().Format("15:04:05"), actualCPU[0])
		*results += usage
	}
}
