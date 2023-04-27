package guacamole

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type guacamole struct {
	username string
	password string
	baseUrl  string // vi skal bruge user, pass, server ip, guac port somehow
}

func (guac *guacamole) getAccessToken() (string, error) {
	form := url.Values{"username": {guac.username}, "password": {guac.password}}
	formData := form.Encode()

	req, err := http.NewRequest("POST", guac.baseUrl+"/guacamole/api/tokens", strings.NewReader(formData))
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

	body, err := ioutil.ReadAll(resp.Body)
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
	return authToken, nil // TODO maybe save the access token incide input guac struct and return that one???
}
