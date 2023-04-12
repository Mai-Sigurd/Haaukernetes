package tests

import (
	"fmt"
	"k8-project/apis"
	"k8-project/deployments"
	"k8-project/namespaces"
	"k8-project/netpol"
	"k8-project/secrets"
	"k8-project/services"
	"k8-project/wireguard"
	"k8s.io/client-go/kubernetes"
	"log"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/mackerelio/go-osstat/cpu"
)

// TODO: defer works or not?!
func setupLog(filename string) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer file.Close()

	log.SetOutput(file)
	currentTime := time.Now()
	log.Printf("Testing started " + currentTime.Format("2006.01.02 15:04:05"))
}

func setUpKubernetesResources(clientSet kubernetes.Clientset, teamName string) {
	namespaces.CreateNamespace(clientSet, teamName)
	secrets.CreateImageRepositorySecret(clientSet, teamName)
	netpol.CreateChallengeIngressPolicy(clientSet, teamName)
	netpol.CreateEgressPolicy(clientSet, teamName)
	wireguard.StartWireguard(clientSet, teamName, "2A/Rj6X3+YxP6lXOv2BgbRQfpCn5z6Ob8scKhxiCRyM=") //random publickey
	netpol.AddWireguardToChallengeIngressPolicy(clientSet, teamName)
}

func startChallenge(challengeName string, clientSet kubernetes.Clientset, namespace string) {
	challengePorts := ports[challengeName]
	podLabels := make(map[string]string)
	podLabels["app"] = challengeName
	podLabels["type"] = "challenge"
	deployments.CreateLocalDeployment(clientSet, namespace, challengeName, challengePorts, podLabels)
	services.CreateService(clientSet, namespace, challengeName, challengePorts)
}

func startAllChallenges(clientSet kubernetes.Clientset, namespace string) {
	for key, _ := range ports {
		startChallenge(key, clientSet, namespace)
	}
}

// General load (resources used for new user, kali docker(simple vs kali many tools), wireguard, guacamole, etc)
func TestGeneralLoad(t *testing.T) {
	setupLog("General-load")

	// CPU Load before starting
	before, err := cpu.Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return
	}

	// Starting the kuberneets
	clientSet := getClientSet()
	personA := "PersonA"

	setUpKubernetesResources(*clientSet, personA)

	//TODO do we wanna start all challenges for this test?
	startAllChallenges(*clientSet, personA)
	apis.StartKali(*clientSet, personA)

	time.Sleep(5 * time.Second)

	// CPU Load after starting
	after, err1 := cpu.Get()
	if err1 != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err1)
		return
	}

	// TODO test om det er det output vi rent faktisk vil have
	log.Printf("cpu system before test: %f %%\n", float64(before.System))
	log.Printf("cpu system after test: %f %%\n", float64(after.System))

}

// Find out how many users there can be run on a minimal kubernetes requirements server setup (with an amount of challenges running)
func TestMinimalKubernetesSetup(t *testing.T) {
	setupLog("Minimal-k8s")
	//  TODO How do we actually stress test Kubernetes?
	//now for testi westi
}

// Find out how much resource usage there is for decently size competition (maybe the amount of people of who participate in cybermesterskaberne).
func TestChampionshipLoad(t *testing.T) {
	setupLog("Championship")
	clientSet := getClientSet()

	// TODO correct amount
	const amountOfPeople = 100
	people := [amountOfPeople]string{}

	for i := 0; i < amountOfPeople; i++ {
		is := strconv.Itoa(i)
		personI := "Person" + is
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

// Research usage of different amount of open challenges, like max 5 vs. all challenges running
func TestChallengeLoad(t *testing.T) {
	setupLog("Challenge-load")
	//now for testi westi
}
