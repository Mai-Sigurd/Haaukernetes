package tests

import (
	"k8s-project/namespaces"
	"k8s-project/utils"
	"testing"
	"time"
)

const waitTime10 = 10 * time.Second

// General load (resources used for new user, kali docker(simple vs kali many tools), wireguard, guacamole, etc)
func Test6ChallengesAndWireguard(t *testing.T) {
	utils.SetLogTest("6ChallengesAndWireguardOneUser")

	utils.TestLogger.Println("Test started")

	// Starting the kuberneets
	clientSet := getClientSet()
	person1 := "Test6challengesWireguardOneUser"

	setUpKubernetesResourcesWithWireguard(*clientSet, person1, utils.WireguardEndpoint, utils.WireguardSubnet)

	startAllChallenges(*clientSet, person1)
	utils.TestLogger.Println("Waiting for the const seconds")
	time.Sleep(waitTime10)
	utils.TestLogger.Println("Test ended, deleting namespaces")
	err1 := namespaces.DeleteNamespace(*clientSet, person1)
	if err1 != nil {
		utils.TestLogger.Println(err1.Error())
	}
	utils.TestLogger.Println("Deleted namespaces")
}

func Test6ChallengesAndKali(t *testing.T) {
	utils.SetLogTest("6ChallengesAndKaliOneUser")
	utils.TestLogger.Println("Test started")

	// Starting the kuberneets
	clientSet := getClientSet()
	person1 := "Test6challengesKaliOneUser"
	setUpKubernetesResourcesWithKali(*clientSet, person1)

	startAllChallenges(*clientSet, person1)
	utils.TestLogger.Println("Waiting for the const seconds")
	time.Sleep(waitTime10)
	utils.TestLogger.Println("Test ended, deleting namespaces")

	err1 := namespaces.DeleteNamespace(*clientSet, person1)
	if err1 != nil {
		utils.TestLogger.Println(err1.Error())
	}
	utils.TestLogger.Println("Deleted namespaces")
}

// Research usage of different amount of open challenges, like max 5 vs. all challenges running
func Test30ChallengesWireguard(t *testing.T) {
	utils.SetLogTest("30challengesWireguardOneUser")
	utils.TestLogger.Println("Test started")

	// Starting the kuberneets
	clientSet := getClientSet()
	person1 := "Test30challengesWireguardOneUser"

	setUpKubernetesResourcesWithWireguard(*clientSet, person1, utils.WireguardEndpoint, utils.WireguardSubnet)

	startAllChallengesWithDuplicates(*clientSet, person1)
	utils.TestLogger.Println("Waiting for the const seconds")
	time.Sleep(waitTime10)
	utils.TestLogger.Println("Test ended, deleting namespaces")
	err1 := namespaces.DeleteNamespace(*clientSet, person1)
	if err1 != nil {
		utils.TestLogger.Println(err1.Error())
	}
	utils.TestLogger.Println("Waiting for the const seconds")
	time.Sleep(waitTime10)
	utils.TestLogger.Println("Deleted namespaces")

}

func Test30ChallengesKali(t *testing.T) {
	utils.SetLogTest("30challengesKaliOneUser")
	utils.TestLogger.Println("Test started")

	// Starting the kuberneets
	clientSet := getClientSet()
	person1 := "Test30challengesKaliOneUser"

	setUpKubernetesResourcesWithKali(*clientSet, person1)

	startAllChallengesWithDuplicates(*clientSet, person1)
	utils.TestLogger.Println("Waiting for the const seconds")
	time.Sleep(waitTime10)
	utils.TestLogger.Println("Test ended, deleting namespaces")

	err1 := namespaces.DeleteNamespace(*clientSet, person1)
	if err1 != nil {
		utils.TestLogger.Println(err1.Error())
	}
	utils.TestLogger.Println("Waiting for the const seconds")
	time.Sleep(waitTime10)
	utils.TestLogger.Println("Deleted everything")
}
