package guacamole

import ( // todo we need to change default guac user somehow to not have it exposed to the whole world
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func (guac *Guacamole) getAuthToken() (string, error) {
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
	return authToken, nil // TODO maybe save the access token inside input guac struct and return that one???
}

func (guac *Guacamole) createUser(u UserInfo) {

	payload, err := json.Marshal(u)
	if err != nil {
		fmt.Println("Error marshaling JSON payload:", err) // TODO error handling
		return
	}

	req, err := http.NewRequest("POST", guac.BaseUrl+"/api/session/data/postgresql/users?token="+guac.AuthToken, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending HTTP request:", err) // TODO error handling
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err) // TODO error handling
		return
	}

	fmt.Println("Response Body:", string(body)) // TODO error handling and do something with it
}

func (guac *Guacamole) createConnection(kaliHostname string, kaliPort string) {
	param := RDPParameters{
		Username:   "Kali",
		Password:   "Kali",
		Hostname:   kaliHostname,
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
		return
	}

	req, err := http.NewRequest("POST", guac.BaseUrl+"/api/session/data/postgresql/connections?token="+guac.AuthToken, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending HTTP request:", err) // TODO error handling
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err) // TODO error handling
		return
	}

	var responseMap map[string]interface{}

	err = json.Unmarshal(body, &responseMap)
	if err != nil {
		fmt.Println("Error decoding response body:", err) // TODO error handling
		return
	}

	identifier := responseMap["identifier"].(string)

	fmt.Println("Identifier", identifier) // TODO error handling and do something with it
}
