package tests

import (
	"k8-project/namespaces"
	"k8-project/utils"
	"log"
	"testing"
	"time"
)

// General load (resources used for new user, kali docker(simple vs kali many tools), wireguard, guacamole, etc)
func Test6ChallengesAndWireguard(t *testing.T) {
	file := setupLog("6ChallengesAndWireguardOneUser")
	defer file.Close()
	log.Printf("Test started")

	// Starting the kuberneets
	clientSet := getClientSet()
	person1 := "testperson1"

	settings := utils.ReadYaml("settings-test.yaml")
	setUpKubernetesResourcesWithWireguard(*clientSet, person1, settings.Endpoint, settings.Subnet)

	startAllChallenges(*clientSet, person1)
	time.Sleep(10 * time.Second)

	err1 := namespaces.DeleteNamespace(*clientSet, person1)
	utils.ErrHandler(err1)
	log.Printf("Test ended")
}

func Test6ChallengesAndKali(t *testing.T) {
	file := setupLog("6ChallengesAndKaliOneUser")
	defer file.Close()
	log.Printf("Test started")

	// Starting the kuberneets
	clientSet := getClientSet()
	person1 := "testperson1"
	setUpKubernetesResourcesWithKali(*clientSet, person1)

	startAllChallenges(*clientSet, person1)
	time.Sleep(10 * time.Second)

	err1 := namespaces.DeleteNamespace(*clientSet, person1)
	utils.ErrHandler(err1)
	log.Printf("Test ended")
}

// Research usage of different amount of open challenges, like max 5 vs. all challenges running
func Test30ChallengesWireguard(t *testing.T) {
	file := setupLog("30challengesWireguardOneUser")
	defer file.Close()
	log.Printf("Test started")

	// Starting the kuberneets
	clientSet := getClientSet()
	person1 := "testperson1"

	settings := utils.ReadYaml("settings-test.yaml")
	setUpKubernetesResourcesWithWireguard(*clientSet, person1, settings.Endpoint, settings.Subnet)

	startAllChallengesWithDuplicates(*clientSet, person1)
	time.Sleep(10 * time.Second)

	err1 := namespaces.DeleteNamespace(*clientSet, person1)
	utils.ErrHandler(err1)
	time.Sleep(30 * time.Second)
	log.Printf("Test ended")
}

func Test30ChallengesKali(t *testing.T) {
	file := setupLog("30challengesKaliOneUser")
	defer file.Close()
	log.Printf("Test started")

	// Starting the kuberneets
	clientSet := getClientSet()
	person1 := "testperson1"

	setUpKubernetesResourcesWithKali(*clientSet, person1)

	startAllChallengesWithDuplicates(*clientSet, person1)
	time.Sleep(10 * time.Second)

	err1 := namespaces.DeleteNamespace(*clientSet, person1)
	utils.ErrHandler(err1)
	time.Sleep(30 * time.Second)
	log.Printf("Test ended")
}
