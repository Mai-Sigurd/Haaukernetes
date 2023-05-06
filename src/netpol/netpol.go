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

func CreateEgressPolicy(clientSet kubernetes.Clientset, namespace string) error {
	policyName := "egress-policy"
	policyTypes := []networking.PolicyType{"Egress"}
	egress := buildEgressRules()
	matchLabels := make(map[string]string)
	matchLabels[utils.NetworkPolicyLabelKey] = utils.NetworkPolicyLabelValue
	return createNetworkPolicy(clientSet, policyName, namespace, policyTypes, egress, nil, matchLabels)
}

func CreateChallengeEgressPolicy(clientSet kubernetes.Clientset, namespace string) error {
	policyName := "challenge-egress-policy"
	policyTypes := []networking.PolicyType{"Egress"}
	matchLabels := make(map[string]string)
	matchLabels[utils.ChallengePodLabelKey] = utils.ChallengePodLabelValue
	return createNetworkPolicy(clientSet, policyName, namespace, policyTypes, nil, nil, matchLabels)
}

func CreateChallengeIngressPolicy(clientSet kubernetes.Clientset, namespace string) error {
	policyName := "ingress-policy"
	policyTypes := []networking.PolicyType{"Ingress"}

	ingress := buildIngressRules()
	matchLabels := make(map[string]string)
	matchLabels[utils.ChallengePodLabelKey] = utils.ChallengePodLabelValue
	return createNetworkPolicy(clientSet, policyName, namespace, policyTypes, nil, ingress, matchLabels)
}

func buildEgressRules() []networking.NetworkPolicyEgressRule {
	getAddress := func(s v1.Protocol) *v1.Protocol { return &s }
	return []networking.NetworkPolicyEgressRule{
		{
			To: []networking.NetworkPolicyPeer{
				{
					PodSelector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							utils.ChallengePodLabelKey: utils.ChallengePodLabelValue,
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
							utils.KubernetesDNSLabelKey: utils.KubernetesDNSLabelValue,
						},
					},
				},
			},
			Ports: []networking.NetworkPolicyPort{
				{
					Port:     &intstr.IntOrString{Type: intstr.Int, IntVal: 53},
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
							utils.KaliPodLabelKey: utils.KaliPodLabelValue,
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
							utils.WireguardPodLabelKey: utils.WireguardPodLabelValue,
						},
					},
				},
			},
		},
	}
}

func createNetworkPolicy(clientSet kubernetes.Clientset, policyName string, namespace string, policyTypes []networking.PolicyType, egress []networking.NetworkPolicyEgressRule, ingress []networking.NetworkPolicyIngressRule, matchLabels map[string]string) error {
	netpol := configureNetworkPolicy(policyName, namespace, policyTypes, egress, ingress, matchLabels)
	networkClient := clientSet.NetworkingV1().NetworkPolicies(namespace)
	result, err := networkClient.Create(context.TODO(), &netpol, metav1.CreateOptions{})

	if err != nil {
		return err
	}

	utils.InfoLogger.Printf("Created network policy of type %q with name %q for namespace %s\n", &result.Spec.PolicyTypes, result.GetObjectMeta().GetName(), namespace)
	return nil
}

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
