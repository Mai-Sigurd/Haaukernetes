package guacamole

import ( // todo we need to change default guac user somehow to not have it exposed to the whole world
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"k8s-project/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func GetGuacamoleSecret(clientSet kubernetes.Clientset) (string, string, error) {
	secret, err := clientSet.CoreV1().Secrets("guacamole").Get(context.TODO(), "guacamole", metav1.GetOptions{}) // TODO HANDLE ERROR
	username := string(secret.Data["guac-user"])
	password := string(secret.Data["guac-password"])
	return username, password, err
}

func GetGuacamoleBaseAddress(clientSet kubernetes.Clientset) string {
	serverIp := os.Getenv("SERVER_IP")
	guacamoleService, _ := utils.FindService(clientSet, "guacamole", "guacamole") // TODO HANDLE ERROR
	port := guacamoleService.Spec.Ports[0].NodePort                               // TODO is it in the form something:something? Then this might not work
	return fmt.Sprintf("http://%s:%d/guacamole", serverIp, port)
}

func (guac *Guacamole) GetAuthToken() (string, error) {
	form := url.Values{
		"username":   {guac.Username},
		"password":   {guac.Password},
		"attributes": {},
	}
	formData := form.Encode()

	req, err := http.NewRequest("POST", guac.BaseUrl+"/api/tokens", strings.NewReader(formData))
	if err != nil {
		fmt.Println("Error creating request:", err) // TODO error handling
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err) // TODO error handling
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err) // TODO error handling
		return "", err
	}

	var responseMap map[string]interface{}

	err = json.Unmarshal(body, &responseMap)
	if err != nil {
		fmt.Println("Error decoding response body:", err) // TODO error handling
		return "", err
	}

	authToken := responseMap["authToken"].(string)
	guac.AuthToken = authToken
	return "", nil // TODO maybe save the access token inside input guac struct and return that one???
}

func (guac *Guacamole) CreateUser(username string, password string) error {

	u := UserInfo{
		Username:   username,
		Password:   password,
		Attributes: CreateUserAttributes{},
	}

	payload, err := json.Marshal(u)
	if err != nil {
		fmt.Println("Error marshaling JSON payload:", err) // TODO error handling
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
		fmt.Println("Error sending HTTP request:", err) // TODO error handling
		return err
	}
	defer resp.Body.Close()

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err) // TODO error handling
		return err
	}

	return nil
}

func (guac *Guacamole) CreateConnection(kaliIp string, kaliPort string) (string, error) {
	param := RDPParameters{
		Username:   "Kali",
		Password:   "Kali",
		Hostname:   kaliIp,
		Port:       kaliPort,
		IgnoreCert: true,
	}

	attr := RDPAttributes{}

	conn := RDPConnection{
		ParentIdentifier: "ROOT",
		Name:             "Kali-RDP",
		Protocol:         "rdp",
		Parameters:       param,
		Attributes:       attr,
	}

	payload, err := json.Marshal(conn)
	if err != nil {
		fmt.Println("Error marshaling JSON payload:", err) // TODO error handling
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
		fmt.Println("Error sending HTTP request:", err) // TODO error handling
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err) // TODO error handling
		return "", err
	}

	var responseMap map[string]interface{}

	err = json.Unmarshal(body, &responseMap)
	if err != nil {
		fmt.Println("Error decoding response body:", err) // TODO error handling
		return "", err
	}

	identifier := responseMap["identifier"].(string)

	return identifier, nil // TODO error handling and do something with it
}

func (guac *Guacamole) AssignConnection(connIdentifier string, username string) (string, error) {
	addConn := []AddConnection{{
		Operation: "add",
		Path:      fmt.Sprintf("/connectionPermissions/%s", connIdentifier),
		Value:     "READ",
	}}

	payload, err := json.Marshal(addConn)
	if err != nil {
		fmt.Println("Error marshaling JSON payload:", err) // TODO error handling
		return "", nil
	}

	patchUrl := fmt.Sprintf("%s/api/session/data/postgresql/users/%s/permissions?token=%s", guac.BaseUrl, username, guac.AuthToken)
	fmt.Println(patchUrl)

	req, err := http.NewRequest("PATCH", patchUrl, bytes.NewBuffer(payload))

	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return "", nil
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending HTTP request:", err) // TODO error handling
		return "", nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err) // TODO error handling
		return "", nil
	}

	fmt.Println("Response Body:", string(body)) // TODO error handling and do something with it
	return "", nil
}
