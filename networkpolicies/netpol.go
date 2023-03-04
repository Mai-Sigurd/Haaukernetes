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



func createExerciseIngressPolicy(clientSet kubernetes.Clientset, teamName string, podName string) {
	policyTypes := []network.PolicyType {}
	ingress := []network.NetworkPolicyIngressRule {} //empty rule because we need the param and nil might screw us but empty rule might 
	//also screw us as all ingress might just be denied...
	egress := []network.NetworkPolicyEgressRule {
		{
			To: []network.NetworkPolicyPeer{}, //fill me 
			Ports: []network.NetworkPolicyPort{}, //fill me
		},
	} 

	createNetworkPolicy(clientSet, podName, "remember to handle podlabel", teamName, policyTypes, ingress, egress)

}
func createKaliEgressPolicy() {
	fmt.Println("Hello World")
}

//demo
func createNetworkPolicy(clientSet kubernetes.Clientset, name string, podLabel string, teamName string, policyTypes []network.PolicyType, ingress []network.NetworkPolicyIngressRule, egress []network.NetworkPolicyEgressRule) {
	networkPolicy := configureNetworkPolicy(name, podLabel, teamName,policyTypes, ingress, egress) 
	fmt.Printf("Creating network policy %s\n", networkPolicy.ObjectMeta.Name)
	networkClient := clientSet.NetworkingV1().NetworkPolicies(teamName)
	result, err := networkClient.Create(context.TODO(), &networkPolicy, metav1.CreateOptions{})
	utils.ErrHandler(err)
	fmt.Printf("Created network policy %q.\n", result.GetObjectMeta().GetName())
}


//how generic do we want to make this? 
//we need to configure both ingress and egress, but mostly in the same way
//for now i've made it veryyyy generic (except for the podselector, hihi)
func configureNetworkPolicy(name string, podLabel string, teamName string, policyTypes []network.PolicyType, ingress []network.NetworkPolicyIngressRule, egress []network.NetworkPolicyEgressRule) network.NetworkPolicy {
	netpol := &network.NetworkPolicy {
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Namespace: teamName,
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
