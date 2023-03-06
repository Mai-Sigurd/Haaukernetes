package netpol

import (
	"context"
	"fmt"
	utils "k8-project/utils"

	v1 "k8s.io/api/core/v1"
	networking "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
)

//!This is not made to be pretty, but as a starting point that works as intended!
//pls refactor into smaller functions etc.

func CreateKaliEgressPolicy(clientSet kubernetes.Clientset, teamName string) {
	//create egress rule
	getAddress := func(s v1.Protocol) *v1.Protocol { return &s }
	policyTypes := []networking.PolicyType{"Egress"}
	egress := []networking.NetworkPolicyEgressRule{
		{
			To: []networking.NetworkPolicyPeer{
				{
					// NamespaceSelector: &metav1.LabelSelector{
					// 	MatchLabels: map[string]string{
					// 		"kubernetes.io/metadata.name": teamName, //target own namespace
					// 	},
					// },
					PodSelector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							"type": "exercise", //exercise pods has this label
						},
					},
				},
			},
		},
		{
			To: []networking.NetworkPolicyPeer{
				{ //this is for allowing internal dns
					NamespaceSelector: &metav1.LabelSelector{},
					PodSelector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							"k8s-app": "kube-dns",
						},
					},
				},
			},
			Ports: []networking.NetworkPolicyPort{
				{
					Port:     &intstr.IntOrString{Type: intstr.Type(intstr.Int), IntVal: 53},
					Protocol: getAddress(v1.ProtocolUDP),
				},
			},
		},
	}
	//configure policy
	netpol := &networking.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "egress-policy",
			Namespace: teamName,
		},
		Spec: networking.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "vnc",
				},
			},
			PolicyTypes: policyTypes,
			Egress:      egress,
		},
	}

	//create policy
	networkClient := clientSet.NetworkingV1().NetworkPolicies(teamName)
	result, err := networkClient.Create(context.TODO(), netpol, metav1.CreateOptions{})
	utils.ErrHandler(err)
	fmt.Printf("Created egress policy: %q for namespace %s", result.GetObjectMeta().GetName(), teamName)
}

func CreateExerciseIngressPolicy(clientSet kubernetes.Clientset, teamName string) {
	//create ingress rule
	policyTypes := []networking.PolicyType{"Ingress"}
	ingress := []networking.NetworkPolicyIngressRule{
		{
			From: []networking.NetworkPolicyPeer{
				{
					//just deleting namespaceselector might be ok, see comment in above function
					//NamespaceSelector: &metav1.LabelSelector{},
					PodSelector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							"app": "vnc", //needs to align with kali labels
						},
					},
				},
			},
		},
	}

	//configure policy
	netpol := &networking.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "ingress-policy",
			Namespace: teamName,
		},
		Spec: networking.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{
				MatchLabels: map[string]string{
					"type": "exercise",
				},
			},
			PolicyTypes: policyTypes,
			Ingress:     ingress,
		},
	}

	//create policy
	networkClient := clientSet.NetworkingV1().NetworkPolicies(teamName)
	result, err := networkClient.Create(context.TODO(), netpol, metav1.CreateOptions{})
	utils.ErrHandler(err)
	fmt.Printf("Created ingress policy: %q for namespace %s", result.GetObjectMeta().GetName(), teamName)
}
