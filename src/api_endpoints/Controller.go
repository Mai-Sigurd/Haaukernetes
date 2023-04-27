package api_endpoints

import "k8s.io/client-go/kubernetes"

type Controller struct {
	ClientSet *kubernetes.Clientset
	Endpoint  string
	Subnet    string
}

type ErrorResponse struct {
	Message string
}

type User struct {
	// Username
	// in: string
	Name string `json:"name"`
}
