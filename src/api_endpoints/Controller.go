package api_endpoints

import (
	"k8s-project/connections/browser/guacamole"
	"k8s.io/client-go/kubernetes"
)

type Controller struct {
	ClientSet *kubernetes.Clientset
	Guacamole guacamole.Guacamole
}

type ErrorResponse struct {
	Message string
}

type User struct {
	// Username
	// in: string
	Name string `json:"name"`
}
