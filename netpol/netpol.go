package netpol

import (
	"context"
	utils "k8s-project/utils"

	v1 "k8s.io/api/core/v1"
	networking "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
)

func CreateEgressPolicy(clientSet kubernetes.Clientset, namespace string) {
	policyName := "egress-policy"
	policyTypes := []networking.PolicyType{"Egress"}
	egress := buildEgressRules()
	matchLabels := make(map[string]string)
	matchLabels["app"] = "kali"
	matchLabels["vpn"] = "wireguard"
	createNetworkPolicy(clientSet, policyName, namespace, policyTypes, egress, nil, matchLabels)
}

func CreateChallengeIngressPolicy(clientSet kubernetes.Clientset, namespace string) {
	policyName := "ingress-policy"
	policyTypes := []networking.PolicyType{"Ingress"}

	ingress := buildIngressRules()
	matchLabels := make(map[string]string)
	matchLabels["type"] = "challenge"
	createNetworkPolicy(clientSet, policyName, namespace, policyTypes, nil, ingress, matchLabels)
}

func buildEgressRules() []networking.NetworkPolicyEgressRule {
	getAddress := func(s v1.Protocol) *v1.Protocol { return &s }
	return []networking.NetworkPolicyEgressRule{
		{
			To: []networking.NetworkPolicyPeer{
				{
					PodSelector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							"type": "challenge",
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
}

func buildIngressRules() []networking.NetworkPolicyIngressRule {
	return []networking.NetworkPolicyIngressRule{
		{
			From: []networking.NetworkPolicyPeer{
				{
					PodSelector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							"app": "kali",
						},
					},
				},
			},
		},
		{
			From: []networking.NetworkPolicyPeer{
				{
					PodSelector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							"vpn": "wireguard",
						},
					},
				},
			},
		},
	}
}

// todo many params not very pretty
func createNetworkPolicy(clientSet kubernetes.Clientset, policyName string, namespace string, policyTypes []networking.PolicyType, egress []networking.NetworkPolicyEgressRule, ingress []networking.NetworkPolicyIngressRule, matchLabels map[string]string) {
	netpol := configureNetworkPolicy(policyName, namespace, policyTypes, egress, ingress, matchLabels)
	networkClient := clientSet.NetworkingV1().NetworkPolicies(namespace)
	result, err := networkClient.Create(context.TODO(), &netpol, metav1.CreateOptions{})
	utils.ErrLogger(err)
	utils.InfoLogger.Printf("Created network policy of type %q with name %q for namespace %s\n", &result.Spec.PolicyTypes, result.GetObjectMeta().GetName(), namespace)
}

// many params not very pretty
func configureNetworkPolicy(policyName string, namespace string, policyTypes []networking.PolicyType, egress []networking.NetworkPolicyEgressRule, ingress []networking.NetworkPolicyIngressRule, matchLabels map[string]string) networking.NetworkPolicy {
	netpol := &networking.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      policyName,
			Namespace: namespace,
		},
		Spec: networking.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{
				MatchLabels: matchLabels,
			},
			PolicyTypes: policyTypes,
			Egress:      egress,
			Ingress:     ingress,
		},
	}
	return *netpol
}
