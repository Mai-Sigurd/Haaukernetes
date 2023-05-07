package api

import (
	"bytes"
	"fmt"
	"io"
	"k8s-project/api_endpoints"
	"k8s-project/utils"
	"net/http"

	"github.com/goccy/go-json"
)

var ipPort = utils.APIPort

func SetIpPort(port string) {
	ipPort = port
}

func GetUser(name string) {
	reqBody := api_endpoints.User{Name: name}
	jsonBody, _ := json.Marshal(reqBody)
	url := "http://localhost:" + ipPort + "/user/"
	resp, err := http.NewRequest("GET", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Println(err.Error())
	}
	generalResponse(resp)
}

func PostUser(name string) {
	reqBody := api_endpoints.User{Name: name}
	jsonBody, _ := json.Marshal(reqBody)

	url := "http://localhost:" + ipPort + "/user/"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Println(err.Error())
	}
	generalResponse(req)

}

func DeleteUser(name string) {
	reqBody := api_endpoints.User{Name: name}
	jsonBody, _ := json.Marshal(reqBody)

	url := "http://localhost:" + ipPort + "/user/"
	resp, err := http.NewRequest("DELETE", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Println(err.Error())
	}
	generalResponse(resp)
}
func GetUsers() {
	url := "http://localhost:" + ipPort + "/users/"
	resp, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err.Error())
	}
	generalResponse(resp)
}

func GetUserChallenges(name string) {
	reqBody := api_endpoints.User{Name: name}
	jsonBody, _ := json.Marshal(reqBody)

	url := "http://localhost:" + ipPort + "/user/challenges/"
	resp, err := http.NewRequest("GET", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Println(err.Error())
	}
	generalResponse(resp)
}

func PostKali(username string, password string) {
	reqBody := api_endpoints.GuacUser{Name: username, Password: password}
	jsonBody, _ := json.Marshal(reqBody)
	url := "http://localhost:" + ipPort + "/kali/"
	resp, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Println(err.Error())
	}
	generalResponse(resp)
}

func PostWireguard(username string, key string) {
	reqBody := api_endpoints.Wireguard{User: username, Key: key}
	jsonBody, _ := json.Marshal(reqBody)
	fmt.Println(reqBody)
	url := "http://localhost:" + ipPort + "/wireguard/"
	resp, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Println(err.Error())
	}
	generalResponse(resp)
}

func DeleteChallenge(user string, challengeName string) {
	reqBody := api_endpoints.DelChallenge{
		ChallengeName: challengeName,
		User:          user,
	}
	jsonBody, _ := json.Marshal(reqBody)
	url := "http://localhost:" + ipPort + "/challenge/"
	resp, err := http.NewRequest("DELETE", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Println(err.Error())
	}
	generalResponse(resp)
}

func PostChallenge(user string, challengeName string, ports []int32) {
	reqBody := api_endpoints.Challenge{
		Ports:         ports,
		ChallengeName: challengeName,
		User:          user,
	}
	jsonBody, _ := json.Marshal(reqBody)

	url := "http://localhost:" + ipPort + "/challenge/"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Println(err.Error())
	}
	generalResponse(req)

}

func generalResponse(req *http.Request) {
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Print(err.Error())
	}
	defer res.Body.Close()
	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		fmt.Print(err.Error())
	}
	fmt.Println(string(body))
}
