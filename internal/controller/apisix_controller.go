package controller

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"os"

)

type Upstream struct{
	Node struct{
		Nodes []struct{
			value struct{
				Name string `json:"name"`
				NodePort string `json:"nodes"`
			}`json:"value"`
		}`json:"nodes"`
	}`json:"node"`
}
var (
	host,apiKey ="",""
)
var existUpstreams = make(map[string]string,0)

func init(){
	env := os.Getenv("in_host")
	getUrl(env)

}
func getUrl(env string){

	if env == "beta" {

	}else{

	}
}

func UrlRequest(method, url string, payload interface{}) (string, error) {
	var req *http.Request
	var err error
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-KEY",apiKey)
	if method == "POST" {
		jsonData, err := json.Marshal(payload)
		if err != nil {
			return "", err
		}
		req, err = http.NewRequest("POST", url, bytes.NewBuffer(jsonData))

	} else if method == "GET" {
		req, err = http.NewRequest("GET", url, nil)
	} else {
		return "", fmt.Errorf("unsupported method: %s", method)
	}

	if err != nil {
		return "", err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}


func IsExistStreams(){
	url := host+"/apisix/admin/upstreams"


	
}

func AddUpstreams(){

}
