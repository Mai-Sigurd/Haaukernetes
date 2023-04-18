package tests

//docs says that tests should live in package/dir of a module, but this is more of an integration/system test, than a simple
//unit test covering 1 module
//https://go.dev/doc/tutorial/add-a-test
//https://pkg.go.dev/testing

import (
	"context"
	"fmt"
	"k8-project/deployments"
	"k8-project/namespaces"
	"k8-project/netpol"
	"k8-project/secrets"
	"k8-project/services"
	"k8-project/utils"
	"k8-project/wireguard"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

//BASIC INFO
//use "t.SkipNow()"" in a test to skip it
//use "go test -v -run FUNCTIONNAME" to only test a single function i.e. "go test -v -run TestResourceUse"

func getClientSet() *kubernetes.Clientset {
	home := homedir.HomeDir()
	kubeConfigPath := filepath.Join(home, ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	utils.ErrHandler(err)
	clientSet, err := kubernetes.NewForConfig(config)
	utils.ErrHandler(err)
	return clientSet
}

func TestCreateAndDeleteNamespace(t *testing.T) {
	clientSet := getClientSet()

	for i := 0; i < 5; i++ {
		namespaces.CreateNamespace(*clientSet, "test"+fmt.Sprint(i))
	}

	for i := 0; i < 5; i++ {
		exists := namespaces.NamespaceExists(*clientSet, "test"+fmt.Sprint(i))
		if !exists {
			t.Error("Namespace not created")
		}
	}

	for i := 0; i < 5; i++ {
		namespaces.DeleteNamespace(*clientSet, "test"+fmt.Sprint(i))
	}

	fmt.Println("Namespace test sleeping to ensure that namespaces have been deleted")
	time.Sleep(10 * time.Second)

	for i := 0; i < 5; i++ {
		exists := namespaces.NamespaceExists(*clientSet, "test"+fmt.Sprint(i))
		if exists {
			t.Error("Namespace not deleted")
		}
	}
}

// waitgroups are used to have concurrency while avoiding using a timer or infinite loop, as goroutines will be killed
// when function returns
func TestResourceUse(t *testing.T) {
	clientSet := getClientSet()
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		teamName := "test" + fmt.Sprint(i)

		wg.Add(1)
		go func() {
			defer wg.Done()
			setUpKubernetesResourcesWithLogon(*clientSet, teamName)
		}()
	}
	wg.Wait()

	//comment out this let all created resources stay alive and wait for manual cleanup
	//manuel cleanup can be done via:
	//'kubectl delete ns $(kubectl get namespaces --no-headers=true -o custom-columns=:metadata.name  | grep test* - )'
	fmt.Println("Sleeping for 60 seconds to let system resources stabilize")
	time.Sleep(60 * time.Second)
	for i := 0; i < 50; i++ {
		teamName := "test" + fmt.Sprint(i)
		namespaces.DeleteNamespace(*clientSet, teamName)
	}

	//can we do any assertion that actually makes sense in this case?
}

func setUpKubernetesResourcesWithLogon(clientSet kubernetes.Clientset, teamName string) {
	challengeName := "logon"
	challengePorts := ports[challengeName]
	podLabels := make(map[string]string)
	podLabels["app"] = challengeName
	podLabels["type"] = "challenge"
	namespaces.CreateNamespace(clientSet, teamName)
	secrets.CreateImageRepositorySecret(clientSet, teamName)
	netpol.CreateChallengeIngressPolicy(clientSet, teamName)
	netpol.CreateEgressPolicy(clientSet, teamName)
	wireguard.StartWireguard(clientSet, teamName, "2A/Rj6X3+YxP6lXOv2BgbRQfpCn5z6Ob8scKhxiCRyM=") //random publickey
	netpol.AddWireguardToChallengeIngressPolicy(clientSet, teamName)
	deployments.CreateDeployment(clientSet, teamName, challengeName, challengePorts, podLabels)
	services.CreateService(clientSet, teamName, challengeName, challengePorts)
}

func TestPing(t *testing.T) {
	clientSet := getClientSet()

	teamA := "team-a"
	teamB := "team-b"

	setUpKubernetesResourcesWithLogon(*clientSet, teamA)
	time.Sleep(10 * time.Second)
	setUpKubernetesResourcesWithLogon(*clientSet, teamB)
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

	namespaces.DeleteNamespace(*clientSet, teamA)
	namespaces.DeleteNamespace(*clientSet, teamB)
}

func findPodIp(pods *v1.PodList) string {
	for i := range pods.Items {
		if strings.Contains(pods.Items[i].Name, "logon") {
			return pods.Items[i].Status.PodIP
		}
	}
	return "IP of wireguard pod not found"
}
