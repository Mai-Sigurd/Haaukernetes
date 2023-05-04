package tests

import (
	"k8s-project/namespaces"
	"k8s-project/utils"
	"testing"
	"time"
)

const waitTime10 = 10 * time.Second

func TestAll6challenges3times(t *testing.T) {
	Test6ChallengesAndWireguard(t)
	time.Sleep(1 * time.Minute)
	Test6ChallengesAndWireguard(t)
	time.Sleep(1 * time.Minute)
	Test6ChallengesAndWireguard(t)
	time.Sleep(1 * time.Minute)

	Test6ChallengesAndKali(t)
	time.Sleep(1 * time.Minute)
	Test6ChallengesAndKali(t)
	time.Sleep(1 * time.Minute)
	Test6ChallengesAndKali(t)
	time.Sleep(1 * time.Minute)
}

func TestAll30challenges3times(t *testing.T) {
	Test30ChallengesWireguard(t)
	time.Sleep(1 * time.Minute)
	Test30ChallengesWireguard(t)
	time.Sleep(1 * time.Minute)
	Test30ChallengesWireguard(t)
	time.Sleep(1 * time.Minute)

	Test30ChallengesKali(t)
	time.Sleep(1 * time.Minute)
	Test30ChallengesKali(t)
	time.Sleep(1 * time.Minute)
	Test30ChallengesKali(t)
	time.Sleep(1 * time.Minute)
}

// General load (resources used for new user, kali docker(simple vs kali many tools), wireguard, guacamole, etc)
func Test6ChallengesAndWireguard(t *testing.T) {
	utils.SetLogTest("6ChallengesAndWireguardOneUser")

	utils.TestLogger.Println("Test started")

	// Starting the kuberneets
	clientSet := getClientSet()
	person1 := "test-6-challenges-wireguard-one-user"

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
	person1 := "test-6-challenges-kali-one-user"
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
	person1 := "test-30-challenges-wireguard-one-user"

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
	person1 := "test-30-challenges-kali-one-user"

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
