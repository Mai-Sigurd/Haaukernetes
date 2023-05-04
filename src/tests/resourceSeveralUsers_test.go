package tests

import (
	"fmt"
	"k8s-project/utils"
	"testing"
	"time"
)

const waitTime2 = 2 * time.Second

// Find out how many users there can be run on a minimal kubernetes requirements server setup (with an amount of challenges running) while we wait in between the starting of namespaces
func TestMaximumLoad(t *testing.T) {
	/// 50/50 kali wireguard
	// Alle namespace kører 5 challenges
	utils.SetLogTest("testMaximumLoad")
	utils.TestLogger.Println("Test started")
	clientSet := getClientSet()
	counter := 0
	user := "maximum-load-user"

	for {
		namespace := fmt.Sprintf(user+"%d", counter)
		if counter%2 == 0 {
			err := setUpKubernetesResourcesWithWireguard(*clientSet, namespace, utils.WireguardEndpoint, utils.WireguardSubnet)
			if err != nil {
				utils.TestLogger.Println(err.Error())
				utils.TestLogger.Printf("Error setting up namespace and wireguard for namespace %s - shutting down test\n", namespace)
				break
			}
		} else {
			err := setUpKubernetesResourcesWithKali(*clientSet, namespace)
			if err != nil {
				utils.TestLogger.Println(err.Error())
				utils.TestLogger.Printf("Error setting up namespace and wireguard for namespace %s - shutting down test\n", namespace)
				break
			}
		}

		err := startAllChallenges(*clientSet, namespace)
		if err != nil {
			utils.TestLogger.Println(err.Error())
			utils.TestLogger.Printf("Error setting starting all challenges for namespace %s - shutting down test\n", namespace)
			break
		}

		utils.TestLogger.Printf("Started namespace %s", namespace)
		counter++
		time.Sleep(waitTime2)
	}

	utils.TestLogger.Printf("Maximum load test done - successfully created %d namespaces \n", counter)
	utils.TestLogger.Println("Maximum load test ended - remember to manually delete namespaces")
}

// Find out how many users there can be run on a minimal kubernetes requirements, stress testing how many namespaces can start at the same time.
// TODO mememory might be relevant
func TestMaximumStartUp(t *testing.T) {
	/// 50/50 kali wireguard
	// Alle namespace kører 5 challenges
	utils.SetLogTest("TestMaximumStartUp")
	utils.TestLogger.Println("Test started")
	clientSet := getClientSet()
	counter := 0
	user := "maximum-startup-user"

	//attempting to use channels to communicate errors as return values from goroutines are not possible

	errorChannel := make(chan string)
	channelOutput := ""
	go func() {
		channelOutput = <-errorChannel
	}()

	for {
		namespace := fmt.Sprintf(user+"%d", counter)
		if counter%2 == 0 {
			go setUpKubernetesResourcesWithWireguardAndChannel(*clientSet, namespace, utils.WireguardEndpoint, utils.WireguardSubnet, errorChannel)
			for channelOutput != "" {
				utils.TestLogger.Println(channelOutput)
				utils.TestLogger.Printf("Error setting up namespace and wireguard for namespace %s - shutting down test\n", namespace)
				break
			}
		} else {
			go setUpKubernetesResourcesWithKaliAndChannel(*clientSet, namespace, errorChannel)
			for channelOutput != "" {
				utils.TestLogger.Println(channelOutput)
				utils.TestLogger.Printf("Error setting up namespace and wireguard for namespace %s - shutting down test\n", namespace)
				break
			}
		}

		err := startAllChallenges(*clientSet, namespace)
		if err != nil {
			utils.TestLogger.Println(err.Error())
			utils.TestLogger.Printf("Error setting starting all challenges for namespace %s - shutting down test\n", namespace)
			break
		}

		utils.TestLogger.Printf("Started namespace %s", namespace)
		counter++
		time.Sleep(waitTime2)
	}

	utils.TestLogger.Printf("Maximum startup test done - successfully created %d namespaces \n", counter)
	utils.TestLogger.Println("Maximum startup test ended - remember to manually delete namespaces")
}
