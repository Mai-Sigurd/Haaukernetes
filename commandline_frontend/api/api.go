package api

import (
	"bytes"
	"fmt"
	"io"
	"k8-project/api_endpoints"
	"net/http"

	"github.com/goccy/go-json"
)

var ipPort = "5113"

func SetIpPort(port string) {
	ipPort = port
}

func GetUser(name string) {
	url := "http://localhost:" + ipPort + "/namespace/" + name
	resp, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err.Error())
	}
	generalResponse(resp)
}

func PostUser(name string) {
	reqBody := api_endpoints.User{Name: name}
	jsonBody, _ := json.Marshal(reqBody)

	url := "http://localhost:" + ipPort + "/namespace/"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Println(err.Error())
	}
	generalResponse(req)

}

func DeleteUser(name string) {
	url := "http://localhost:" + ipPort + "/namespace/" + name
	resp, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		fmt.Println(err.Error())
	}
	generalResponse(resp)
}

func GetKali(namespace string) {
	url := "http://localhost:" + ipPort + "/kali/" + namespace
	resp, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err.Error())
	}
	generalResponse(resp)
}

func PostKali(namespace string) {
	url := "http://localhost:" + ipPort + "/kali/" + namespace
	resp, err := http.NewRequest("POST", url, nil)
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
