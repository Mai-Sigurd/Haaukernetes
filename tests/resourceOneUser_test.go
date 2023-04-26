package tests

import (
	"k8s-project/namespaces"
	"k8s-project/utils"
	"testing"
	"time"
)

// General load (resources used for new user, kali docker(simple vs kali many tools), wireguard, guacamole, etc)
func Test6ChallengesAndWireguard(t *testing.T) {
	utils.SetLogTest("6ChallengesAndWireguardOneUser", false)

	utils.TestLogger.Println("Test started")

	// Starting the kuberneets
	clientSet := getClientSet()
	person1 := "testperson1"

	setUpKubernetesResourcesWithWireguard(*clientSet, person1, utils.WireguardEndpoint, utils.WireguardSubnet)

	startAllChallenges(*clientSet, person1)
	time.Sleep(10 * time.Second)

	err1 := namespaces.DeleteNamespace(*clientSet, person1)
	utils.TestLogger.Println(err1.Error())
	utils.TestLogger.Println("Test ended")
}

func Test6ChallengesAndKali(t *testing.T) {
	utils.SetLogTest("6ChallengesAndKaliOneUser", false)
	utils.TestLogger.Println("Test started")

	// Starting the kuberneets
	clientSet := getClientSet()
	person1 := "testperson1"
	setUpKubernetesResourcesWithKali(*clientSet, person1)

	startAllChallenges(*clientSet, person1)
	time.Sleep(10 * time.Second)

	err1 := namespaces.DeleteNamespace(*clientSet, person1)
	utils.TestLogger.Println(err1.Error())
	utils.TestLogger.Println("Test ended")
}

// Research usage of different amount of open challenges, like max 5 vs. all challenges running
func Test30ChallengesWireguard(t *testing.T) {
	utils.SetLogTest("30challengesWireguardOneUser", false)
	utils.TestLogger.Println("Test started")

	// Starting the kuberneets
	clientSet := getClientSet()
	person1 := "testperson1"

	setUpKubernetesResourcesWithWireguard(*clientSet, person1, utils.WireguardEndpoint, utils.WireguardSubnet)

	startAllChallengesWithDuplicates(*clientSet, person1)
	time.Sleep(10 * time.Second)

	err1 := namespaces.DeleteNamespace(*clientSet, person1)
	utils.TestLogger.Println(err1.Error())
	time.Sleep(30 * time.Second)
	utils.TestLogger.Println("Test ended")

}

func Test30ChallengesKali(t *testing.T) {
	utils.SetLogTest("30challengesKaliOneUser", false)
	utils.TestLogger.Println("Test started")

	// Starting the kuberneets
	clientSet := getClientSet()
	person1 := "testperson1"

	setUpKubernetesResourcesWithKali(*clientSet, person1)

	startAllChallengesWithDuplicates(*clientSet, person1)
	time.Sleep(10 * time.Second)

	err1 := namespaces.DeleteNamespace(*clientSet, person1)
	utils.TestLogger.Println(err1.Error())
	time.Sleep(30 * time.Second)
	utils.TestLogger.Println("Test ended")
}
