package tests

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/mem"
	"k8-project/apis"
	"k8-project/deployments"
	"k8-project/namespaces"
	"k8-project/netpol"
	"k8-project/secrets"
	"k8-project/services"
	"k8-project/utils"
	"k8-project/wireguard"
	"log"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"k8s.io/client-go/kubernetes"
)

func setupLog(filename string) *os.File {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	log.SetOutput(file)
	log.Printf("Testing started ")
	return file
}

func setUpKubernetesResources(clientSet kubernetes.Clientset, teamName string) {
	namespaces.CreateNamespace(clientSet, teamName)
	secrets.CreateImageRepositorySecret(clientSet, teamName)
	netpol.CreateChallengeIngressPolicy(clientSet, teamName)
	netpol.CreateEgressPolicy(clientSet, teamName)
	wireguard.StartWireguard(clientSet, teamName, "2A/Rj6X3+YxP6lXOv2BgbRQfpCn5z6Ob8scKhxiCRyM=") //random publickey
	netpol.AddWireguardToChallengeIngressPolicy(clientSet, teamName)
}

func startChallenge(challengeNameI string, clientSet kubernetes.Clientset, namespace string, challengePorts []int32) {
	podLabels := make(map[string]string)
	podLabels["app"] = challengeNameI
	podLabels["type"] = "challenge"
	deployments.CreateLocalDeployment(clientSet, namespace, challengeNameI, challengePorts, podLabels)
	services.CreateService(clientSet, namespace, challengeNameI, challengePorts)
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

func logCPUWithStoredResult(c chan string, results *string) {
	result := "\n"
	input := ""
	go func() {
		input = <-c
	}()
	for input == "" {
		time.Sleep(500 * time.Millisecond)
		//actualCPU := (1.0 - float64(cpuNow.Idle)/float64(cpuNow.Total)) * 100
		//actualCPU := float64(cpuNow.User) / float64(cpuNow.Total) * 100
		actualCPU, _ := cpu.Percent(500*time.Millisecond, false)
		thing := fmt.Sprintf("%s, %f\n", time.Now().Format("15:04:05"), actualCPU[0])
		result = result + thing
		*results = result
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
		log.Printf(thing)
	}
}

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

	setUpKubernetesResources(*clientSet, personA)

	startAllChallenges(*clientSet, personA)
	apis.StartKali(*clientSet, personA)

	time.Sleep(10 * time.Second)
	comChannel <- "stop"
	log.Println(results)

	namespaces.DeleteNamespace(*clientSet, personA)

}

// TODO delete
func TestMai(t *testing.T) {
	file := setupLog("Mai")
	defer file.Close()
	//
	comChannel := make(chan string)

	go logCPUContiously(comChannel)
	log.Printf("something in the middle happens")
	time.Sleep(5 * time.Second)
	comChannel <- "stop"
	log.Printf("midway")
	go logCPUContiously(comChannel)
	time.Sleep(5 * time.Second)
	comChannel <- "stop"
}

// Find out how many users there can be run on a minimal kubernetes requirements server setup (with an amount of challenges running) while we wait in between the starting of namespaces
func TestMinimalKubernetesSetup(t *testing.T) {
	file := setupLog("Minimal-k8s")
	defer file.Close()
	//
}

// Find out how many users there can be run on a minimal kubernetes requirements, stress testing how many namespaces can start at the same time.
func TestMinimalKubernetesSetupStartup(t *testing.T) {
	file := setupLog("Minimal-k8s")
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
		setUpKubernetesResources(*clientSet, personI)
		startAllChallenges(*clientSet, personI)
	}
	time.Sleep(5 * time.Second)
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

// TODO: add memory also
// Research usage of different amount of open challenges, like max 5 vs. all challenges running
func TestChallengeLoad(t *testing.T) {
	file := setupLog("challenge-load")
	defer file.Close()

	//does not work -> maybe calling 'go setupLog' might keep the file open?
	//setupLog("Challenge-load")
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
	var sixchallengeslog string
	go logCPUWithStoredResult(comChannel, &sixchallengeslog)

	setUpKubernetesResources(*clientSet, teamName)
	startAllChallenges(*clientSet, teamName)

	time.Sleep(30 * time.Second)

	comChannel <- "stop"
	log.Printf("Challenges has been run \n")

	namespaces.DeleteNamespace(*clientSet, teamName)

	log.Printf("Sleeping 90 seconds to allow for namespace to be deleted, loggin cpu contiosuly \n")

	time.Sleep(90 * time.Second)

	log.Printf("CPU 90 seconds after running 6 challenges and deleting namespace, before running with 30 \n")

	//30 challenges
	//-> deployment/pod/service name og imagename skal adskilles i param...
	log.Printf("Setting up namespace etc. for 30 challenges, loggin cpu contiously\n")
	var thirtyChallengeslog string
	go logCPUWithStoredResult(comChannel, &thirtyChallengeslog)
	setUpKubernetesResources(*clientSet, teamName)
	startAllChallengesWithDuplicates(*clientSet, teamName)

	time.Sleep(30 * time.Second)
	comChannel <- "stop"
	log.Printf("Done running 30 challengesl\n")
	log.Printf("Results for running 6 challenges")
	log.Println(sixchallengeslog)
	log.Printf("Results for running 30 challenges")
	log.Println(thirtyChallengeslog)

	namespaces.DeleteNamespace(*clientSet, teamName)
	time.Sleep(30 * time.Second)

}
