package apis

import "k8s.io/client-go/kubernetes"

type Controller struct {
	ClientSet *kubernetes.Clientset
}
