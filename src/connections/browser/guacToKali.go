package browser

import (
	"k8s-project/connections/browser/guacamole"
	"k8s-project/connections/browser/kali"
	"k8s-project/utils"
	"k8s.io/client-go/kubernetes"
)

func SetupGuacAndKali(guac guacamole.Guacamole, clientSet kubernetes.Clientset, name string, password string) (string, error) {
	ip, port, err := kali.StartKali(clientSet, name, "kali")

	err = guac.UpdateAuthToken()
	if err != nil {
		utils.ErrLogger(err)
		return "", err
	}

	err = guac.CreateUser(name, password)
	if err != nil {
		utils.ErrLogger(err)
		return "", err
	}

	connIdentifier, err := guac.CreateConnection(ip, port, name)
	if err != nil {
		utils.ErrLogger(err)
		return "", err
	}

	err = guac.AssignConnection(connIdentifier, name)
	if err != nil {
		utils.ErrLogger(err)
		return "", err
	}

	message := "You can now RDP into your Kali by visiting the Guacamole interface at: " + guac.BaseUrl
	utils.InfoLogger.Printf("Message to user " + name + ": " + message)
	return message, nil
}
