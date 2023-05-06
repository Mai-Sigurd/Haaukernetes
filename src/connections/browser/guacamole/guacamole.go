package guacamole

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"k8s-project/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func GetGuacamolePasswordSecret(clientSet kubernetes.Clientset) (string, error) {
	secret, err := clientSet.CoreV1().Secrets("guacamole").Get(context.TODO(), "guacamole", metav1.GetOptions{})
	password := string(secret.Data["guac-password"])
	utils.InfoLogger.Printf("Retrieved Guacamole secret")
	return password, err
}

func GetGuacamoleBaseAddress(clientSet kubernetes.Clientset) (string, error) {
	serverIp := os.Getenv("SERVER_IP")
	guacamoleService, err := utils.FindService(clientSet, "guacamole", "guacamole")
	port := guacamoleService.Spec.Ports[0].NodePort
	baseAddress := fmt.Sprintf("http://%s:%d/guacamole", serverIp, port)
	utils.InfoLogger.Printf("Retrieved Guacamole base address: " + baseAddress)
	return baseAddress, err
}

func (guac *Guacamole) UpdateAuthToken() error {
	form := url.Values{
		"username":   {guac.Username},
		"password":   {guac.Password},
		"attributes": {},
	}
	formData := form.Encode()

	req, err := http.NewRequest("POST", guac.BaseUrl+"/api/tokens", strings.NewReader(formData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New("could not update Guacamole auth token: " + string(body))
	}

	var responseMap map[string]interface{}

	err = json.Unmarshal(body, &responseMap)
	if err != nil {
		return err
	}

	authToken := responseMap["authToken"].(string)
	guac.AuthToken = authToken
	utils.InfoLogger.Printf("Updated Guacamole auth token")
	return nil
}

func (guac *Guacamole) UpdateDefaultGuacAdminPassword(clientSet kubernetes.Clientset, oldPassword string) error {
	err := guac.UpdateAuthToken()
	if err != nil {
		return err
	}

	guacPassword, err := GetGuacamolePasswordSecret(clientSet)
	guac.Password = guacPassword
	if err != nil {
		return err
	}

	u := UpdateUser{
		OldPassword: oldPassword,
		NewPassword: guac.Password,
	}

	payload, err := json.Marshal(u)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", guac.BaseUrl+"/api/session/data/postgresql/users/"+guac.Username+"/password?token="+guac.AuthToken, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 204 {
		return errors.New("could not update Guacamole password for user " + guac.Username + ": " + string(body))
	}

	utils.InfoLogger.Printf("Updated Guacamole admin password for user" + guac.Username)
	return nil
}

func (guac *Guacamole) CreateUser(username string, password string) error {
	u := UserInfo{
		Username:   username,
		Password:   password,
		Attributes: CreateUserAttributes{},
	}

	payload, err := json.Marshal(u)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", guac.BaseUrl+"/api/session/data/postgresql/users?token="+guac.AuthToken, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New("could not create user " + username + ": " + string(respBody))
	}

	utils.InfoLogger.Printf("Created Guacamole user " + username + ": " + string(respBody))
	return nil
}

func (guac *Guacamole) CreateConnection(kaliIp string, kaliPort int32, username string) (string, error) {
	param := RDPParameters{
		Username:   "kali",
		Password:   "kali",
		Hostname:   kaliIp,
		Port:       strconv.Itoa(int(kaliPort)),
		IgnoreCert: true,
	}

	attr := RDPAttributes{}

	conn := RDPConnection{
		ParentIdentifier: "ROOT",
		Name:             fmt.Sprintf("kali-%s-%s", username, kaliIp),
		Protocol:         "rdp",
		Parameters:       param,
		Attributes:       attr,
	}

	payload, err := json.Marshal(conn)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", guac.BaseUrl+"/api/session/data/postgresql/connections?token="+guac.AuthToken, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", errors.New("could not create connection for user " + username + ": " + string(body))
	}

	var responseMap map[string]interface{}

	err = json.Unmarshal(body, &responseMap)
	if err != nil {
		return "", err
	}

	identifier := responseMap["identifier"].(string)
	utils.InfoLogger.Printf("Created Kali RDP connection for user: " + username)

	return identifier, nil
}

func (guac *Guacamole) AssignConnection(connIdentifier string, username string) error {
	addConn := []AddConnection{{
		Operation: "add",
		Path:      fmt.Sprintf("/connectionPermissions/%s", connIdentifier),
		Value:     "READ",
	}}

	payload, err := json.Marshal(addConn)
	if err != nil {
		return nil
	}

	patchUrl := fmt.Sprintf("%s/api/session/data/postgresql/users/%s/permissions?token=%s", guac.BaseUrl, username, guac.AuthToken)
	fmt.Println(patchUrl)

	req, err := http.NewRequest("PATCH", patchUrl, bytes.NewBuffer(payload))

	if err != nil {
		return nil
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	utils.InfoLogger.Printf("Assigned Kali connection to user " + username + ": " + string(body))
	return nil
}
