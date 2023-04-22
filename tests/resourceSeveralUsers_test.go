package tests

import (
	"k8-project/namespaces"
	"k8-project/utils"
	"log"
	"strconv"
	"testing"
	"time"
)

// Find out how many users there can be run on a minimal kubernetes requirements server setup (with an amount of challenges running) while we wait in between the starting of namespaces
// TODO mememory might be relevant
func TestMaximumLoad(t *testing.T) {
	/// 50/50 kali wireguard
	// Alle namespace kører 5 challenges
	utils.SetLogTest("Minimal-k8s-den-anden", false)

	//
}

// Find out how many users there can be run on a minimal kubernetes requirements, stress testing how many namespaces can start at the same time.
// TODO mememory might be relevant
func TestMaximumStartUp(t *testing.T) {
	/// 50/50 kali wireguard
	// Alle namespace kører 5 challenges
	utils.SetLogTest("Minimal-k8s-den-ene", false)

	//
}

// Find out how much resource usage there is for decently size competition (maybe the amount of people of who participate in cybermesterskaberne).
func TestDeprecatedFocusOnAboveTESTS(t *testing.T) {
	//TODO DeprecatedFocusOnAboveTESTS but reuse code maybe

	utils.SetLogTest("Championship", false)

	clientSet := getClientSet()

	comChannel := make(chan string)
	var results string
	go logCPUWithStoredResult(comChannel, &results)

	const amountOfPeople = 350
	people := [amountOfPeople]string{}

	settings := utils.ReadYaml("settings-test.yaml")

	for i := 0; i < amountOfPeople; i++ {
		is := strconv.Itoa(i)
		personI := "person" + is
		people[i] = personI
		setUpKubernetesResourcesWithWireguard(*clientSet, personI, settings.Endpoint, settings.Subnet)
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
