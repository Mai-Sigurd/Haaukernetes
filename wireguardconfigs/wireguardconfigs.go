package wireguardconfigs

import (
	"context"
	"k8-project/utils"
	"regexp"
	"strconv"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const clientConfig = `
[Interface]
Address = 10.33.0.2/32
PrivateKey =
DNS =

[Peer]
PublicKey =
Endpoint =
AllowedIPs =
PersistentKeepalive = 25
`

const serverConfig = `
[Interface]
Address = 10.33.0.1/24
ListenPort = 51820
PrivateKey =
PostUp = iptables -t nat -A POSTROUTING -s 10.33.0.0/24 -o eth0 -j MASQUERADE
PostDown = iptables -t nat -D POSTROUTING -s 10.33.0.0/24 -o eth0 -j MASQUERADE

[Peer]
PublicKey =
AllowedIPs = 10.33.0.2/32
`

func GetClientConfig(clientSet kubernetes.Clientset, serverPublicKey string, nodeport int32, endpoint string, subnet string) string {
	configWithIPsAndEndpoint := addAllowedIpsAndEndpointToClientConfig(clientSet, addNodePort(nodeport, endpoint), subnet)
	return replacePublicKey(serverPublicKey, configWithIPsAndEndpoint)
}

func GetServerConfig(privateKey string, publicKey string) string {
	return addKeysToConfig(privateKey, publicKey, serverConfig)
}

func addAllowedIpsAndEndpointToClientConfig(clientSet kubernetes.Clientset, endpoint string, subnet string) string {
	endpointRegex := regexp.MustCompile("Endpoint =")
	endpointReplaced := endpointRegex.ReplaceAllString(clientConfig, "Endpoint = "+endpoint)
	allowedIpsRegex := regexp.MustCompile("AllowedIPs =")
	allowedIpsReplaced := allowedIpsRegex.ReplaceAllString(endpointReplaced, "AllowedIPs = "+subnet)

	dnsIPReplaced := replaceKubeDNSIP(clientSet, allowedIpsReplaced)
	return dnsIPReplaced
}

func addKeysToConfig(privateKey string, publicKey string, conf string) string {
	privateKeyReplaced := replacePrivateKey(privateKey, conf)
	bothKeysReplaced := replacePublicKey(publicKey, privateKeyReplaced)
	return bothKeysReplaced
}

func replacePrivateKey(privateKey string, conf string) string {
	privateKeyRegex := regexp.MustCompile("PrivateKey =")
	return privateKeyRegex.ReplaceAllString(conf, "PrivateKey = "+privateKey)
}

func replacePublicKey(publicKey string, conf string) string {
	publicKeyRegex := regexp.MustCompile("PublicKey =")
	return publicKeyRegex.ReplaceAllString(conf, "PublicKey = "+publicKey)
}

func addNodePort(nodePort int32, endpoint string) string {
	return endpoint + strconv.Itoa(int(nodePort))
}

func getKubeDnsIP(clientSet kubernetes.Clientset) string {
	dnsService, err := clientSet.CoreV1().Services("kube-system").Get(context.TODO(), "kube-dns", metav1.GetOptions{})
	utils.ErrLogger(err)
	return dnsService.Spec.ClusterIP
}

func replaceKubeDNSIP(clientSet kubernetes.Clientset, conf string) string {
	DNSRegex := regexp.MustCompile("DNS =")
	return DNSRegex.ReplaceAllString(conf, "DNS = "+getKubeDnsIP(clientSet))
}
