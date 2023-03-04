package networkpolicies

import (
	"context"
	"fmt"
	utils "k8-project/utils"

	//appsv1 "k8s.io/api/apps/v1"
	//apiv1 "k8s.io/api/core/v1"
	network "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	//networkv1 "k8s.io/client-go/kubernetes/typed/networking/v1"
)



func createExerciseIngressPolicy() {
	fmt.Println("Hello World")
}
func createKaliEgressPolicy() {
	fmt.Println("Hello World")
}

//demo
func createNetworkPolicy(clientSet kubernetes.Clientset, teamName string) {
	networkPolicy := &network.NetworkPolicy {} //placeholder 
	fmt.Printf("Creating network policy %s\n", networkPolicy.ObjectMeta.Name)
	networkClient := clientSet.NetworkingV1().NetworkPolicies(teamName)
	result, err := networkClient.Create(context.TODO(), networkPolicy, metav1.CreateOptions{})
	utils.ErrHandler(err)
	fmt.Printf("Created network policy %q.\n", result.GetObjectMeta().GetName())
}


//how generic do we want to make this? 
//we need to configure both ingress and egress, but mostly in the same way
//for now i've made it veryyyy generic (except for the podselector, hihi)
func configureNetworkPolicy(name string, podLabel string, nameSpace string, policyTypes []network.PolicyType, ingress []network.NetworkPolicyIngressRule, egress []network.NetworkPolicyEgressRule) network.NetworkPolicy {
	netpol := &network.NetworkPolicy {
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: network.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": podLabel,
				},
			},
			PolicyTypes: policyTypes,
			Ingress: ingress, 
			Egress: egress, 
		},
	}
        return *netpol
}
