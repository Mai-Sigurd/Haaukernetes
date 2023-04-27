package utils

const (
	Port                    = "33333"
	ImageRepoUrl            = "registry.digitalocean.com/haaukins-bsc/"
	WireguardImage          = "wireguard-go"
	WireguardEndpoint       = "164.92.194.69" // public endpoint for k8s cluster
	WireguardSubnet         = "10.96.0.0/12"  // pod CIDR for cluster
	WireguardPodLabelKey    = "vpn"
	WireguardPodLabelValue  = "wireguard"
	KaliPodLabelKey         = "app"
	KaliPodLabelValue       = "kali"
	ChallengePodLabelKey    = "type"
	ChallengePodLabelValue  = "challenge"
	KubernetesDNSLabelKey   = "k8s-app"
	KubernetesDNSLabelValue = "kube-dns"
)
