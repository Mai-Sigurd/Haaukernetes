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

type guacamole struct {
	Username  string
	Password  string
	BaseUrl   string // vi skal bruge user, pass, server ip, guac port somehow
	AuthToken string
}

type createUserAttributes struct {
	Disabled          string `json:"disabled"`
	Expired           string `json:"expired"`
	AccessWindowStart string `json:"access-window-start"`
	AccessWindowEnd   string `json:"access-window-end"`
	ValidFrom         string `json:"valid-from"`
	ValidUntil        string `json:"valid-until"`
	TimeZone          string `json:"timezone"`
}

type userInfo struct {
	Username   string               `json:"username"`
	Password   string               `json:"password"`
	Attributes createUserAttributes `json:"attributes"`
}

func (guac *guacamole) getAuthToken() (string, error) {
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

func (guac *guacamole) createUser(username string, password string) {
	creds := userInfo{
		Username: username,
		Password: password,
	}

	payload, err := json.Marshal(creds)
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

	fmt.Println("Response Body:", string(body)) // TODO error handling
}
