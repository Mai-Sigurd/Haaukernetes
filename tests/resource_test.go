package tests

import (
	"log"
	"os"
	"testing"
	"time"
)

//TODO: defer works or not?!
func setupLog(filename string) {
	currentTime := time.Now()
	file, err := os.OpenFile(currentTime.Format("2006.01.02 15:04:05")+filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer file.Close()

	log.SetOutput(file)
}

//test functions follow the order in our notion file containing test cases

func TestGeneralLoad(t *testing.T) {
	setupLog("General-load")
	//now for testi westi
}

func TestMinimalKubernetesSetup(t *testing.T) {
	setupLog("Minimal-k8s")
	//now for testi westi
}

func TestChampionshipLoad(t *testing.T) {
	setupLog("Championship")
	//now for testi westi
}

func TestChallengeLoad(t *testing.T) {
	setupLog("Challenge-load")
	//now for testi westi
}
