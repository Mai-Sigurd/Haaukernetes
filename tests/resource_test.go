package tests

import (
	"fmt"
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

	"k8s.io/client-go/kubernetes"

	"github.com/mackerelio/go-osstat/cpu"
)

// TODO: defer works or not?!
func setupLog(filename string) *os.File {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	log.SetOutput(file)
	currentTime := time.Now()
	log.Printf("Testing started " + currentTime.Format("2006.01.02 15:04:05"))
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
	log.Printf("Start all challenges only starts 6")
	for key := range ports {
		startChallenge(key, clientSet, namespace, ports[key])
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

func logCPU(c chan string, results *[]string) {
	var result []string
	input := ""
	go func() {
		input = <-c
	}()
	for input == "" {
		time.Sleep(500 * time.Millisecond)
		cpuNow, err := cpu.Get()
		if err != nil {
			log.Fatalf("%s\n", err)
		}
		thing := fmt.Sprintf("%s, %f", time.Now().Format("2006.01.02 15:04:0"), float64(cpuNow.System))
		result = append(result, thing)
		results = &result
	}
}

// General load (resources used for new user, kali docker(simple vs kali many tools), wireguard, guacamole, etc)
func TestGeneralLoad(t *testing.T) {
	file := setupLog("General-load")
	defer file.Close()
	// CPU Load before starting
	comChannel := make(chan string)
	var results []string
	go logCPU(comChannel, &results)

	// Starting the kuberneets
	clientSet := getClientSet()
	personA := "persona"

	setUpKubernetesResources(*clientSet, personA)

	startAllChallenges(*clientSet, personA)
	apis.StartKali(*clientSet, personA)

	time.Sleep(10 * time.Second)
	comChannel <- "stop"
	log.Println(results)

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

	const amountOfPeople = 350
	people := [amountOfPeople]string{}

	for i := 0; i < amountOfPeople; i++ {
		is := strconv.Itoa(i)
		personI := "person" + is
		people[i] = personI
		setUpKubernetesResources(*clientSet, personI)
		startAllChallenges(*clientSet, personI)
	}

	// CPU Load after starting
	after, err := cpu.Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return
	}
	// TODO igen hmm, er det vi gerne vil have
	log.Printf("cpu system after test: %f %%\n", float64(after.System))

}

// TODO: add memory also
// Research usage of different amount of open challenges, like max 5 vs. all challenges running
func TestChallengeLoad(t *testing.T) {
	//does not work -> maybe calling 'go setupLog' might keep the file open?
	//setupLog("Challenge-load")

	file, err := os.OpenFile("Challenge-load", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer file.Close()

	log.SetOutput(file)
	currentTime := time.Now()
	log.Printf("Testing started " + currentTime.Format("2006.01.02 15:04:05"))

	clientSet := getClientSet()
	teamName := "test"

	//6 challenges

	cpuBeforeFew, err := cpu.Get()
	utils.ErrHandler(err)
	log.Printf("CPU before running 6 challenges %f\n", float64(cpuBeforeFew.System))

	log.Printf("Setting up namespace etc. for 6 challenges\n")
	setUpKubernetesResources(*clientSet, teamName)
	startAllChallenges(*clientSet, teamName)

	time.Sleep(30 * time.Second)
	cpuAfterFew, err := cpu.Get()
	utils.ErrHandler(err)
	log.Printf("CPU after/while running 6 challenges %f\n", float64(cpuAfterFew.System))

	namespaces.DeleteNamespace(*clientSet, teamName)
	log.Printf("Sleeping 90 seconds to allow for namespace to be deleted")
	time.Sleep(90 * time.Second)

	cpuBeforeMany, err := cpu.Get()
	utils.ErrHandler(err)
	log.Printf("CPU 90 seconds after running 6 challenges and deleting namespace, before running 30 %f\n", float64(cpuBeforeMany.System))

	//30 challenges
	//TODO: det er et problem at vi prøver at starte "logon1" fordi den hiver fra map...
	//-> deployment/pod/service name og imagename skal adskilles i param...
	log.Printf("Setting up namespace etc. for 30 challenges\n")
	setUpKubernetesResources(*clientSet, teamName)
	startAllChallengesWithDuplicates(*clientSet, teamName)

	time.Sleep(30 * time.Second)
	cpuAfterMany, err := cpu.Get()
	utils.ErrHandler(err)
	log.Printf("CPU after/while running 30 challenges %f\n", float64(cpuAfterMany.System))

	namespaces.DeleteNamespace(*clientSet, teamName)
	time.Sleep(30 * time.Second)

}
