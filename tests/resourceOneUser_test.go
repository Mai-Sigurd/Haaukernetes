package tests

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"k8-project/namespaces"
	"k8-project/utils"
	"log"
	"os/exec"
	"testing"
	"time"
)

// General load (resources used for new user, kali docker(simple vs kali many tools), wireguard, guacamole, etc)
func Test6ChallengesAndWireguardOneUser(t *testing.T) {
	file := setupLog("6ChallengesAndWireguardOneUser")
	//TODO should also log memory
	defer file.Close()
	// CPU Load before starting
	comChannel := make(chan string)
	var results string
	go logCPUWithStoredResult(comChannel, &results)

	// Starting the kuberneets
	clientSet := getClientSet()
	personA := "persona"

	settings := utils.ReadYaml("settings-test.yaml")
	setUpKubernetesResourcesWithWireguard(*clientSet, personA, settings.Endpoint, settings.Subnet)

	startAllChallenges(*clientSet, personA)
	time.Sleep(10 * time.Second)
	comChannel <- "stop"
	log.Println(results)

	namespaces.DeleteNamespace(*clientSet, personA)
}

func Test6ChallengesAndKaliOneUser(t *testing.T) {
	file := setupLog("6ChallengesAndKaliOneUser")
	//TODO do it
	defer file.Close()
}

// Research usage of different amount of open challenges, like max 5 vs. all challenges running
func Test30ChallengesWireguardOneUser(t *testing.T) {
	file := setupLog("challenge-load")
	defer file.Close()

	// TODO remove run 6 challenges
	v, _ := mem.VirtualMemory()
	// almost every return value is a struct
	//TODO incooperate mem in function
	fmt.Printf("Total: %v, Free:%v, UsedPercent:%f%%\n", v.Total, v.Free, v.UsedPercent)

	currentTime := time.Now()
	log.Printf("Testing started " + currentTime.Format("15:04:05"))

	clientSet := getClientSet()
	teamName := "test"

	//6 challenges

	cpuBeforeFew, err := cpu.Percent(500*time.Millisecond, false)
	utils.ErrHandler(err)
	log.Printf("CPU before running 6 challenges %f\n", cpuBeforeFew[0])

	log.Printf("Setting up namespace etc. for 6 challenges, loggin cpu contiously\n ")

	comChannel := make(chan string)
	memChannel := make(chan string) //might be able to use same channel but with pointer?
	var sixchallengeslog string
	var sixChallengesMemLog string
	go logCPUWithStoredResult(comChannel, &sixchallengeslog)
	go logMemoryWithStoredResult(memChannel, &sixChallengesMemLog)

	settings := utils.ReadYaml("settings-test.yaml")
	setUpKubernetesResourcesWithWireguard(*clientSet, teamName, settings.Endpoint, settings.Subnet)
	startAllChallenges(*clientSet, teamName)

	time.Sleep(30 * time.Second)

	comChannel <- "stop"
	memChannel <- "stop"
	log.Printf("Challenges has been run \n")

	namespaces.DeleteNamespace(*clientSet, teamName)

	//watch for namespace to be removed definetely
	res, _ := exec.Command("/bin/sh", "-c", fmt.Sprintf("watch -e \"kubectl get ns | grep -m 1 %s\" ", teamName)).Output()
	_ = res //ignore

	//30 challenges
	//-> deployment/pod/service name og imagename skal adskilles i param...
	log.Printf("Setting up namespace etc. for 30 challenges, loggin cpu contiously\n")
	var thirtyChallengeslog string
	var thirtyChallengesMemLog string
	go logCPUWithStoredResult(comChannel, &thirtyChallengeslog)
	go logMemoryWithStoredResult(memChannel, &thirtyChallengesMemLog)
	setUpKubernetesResourcesWithWireguard(*clientSet, teamName, settings.Endpoint, settings.Subnet)
	startAllChallengesWithDuplicates(*clientSet, teamName)

	time.Sleep(30 * time.Second)
	comChannel <- "stop"
	memChannel <- "stop"
	log.Printf("Done running 30 challenges\n")
	log.Printf("Results for running 6 challenges")
	log.Println(sixchallengeslog)
	log.Println("MEMORY")
	log.Println(sixChallengesMemLog)
	log.Printf("Results for running 30 challenges")
	log.Println(thirtyChallengeslog)
	log.Println("MEMORY")
	log.Println(thirtyChallengesMemLog)

	namespaces.DeleteNamespace(*clientSet, teamName)
	time.Sleep(30 * time.Second)
}

func Test30ChallengesKaliOneUser(t *testing.T) {
	//TODO
	file := setupLog("challenge-load")
	defer file.Close()
}
