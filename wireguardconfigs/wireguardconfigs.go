package wireguardconfigs

import (
	"regexp"
	"strconv"
)

// public endpoint for k8s cluster
const endpoint = "164.92.194.69:"

// pod CIDR for cluster
const subnet = "10.96.0.0/12"

const clientConfig = `
[Interface]
# Assign you an IP (that's not in use) and add it to server configmap
Address = 10.33.0.2/32
PrivateKey =
#DNS = 10.96.0.10

[Peer]
# Wireguard server public key
PublicKey =
Endpoint =
AllowedIPs =
PersistentKeepalive = 25
`

// spaghetti -> should this be a secret of some kind? maybe it's own file? TODO
// also: will address-subnet have to be dynamic somehow?
// also: commented out postup line is commented out because i couldnt get the path to privatekey working, instead it is "hardcoded"
const serverConfig = `
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

func GetClientConfig(serverPublicKey string, nodeport int32) string {
	configWithIPsAndEndpoint := addAllowedIpsAndEndpointToClientConfig(addNodePort(nodeport, endpoint), subnet)
	return replacePublicKey(serverPublicKey, configWithIPsAndEndpoint)
}

func GetServerConfig(privateKey string, publicKey string) string {
	return addKeysToConfig(privateKey, publicKey, serverConfig)
}

// Todo Purpose?
func GetEndpoint() string {
	return endpoint
}

// Todo Purpose?
func GetSubnet() string {
	return subnet
}

func addAllowedIpsAndEndpointToClientConfig(endpoint string, subnet string) string {
	endpointRegex := regexp.MustCompile("Endpoint =")
	endpointReplaced := endpointRegex.ReplaceAllString(clientConfig, "Endpoint = "+endpoint)
	allowedIpsRegex := regexp.MustCompile("AllowedIPs =")
	allowedIpsReplaced := allowedIpsRegex.ReplaceAllString(endpointReplaced, "AllowedIPs = "+subnet)
	return allowedIpsReplaced
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
