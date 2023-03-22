package configmap

import (
	"context"
	"fmt"
	"k8-project/utils"
	"regexp"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var clientConf = `
	[Interface]
	# Assign you an IP (that's not in use) and add it to server configmap
	Address = 10.33.0.2/32
	PrivateKey =
	#DNS = 10.96.0.10
	
	[Peer]
	# Wireguard server public key
	PublicKey =
	Endpoint = 164.92.194.69:31820
	AllowedIPs = 
	PersistentKeepalive = 25
	`

//spaghetti -> should this be a secret of some kind? maybe it's own file?
//also: will address-subnet have to be dynamic somehow?
//also: commented out postup line is commented out because i couldnt get the path to privatekey working, instead it is "hardcoded"
var serverConf = `
	[Interface]
	Address = 10.33.0.1/24
	ListenPort = 51820
	PrivateKey =
	PostUp = iptables -t nat -A POSTROUTING -s 10.33.0.0/24 -o eth0 -j MASQUERADE
	#PostUp = wg set wg0 private-key /etc/wireguard/privatekey && iptables -t nat -A POSTROUTING -s 10.33.0.0/24 -o eth0 -j MASQUERADE
	PostDown = iptables -t nat -D POSTROUTING -s 10.33.0.0/24 -o eth0 -j MASQUERADE

	[Peer]
	PublicKey =
	AllowedIPs = 10.33.0.2/32
	`

func CreateWireGuardConfigMap(clientSet kubernetes.Clientset, teamName string, serverPrivateKey string, clientPublicKey string) {
	data := make(map[string]string)
	data["wg0.conf"] = generateConfig(serverPrivateKey, clientPublicKey)
	configMap := configureConfigMap("wg-configmap", teamName, data)
	CreateConfigMap(clientSet, teamName, configMap)
}

func CreateConfigMap(clientSet kubernetes.Clientset, teamName string, configMap v1.ConfigMap) {
	configMapClient := clientSet.CoreV1().ConfigMaps(teamName)
	result, err := configMapClient.Create(context.TODO(), &configMap, metav1.CreateOptions{})
	utils.ErrHandler(err)
	fmt.Printf("Created configmap %q.\n", result.GetObjectMeta().GetName())
}

func configureConfigMap(name string, namespace string, data map[string]string) v1.ConfigMap {
	configmap := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Data: data,
	}

	return *configmap
}
func addAllowedIps(endpoint string, subnet string) string {
	endpointRegex := regexp.MustCompile("Endpoint =")
	endpointReplaced := endpointRegex.ReplaceAllString(clientConf, "Endpoint ="+endpoint)
	allowedIpsRegex := regexp.MustCompile("AllowedIPs =")
	allowedIpsReplaced := allowedIpsRegex.ReplaceAllString(endpointReplaced, "AllowedIPs ="+subnet)
	return allowedIpsReplaced
}

func generateConfig(serverPrivateKey string, clientPublicKey string) string {
	privateKeyReplaced := replacePrivateKey(serverPrivateKey, serverConf)
	bothKeysReplaced := replacePublicKey(clientPublicKey, privateKeyReplaced)
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
