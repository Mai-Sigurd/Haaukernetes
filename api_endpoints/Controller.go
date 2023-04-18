package api_endpoints

import "k8s.io/client-go/kubernetes"

type Controller struct {
	ClientSet *kubernetes.Clientset
}

type ErrorResponse struct {
	Message string
}
