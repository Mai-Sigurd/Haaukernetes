package netpol

import (
	"context"
	"fmt"
	utils "k8-project/utils"
	"strings"

	v1 "k8s.io/api/core/v1"
	networking "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
)

func CreateEgressPolicy(clientSet kubernetes.Clientset, teamName string) {
	policyName := "egress-policy"
	policyTypes := []networking.PolicyType{"Egress"}
	egress := buildEgressRules()
	matchLabels := make(map[string]string)
	matchLabels["app"] = "kali-vnc"
	matchLabels["vpn"] = "wireguard"
	createNetworkPolicy(clientSet, policyName, teamName, policyTypes, egress, nil, matchLabels)
}

func CreateChallengeIngressPolicy(clientSet kubernetes.Clientset, teamName string) {
	policyName := "ingress-policy"
	policyTypes := []networking.PolicyType{"Ingress"}

	//dynamic handling of pod ip because label selectors are not working for wierd reasons
	podClient := clientSet.CoreV1().Pods(teamName)
	pods, err := podClient.List(context.TODO(), metav1.ListOptions{})
	utils.ErrHandler(err)

	//tried to be upfront about not necessarily knowing which index the pod is at (we could make sure that it's at 0 always but meh)
	podIP := findPodIp(pods)
	fmt.Printf("Wireguard pod ip for namespace %s: %s\n", teamName, podIP)

	ingress := buildIngressRules(podIP)
	matchLabels := make(map[string]string)
	matchLabels["type"] = "challenge"
	createNetworkPolicy(clientSet, policyName, teamName, policyTypes, nil, ingress, matchLabels)
}

func findPodIp(pods *v1.PodList) string {
	for i := range pods.Items {
		if strings.Contains(pods.Items[i].Name, "wireguard") {
			return pods.Items[i].Status.PodIP
		}
	}
	return ""
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

func buildIngressRules(podIP string) []networking.NetworkPolicyIngressRule {
	return []networking.NetworkPolicyIngressRule{
		{
			From: []networking.NetworkPolicyPeer{
				{
					PodSelector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							"app": "kali-vnc",
							"vpn": "wireguard",
						},
					},
				},
				//the labelselector above does not work for wireguard, for unknown reasons... so we use ip directly instead
				{
					IPBlock: &networking.IPBlock{
						CIDR: podIP + "/32",
					},
				},
			},
		},
	}
}

// many params not very pretty
func createNetworkPolicy(clientSet kubernetes.Clientset, policyName string, teamName string, policyTypes []networking.PolicyType, egress []networking.NetworkPolicyEgressRule, ingress []networking.NetworkPolicyIngressRule, matchLabels map[string]string) {
	netpol := configureNetworkPolicy(policyName, teamName, policyTypes, egress, ingress, matchLabels)
	networkClient := clientSet.NetworkingV1().NetworkPolicies(teamName)
	result, err := networkClient.Create(context.TODO(), &netpol, metav1.CreateOptions{})
	utils.ErrHandler(err)
	fmt.Printf("Created network policy of type %q with name %q for namespace %s", &result.Spec.PolicyTypes, result.GetObjectMeta().GetName(), teamName)
}

// many params not very pretty
func configureNetworkPolicy(policyName string, teamName string, policyTypes []networking.PolicyType, egress []networking.NetworkPolicyEgressRule, ingress []networking.NetworkPolicyIngressRule, matchLabels map[string]string) networking.NetworkPolicy {
	netpol := &networking.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      policyName,
			Namespace: teamName,
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
