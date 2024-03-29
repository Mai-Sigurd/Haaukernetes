package wireguard

import (
	"k8s-project/configmap"
	"k8s-project/connections/vpn/wireguardconfig"
	"k8s-project/deployments"
	"k8s-project/secrets"
	"k8s-project/services"
	"k8s-project/utils"
	"os/exec"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
)

// StartWireguard
// clientpublickey should come from caller i.e. api call
// clientprivatekey should be inserted to file by client itself
func StartWireguard(clientSet kubernetes.Clientset, namespace string, clientPublicKey string, endpoint string, subnet string) (string, error) {
	serverPrivateKey, serverPublicKey := createKeys()
	err := configmap.CreateWireGuardConfigMap(clientSet, namespace, serverPrivateKey, clientPublicKey)
	if err != nil {
		utils.ErrLogger(err)
		return "", err
	}
	err = secrets.CreateWireGuardSecret(clientSet, namespace, serverPrivateKey)
	if err != nil {
		utils.ErrLogger(err)
		return "", err
	}
	deployment := configureWireGuardDeployment()
	err = deployments.CreatePrebuiltDeployment(clientSet, namespace, deployment)
	if err != nil {
		utils.ErrLogger(err)
		return "", err
	}
	service := configureWireguardNodePortService(namespace)
	createdService, err := services.CreatePrebuiltService(clientSet, namespace, *service)
	if err != nil {
		utils.ErrLogger(err)
		return "", err
	}
	clientConf := wireguardconfig.GetClientConfig(clientSet, serverPublicKey, createdService.Spec.Ports[0].NodePort, endpoint, subnet)

	utils.InfoLogger.Printf("Wireguard successfully started for user: %s\n", namespace)
	return clientConf, nil
}

func createKeys() (string, string) {
	priv, err := exec.Command("/bin/sh", "-c", "docker run --rm -i masipcat/wireguard-go wg genkey").Output()
	utils.ErrLogger(err)
	pub, err := exec.Command("/bin/sh", "-c", "echo '"+string(priv)+"' | docker run --rm -i masipcat/wireguard-go wg pubkey").Output()
	utils.ErrLogger(err)
	return string(priv), string(pub)
}

func configureWireguardNodePortService(namespace string) *apiv1.Service {
	service := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      utils.WireguardPodLabelValue,
			Namespace: namespace,
			Labels: map[string]string{
				utils.WireguardPodLabelKey: utils.WireguardPodLabelValue,
			},
		},
		Spec: apiv1.ServiceSpec{
			Type: apiv1.ServiceTypeNodePort,
			Ports: []apiv1.ServicePort{
				{
					Name:       "wg",
					Protocol:   apiv1.ProtocolUDP,
					Port:       51820,
					TargetPort: intstr.FromInt(51820),
				},
			},
			Selector: map[string]string{
				utils.WireguardPodLabelKey: utils.WireguardPodLabelValue,
			},
			ClusterIP: "",
		},
	}
	return service
}

func configureWireGuardDeployment() *appsv1.Deployment {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: utils.WireguardPodLabelValue,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: utils.Int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					utils.WireguardPodLabelKey:  utils.WireguardPodLabelValue,
					utils.NetworkPolicyLabelKey: utils.NetworkPolicyLabelValue,
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						utils.WireguardPodLabelKey:  utils.WireguardPodLabelValue,
						utils.NetworkPolicyLabelKey: utils.NetworkPolicyLabelValue,
					},
				},
				Spec: apiv1.PodSpec{
					InitContainers: []apiv1.Container{
						{
							Name:    "sysctls",
							Image:   "busybox",
							Command: []string{"sh", "-c", "sysctl -w net.ipv4.ip_forward=1 && sysctl -w net.ipv4.conf.all.forwarding=1"},
							SecurityContext: &apiv1.SecurityContext{
								Capabilities: &apiv1.Capabilities{
									Add: []apiv1.Capability{
										"NET_ADMIN",
									},
								},
								Privileged: utils.BoolPtr(true),
							},
						},
					},
					Containers: []apiv1.Container{
						{
							Name:    utils.WireguardPodLabelValue,
							Image:   utils.ImageRepoUrl + utils.WireguardImage,
							Command: []string{"sh", "-c", "echo 'Public key '$(wg pubkey < /etc/wireguard/privatekey)'' && /entrypoint.sh"},
							Ports: []apiv1.ContainerPort{
								{
									ContainerPort: 51820,
									Protocol:      apiv1.ProtocolUDP,
									Name:          utils.WireguardPodLabelValue,
								},
							},
							Env: []apiv1.EnvVar{
								{
									Name:  "LOG_LEVEL",
									Value: "info",
								},
							},
							SecurityContext: &apiv1.SecurityContext{
								Capabilities: &apiv1.Capabilities{
									Add: []apiv1.Capability{
										"NET_ADMIN",
									},
								},
							},
							VolumeMounts: []apiv1.VolumeMount{
								{
									Name:      "cfgmap",
									MountPath: "/etc/wireguard/wg0.conf",
									SubPath:   "wg0.conf",
								},
								{
									Name:      "secret",
									MountPath: "/etc/wireguard/privatekey",
									SubPath:   "privatekey",
								},
							},
						},
					},
					ImagePullSecrets: []apiv1.LocalObjectReference{
						{
							Name: "regcred",
						},
					},
					Volumes: []apiv1.Volume{
						{
							Name: "cfgmap",
							VolumeSource: apiv1.VolumeSource{
								ConfigMap: &apiv1.ConfigMapVolumeSource{
									LocalObjectReference: apiv1.LocalObjectReference{
										Name: "wg-configmap",
									},
								},
							},
						},
						{
							Name: "secret",
							VolumeSource: apiv1.VolumeSource{
								Secret: &apiv1.SecretVolumeSource{
									SecretName: "wg-secret",
								},
							},
						},
					},
				},
			},
		},
	}
	return deployment
}
