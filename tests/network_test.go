package tests

import (
	"context"
	"fmt"
	"k8-project/namespaces"
	"k8-project/utils"
	"os/exec"
	"strings"
	"testing"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestPing(t *testing.T) {
	clientSet := getClientSet()

	teamA := "team-a"
	teamB := "team-b"

	settings := utils.ReadYaml("settings-test.yaml")

	setUpKubernetesResourcesWithWireguard(*clientSet, teamA, settings.Endpoint, settings.Subnet)
	name := "logon"
	imageName := "logon"
	startChallenge(name, imageName, *clientSet, teamA, ports[imageName])
	time.Sleep(10 * time.Second)
	setUpKubernetesResourcesWithWireguard(*clientSet, teamB, settings.Endpoint, settings.Subnet)
	startChallenge(name, imageName, *clientSet, teamB, ports[imageName])
	time.Sleep(10 * time.Second)

	podClientA := clientSet.CoreV1().Pods(teamA)
	podsA, err := podClientA.List(context.TODO(), metav1.ListOptions{})
	utils.ErrHandler(err)

	podClientB := clientSet.CoreV1().Pods(teamB)
	podsB, err := podClientB.List(context.TODO(), metav1.ListOptions{})
	utils.ErrHandler(err)

	logonIPA := findPodIp(podsA)
	outA, err := exec.Command("/bin/sh", "-c", "kubectl -n team-a exec -it deployment/wireguard -- ping -c 5 "+logonIPA).Output()
	if err != nil {
		fmt.Println(err.Error())
		t.Error("Intra-namespace pinging should not fail")
	}
	fmt.Println("Output from A pinging A:\n" + string(outA))

	logonIPB := findPodIp(podsB)
	outB, err := exec.Command("/bin/sh", "-c", "kubectl -n team-b exec -it deployment/wireguard -- ping -c 5 "+logonIPB).Output()
	if err != nil {
		fmt.Println(err.Error())
		t.Error("Intra-namespace pinging should not fail")
	}
	fmt.Println("Output from B pinging B:\n" + string(outB))

	crossB, _ := exec.Command("/bin/sh", "-c", "kubectl -n team-a exec -it deployment/wireguard -- ping -c 5 "+logonIPB).Output()
	fmt.Println("Output from A pinging B\n" + string(crossB))
	if !strings.Contains(string(crossB), "100% packet loss") {
		t.Error("100% of the packets should be lost during inter-namespace pinging")
	}

	crossA, _ := exec.Command("/bin/sh", "-c", "kubectl -n team-b exec -it deployment/wireguard -- ping -c 5 "+logonIPA).Output()
	fmt.Println("Output from B pinging A\n" + string(crossA))
	if !strings.Contains(string(crossA), "100% packet loss") {
		t.Error("100% of the packets should be lost during inter-namespace pinging")
	}

	err1 := namespaces.DeleteNamespace(*clientSet, teamA)
	err2 := namespaces.DeleteNamespace(*clientSet, teamB)
	utils.ErrHandler(err1)
	utils.ErrHandler(err2)
}

func findPodIp(pods *v1.PodList) string {
	for i := range pods.Items {
		if strings.Contains(pods.Items[i].Name, "logon") {
			return pods.Items[i].Status.PodIP
		}
	}
	return "IP of wireguard pod not found"
}
