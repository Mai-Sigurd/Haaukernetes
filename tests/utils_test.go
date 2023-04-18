package tests

import (
	"fmt"
	"k8-project/api_endpoints"
	"log"
	"os"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"k8s.io/client-go/kubernetes"
)

var ports = map[string][]int32{"logon": {80}, "heartbleed": {443}, "for-fun-and-profit": {22}, "hide-and-seek": {13371}, "program-behaviour": {20, 21, 12020, 12021, 12022, 12023, 12024, 12025}, "reverseapk": {80}}

// Kubernetes

func setUpKubernetesResourcesWithWireguard(clientSet kubernetes.Clientset, namespace string) {
	_ = api_endpoints.PostNamespaceKubernetes(clientSet, namespace)
	api_endpoints.StartWireguardKubernetes(clientSet, namespace, "2A/Rj6X3+YxP6lXOv2BgbRQfpCn5z6Ob8scKhxiCRyM=") //random publickey
}

func startChallenge(name string, imageName string, clientSet kubernetes.Clientset, namespace string, challengePorts []int32) {
	api_endpoints.PostChallengeKubernetes(clientSet, namespace, name, imageName, challengePorts)
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

// Logs
func setupLog(filename string) *os.File {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	log.SetOutput(file)
	log.Printf("Testing started ")
	return file
}

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
func logMemoryWithStoredResult(c chan string, results *string) {
	input := ""
	go func() {
		input = <-c
	}()
	for input == "" {
		time.Sleep(500 * time.Millisecond)
		memory, _ := mem.VirtualMemory()
		usage := fmt.Sprintf("Total: %v, Free:%v, UsedPercent:%f%%\n", memory.Total, memory.Free, memory.UsedPercent)
		*results += usage
	}
}

// TODO use or delete
//func logCPUContiously(c chan string) {
//	input := ""
//	go func() {
//		input = <-c
//	}()
//	for input == "" {
//		actualCPU, _ := cpu.Percent(500*time.Millisecond, false)
//		thing := fmt.Sprintf("%f\n", actualCPU[0])
//		log.Print(thing)
//	}
//}
