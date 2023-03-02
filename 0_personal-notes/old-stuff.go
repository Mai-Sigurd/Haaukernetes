package __personal_notes

import (
	"fmt"
	"k8-project/deployments"
	"k8-project/utils"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
)

func old_main() {
	home := homedir.HomeDir()
	kubeconfig := filepath.Join(home, ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	utils.ErrHandler(err)
	clientset, err := kubernetes.NewForConfig(config)
	utils.ErrHandler(err)

	deploymentsClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)
	//deploymentsClient := clientset.AppsV1().Deployments("user-a")
	deployments.ListDeployments(deploymentsClient)

	//createDeployment(deploymentsClient, logon())

	//name := "user-a-logon"
	//delete_deployment(deploymentsClient, name)

	//listDeployments(deploymentsClient)

	//create_namespace(*clientset, *namespace_test())
}

func oldoldmain() {

	home := homedir.HomeDir()
	kubeconfig := filepath.Join(home, ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	utils.ErrHandler(err)
	clientset, err := kubernetes.NewForConfig(config)
	utils.ErrHandler(err)
	fmt.Print(clientset)
	//main.browser(*clientset)
}
