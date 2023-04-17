package tests

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"k8-project/apis"
	"k8s.io/client-go/kubernetes"
	"log"
	"os"
	"time"
)

var ports = map[string][]int32{"logon": {80}, "heartbleed": {443}, "for-fun-and-profit": {22}, "hide-and-seek": {13371}, "program-behaviour": {20, 21, 12020, 12021, 12022, 12023, 12024, 12025}, "reverseapk": {80}}

// Kubernetes

func setUpKubernetesResourcesWithWireguard(clientSet kubernetes.Clientset, namespace string) {
	_ = apis.PostNamespaceKubernetes(clientSet, namespace)
	apis.StartWireguardKubernetes(clientSet, namespace, "2A/Rj6X3+YxP6lXOv2BgbRQfpCn5z6Ob8scKhxiCRyM=") //random publickey
}

func startChallenge(challengeNameI string, clientSet kubernetes.Clientset, namespace string, challengePorts []int32) {
	apis.PostChallengeKubernetes(clientSet, namespace, challengeNameI, challengePorts)
}

func startAllChallenges(clientSet kubernetes.Clientset, namespace string) {
	log.Printf("Start 6 challenges")
	for key := range ports {
		startChallenge(fmt.Sprintf(key+"%d", 1), clientSet, namespace, ports[key])
	}
}

func startAllChallengesWithDuplicates(clientSet kubernetes.Clientset, namespace string) {
	log.Printf("Starting 5x6 challenges")
	for key := range ports {
		for i := 1; i < 6; i++ {
			challengePorts := ports[key]
			startChallenge(fmt.Sprintf(key+"%d", i), clientSet, namespace, challengePorts)
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

func logCPUContiously(c chan string) {
	input := ""
	go func() {
		input = <-c
	}()
	for input == "" {
		//actualCPU := (1.0 - float64(cpuNow.Idle)/float64(cpuNow.Total)) * 100
		//actualCPU := float64(cpuNow.User) / float64(cpuNow.Total) * 100
		actualCPU, _ := cpu.Percent(500*time.Millisecond, false)
		thing := fmt.Sprintf("%f\n", actualCPU[0])
		log.Print(thing)
	}
}
