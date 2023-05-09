package tests

import (
	"k8s-project/namespaces"
	"k8s-project/utils"
	"testing"
	"time"
)

func TestBaselineKali(t *testing.T) {
	utils.SetLogTest("baselineKali")
	utils.TestLogger.Println("Test started")

	// Starting Kubernetes
	clientSet := getClientSet()
	person1 := "baseline-kali"
	setUpKubernetesResourcesWithKali(*clientSet, person1)

	utils.TestLogger.Println("Waiting for the const seconds")
	time.Sleep(wait3minutes)
	utils.TestLogger.Println("Test ended, deleting namespaces")

	err1 := namespaces.DeleteNamespace(*clientSet, person1)
	if err1 != nil {
		utils.TestLogger.Println(err1.Error())
	}
	utils.TestLogger.Println("Test function exits")
}

func TestBaselineWireguard(t *testing.T) {
	utils.SetLogTest("baselineWireguard")
	utils.TestLogger.Println("Test started")

	// Starting Kubernetes
	clientSet := getClientSet()
	person1 := "baseline-wireguard"
	setUpKubernetesResourcesWithWireguard(*clientSet, person1, utils.WireguardEndpoint, utils.WireguardSubnet)

	utils.TestLogger.Println("Waiting for the const seconds")
	time.Sleep(wait3minutes)
	utils.TestLogger.Println("Test ended, deleting namespaces")

	err1 := namespaces.DeleteNamespace(*clientSet, person1)
	if err1 != nil {
		utils.TestLogger.Println(err1.Error())
	}
	utils.TestLogger.Println("Test function exits")
}
