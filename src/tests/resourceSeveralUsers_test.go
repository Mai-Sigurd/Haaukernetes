package tests

import (
	"fmt"
	"k8s-project/namespaces"
	"k8s-project/utils"
	"testing"
	"time"
)

//This works in the sense that the test terminates when errors are returned from k8s api (e.g. resources are sparse) but the actual
//deletion of namespaces and resources goes on even after the test exits....
//
// Find out how many users there can be run on a minimal kubernetes requirements server setup (with an amount of challenges running) while we wait in between the starting of namespaces
func TestMaximumLoad(t *testing.T) {
	/// 50/50 kali wireguard
	// Alle namespace kører 5 challenges
	utils.SetLogTest("Minimal-k8s-den-anden", true)

	clientSet := getClientSet()
	counter := 0
	teamName := "maximum-load-team"

	for {
		team := fmt.Sprintf(teamName+"%d", counter)
		if counter%2 == 0 {
			err := setUpKubernetesResourcesWithWireguard(*clientSet, team, utils.WireguardEndpoint, utils.WireguardSubnet)
			if err != nil {
				utils.TestLogger.Println(err.Error())
				utils.TestLogger.Printf("Error setting up namespace and wireguard for namespace %s - shutting down test\n", team)
				break
			}
		} else {
			err := setUpKubernetesResourcesWithKali(*clientSet, team)
			if err != nil {
				utils.TestLogger.Println(err.Error())
				utils.TestLogger.Printf("Error setting up namespace and wireguard for namespace %s - shutting down test\n", team)
				break
			}
		}

		err := startAllChallenges(*clientSet, team)
		if err != nil {
			utils.TestLogger.Println(err.Error())
			utils.TestLogger.Printf("Error setting starting all challenges for namespace %s - shutting down test\n", team)
			break
		}

		counter++
		time.Sleep(2 * time.Second)
	}

	utils.TestLogger.Printf("Maximum load test done - successfully created %d namespaces \n", counter)
	utils.TestLogger.Println("Deleting test namespaces")

	for i := 0; i < counter; i++ {
		namespaces.DeleteNamespace(*clientSet, fmt.Sprintf(teamName+"%d", i))
	}
}

// Find out how many users there can be run on a minimal kubernetes requirements, stress testing how many namespaces can start at the same time.
// TODO mememory might be relevant
func TestMaximumStartUp(t *testing.T) {
	/// 50/50 kali wireguard
	// Alle namespace kører 5 challenges
	utils.SetLogTest("Minimal-k8s-den-ene", true)

	clientSet := getClientSet()
	counter := 0
	teamName := "maximum-startup-team"

	//attempting to use channels to communicate errors as return values from goroutines are not possible

	errorChannel := make(chan string)
	channelOutput := ""
	go func() {
		channelOutput = <-errorChannel
	}()

	for {
		team := fmt.Sprintf(teamName+"%d", counter)
		if counter%2 == 0 {
			go setUpKubernetesResourcesWithWireguardAndChannel(*clientSet, team, utils.WireguardEndpoint, utils.WireguardSubnet, errorChannel)
			for channelOutput != "" {
				utils.TestLogger.Println(channelOutput)
				utils.TestLogger.Printf("Error setting up namespace and wireguard for namespace %s - shutting down test\n", team)
				break
			}
		} else {
			go setUpKubernetesResourcesWithKaliAndChannel(*clientSet, team, errorChannel)
			for channelOutput != "" {
				utils.TestLogger.Println(channelOutput)
				utils.TestLogger.Printf("Error setting up namespace and wireguard for namespace %s - shutting down test\n", team)
				break
			}
		}

		err := startAllChallenges(*clientSet, team)
		if err != nil {
			utils.TestLogger.Println(err.Error())
			utils.TestLogger.Printf("Error setting starting all challenges for namespace %s - shutting down test\n", team)
			break
		}

		counter++
		time.Sleep(2 * time.Second)
	}

	utils.TestLogger.Printf("Maximum load test done - successfully created %d namespaces \n", counter)
	utils.TestLogger.Println("Deleting test namespaces")

	for i := 0; i < counter; i++ {
		namespaces.DeleteNamespace(*clientSet, fmt.Sprintf(teamName+"%d", i))
	}
}
