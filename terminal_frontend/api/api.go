package api

import (
	"bytes"
	"fmt"
	"github.com/goccy/go-json"
	"k8-project/apis"
	"net/http"
)

var ipPort = "5113"

func SetIpPort(port string) {
	ipPort = port
}

func GetNamespace(name string) {
	url := "http://localhost:" + ipPort + "/namespace/" + name
	resp, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err.Error())
	}
	generalResponse(resp)
}

func PostNamespace(name string) {
	reqBody := apis.Namespace{Name: name}
	jsonBody, _ := json.Marshal(reqBody)

	url := "http://localhost:" + ipPort + "/namespace/"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Println(err.Error())
	}
	generalResponse(req)

}

func DeleteNamespace(name string) {
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

func DeleteChallenge(namespace string, challengeName string) {
	reqBody := apis.DelChallenge{
		ChallengeName: challengeName,
		Namespace:     namespace,
	}
	jsonBody, _ := json.Marshal(reqBody)
	url := "http://localhost:" + ipPort + "/challenge/"
	resp, err := http.NewRequest("DELETE", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Println(err.Error())
	}
	generalResponse(resp)
}

func PostChallenge(namespace string, challengeName string, port int32) {
	reqBody := apis.Challenge{
		Port:          port,
		ChallengeName: challengeName,
		Namespace:     namespace,
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
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		fmt.Print(err.Error())
	}
	fmt.Println(string(body))
}
