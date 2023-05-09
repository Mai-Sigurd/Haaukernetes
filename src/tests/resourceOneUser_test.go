package tests

import (
	"k8s-project/namespaces"
	"k8s-project/utils"
	"testing"
	"time"
)

// General load (resources used for new user, kali docker(simple vs kali many tools), wireguard, guacamole, etc)
func Test5ChallengesAndKali(t *testing.T) {
	utils.SetLogTest("5ChallengesAndKaliOneUser")
	utils.TestLogger.Println("Test started")

	// Starting the Kubernetes
	clientSet := getClientSet()
	person1 := "test-5-challenges-kali-one-user"
	setUpKubernetesResourcesWithKali(*clientSet, person1)

	startAllChallenges(*clientSet, person1)
	utils.TestLogger.Println("Waiting for the const seconds")
	time.Sleep(wait3minutes)
	utils.TestLogger.Println("Test ended, deleting namespaces")

	err1 := namespaces.DeleteNamespace(*clientSet, person1)
	if err1 != nil {
		utils.TestLogger.Println(err1.Error())
	}
	time.Sleep(wait3minutes)
	utils.TestLogger.Println("Test function exits")
}

func Test5ChallengesAndWireguard(t *testing.T) {
	utils.SetLogTest("5ChallengesAndWireguardOneUser")

	utils.TestLogger.Println("Test started")

	// Starting the Kubernetes
	clientSet := getClientSet()
	person1 := "test-5-challenges-wireguard-one-user"

	setUpKubernetesResourcesWithWireguard(*clientSet, person1, utils.WireguardEndpoint, utils.WireguardSubnet)
	startAllChallenges(*clientSet, person1)
	utils.TestLogger.Println("Waiting for the const seconds")
	time.Sleep(wait3minutes)
	utils.TestLogger.Println("Test ended, deleting namespaces, waiting for 3 minutes")
	err1 := namespaces.DeleteNamespace(*clientSet, person1)
	if err1 != nil {
		utils.TestLogger.Println(err1.Error())
	}
	time.Sleep(wait3minutes)
	utils.TestLogger.Println("Test function exits")
}

// Research usage of different amount of open challenges, like max 5 vs. all challenges running
func Test30ChallengesKali(t *testing.T) {
	utils.SetLogTest("30challengesKaliOneUser")
	utils.TestLogger.Println("Test started")

	// Starting Kubernetes
	clientSet := getClientSet()
	person1 := "test-30-challenges-kali-one-user"

	setUpKubernetesResourcesWithKali(*clientSet, person1)

	startAllChallengesWithDuplicates(*clientSet, person1)
	utils.TestLogger.Println("Waiting for the const seconds")
	time.Sleep(wait3minutes)
	utils.TestLogger.Println("Test ended, deleting namespaces")

	err1 := namespaces.DeleteNamespace(*clientSet, person1)
	if err1 != nil {
		utils.TestLogger.Println(err1.Error())
	}
	utils.TestLogger.Println("Waiting for the const seconds")
	time.Sleep(wait3minutes)
	utils.TestLogger.Println("Test function exits")
}

func Test30ChallengesWireguard(t *testing.T) {
	utils.SetLogTest("30challengesWireguardOneUser")
	utils.TestLogger.Println("Test started")

	// Starting Kubernetes
	clientSet := getClientSet()
	person1 := "test-30-challenges-wireguard-one-user"

	setUpKubernetesResourcesWithWireguard(*clientSet, person1, utils.WireguardEndpoint, utils.WireguardSubnet)

	startAllChallengesWithDuplicates(*clientSet, person1)
	utils.TestLogger.Println("Waiting for the const seconds")
	time.Sleep(wait3minutes)
	utils.TestLogger.Println("Test ended, deleting namespaces")
	err1 := namespaces.DeleteNamespace(*clientSet, person1)
	if err1 != nil {
		utils.TestLogger.Println(err1.Error())
	}
	utils.TestLogger.Println("Waiting for the const seconds")
	time.Sleep(wait3minutes)
	utils.TestLogger.Println("Test function exits")

}
