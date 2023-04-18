package tests

import (
	"fmt"
	"k8-project/apis"
	"k8-project/namespaces"
	"k8-project/utils"
	"log"
	"os/exec"
	"strconv"
	"testing"
	"time"

	"github.com/shirou/gopsutil/v3/mem"

	"github.com/shirou/gopsutil/v3/cpu"
)

// General load (resources used for new user, kali docker(simple vs kali many tools), wireguard, guacamole, etc)
func TestGeneralLoad(t *testing.T) {
	file := setupLog("General-load")
	defer file.Close()
	// CPU Load before starting
	comChannel := make(chan string)
	var results string
	go logCPUWithStoredResult(comChannel, &results)

	// Starting the kuberneets
	clientSet := getClientSet()
	personA := "persona"

	setUpKubernetesResourcesWithWireguard(*clientSet, personA)

	startAllChallenges(*clientSet, personA)
	apis.PostKaliKubernetes(*clientSet, personA)

	time.Sleep(10 * time.Second)
	comChannel <- "stop"
	log.Println(results)

	namespaces.DeleteNamespace(*clientSet, personA)

}

// Find out how many users there can be run on a minimal kubernetes requirements server setup (with an amount of challenges running) while we wait in between the starting of namespaces
// TODO mememory might be relevant
func TestMinimalKubernetesSetup(t *testing.T) {
	file := setupLog("Minimal-k8s-den-anden")
	defer file.Close()
	//
}

// Find out how many users there can be run on a minimal kubernetes requirements, stress testing how many namespaces can start at the same time.
// TODO mememory might be relevant
func TestMinimalKubernetesSetupStartup(t *testing.T) {
	file := setupLog("Minimal-k8s-den-ene")
	defer file.Close()
	//
}

// Find out how much resource usage there is for decently size competition (maybe the amount of people of who participate in cybermesterskaberne).
func TestChampionshipLoad(t *testing.T) {

	file := setupLog("Championship")
	defer file.Close()
	clientSet := getClientSet()

	comChannel := make(chan string)
	var results string
	go logCPUWithStoredResult(comChannel, &results)

	const amountOfPeople = 350
	people := [amountOfPeople]string{}

	for i := 0; i < amountOfPeople; i++ {
		is := strconv.Itoa(i)
		personI := "person" + is
		people[i] = personI
		setUpKubernetesResourcesWithWireguard(*clientSet, personI)
		startAllChallenges(*clientSet, personI)
	}
	time.Sleep(30 * time.Second)
	comChannel <- "stop"
	log.Println(results)

	for i := 0; i < amountOfPeople; i++ {
		is := strconv.Itoa(i)
		personI := "person" + is
		people[i] = personI
		namespaces.DeleteNamespace(*clientSet, personI)
	}
	time.Sleep(5 * time.Second)

}

// Research usage of different amount of open challenges, like max 5 vs. all challenges running
func TestChallengeLoad(t *testing.T) {
	file := setupLog("challenge-load")
	defer file.Close()

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

	setUpKubernetesResourcesWithWireguard(*clientSet, teamName)
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
	setUpKubernetesResourcesWithWireguard(*clientSet, teamName)
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
