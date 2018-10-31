package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var restURL = "https://httpbin.org"
var marmotURL = "https://admin.staging.bluemix.net:7663/api"

func main() {
	// sendGet()
	// sendPost1()
	// sendPost2()
	sendBasicAuth()
	//loginMarmot()
}

func loginMarmot() {
	// jsonData := map[string]string{"firstname": "Nic", "lastname": "Raboy"}
	// jsonValue, _ := json.Marshal(jsonData)
	authString := "/login/ibmid"
	restEndpoint := marmotURL + authString
	request, _ := http.NewRequest("POST", restEndpoint, bytes.NewBuffer(nil))
	request.Header.Set("Content-Type", "application/json")
	request.SetBasicAuth("Stefan.Zink@de.ibm.com", "i1Oetz.i")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		// data, _ := ioutil.ReadAll(response.Body)
		// fmt.Println(string(data))

		// dump, err := httputil.DumpResponse(response, true)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// fmt.Printf("%q", dump)

		fmt.Println(response.StatusCode)
		log.Println(fmt.Scan("%+v\n", response.Header))
	}

}

func sendBasicAuth() {
	// jsonData := map[string]string{"firstname": "Nic", "lastname": "Raboy"}
	// jsonValue, _ := json.Marshal(jsonData)
	authString := "/basic-auth/user/pass"
	restEndpoint := restURL + authString
	log.Print("restEndpoint = " + restEndpoint)
	request, _ := http.NewRequest("GET", restEndpoint, bytes.NewBuffer(nil))
	request.Header.Set("Content-Type", "application/json")
	request.SetBasicAuth("user", "pass") // -> 200
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data))
		// header, _ := ioutil.ReadAll(response.Header)
		// fmt.Println(string(response.Header))

		// dump, err := httputil.DumpResponse(response, true)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// fmt.Printf("%q", dump)

		fmt.Println(response.StatusCode)
		// fmt.Printf("%+v\n", response.Header)
		log.Println(fmt.Sprintf("%+v\n", response.Header))

	}

}

func sendGet() {
	fmt.Println("Starting the application...")
	response, err := http.Get(restURL + "/ip")
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data))
	}
}

func sendPost1() {
	jsonData := map[string]string{"firstname": "Nic", "lastname": "Raboy"}
	jsonValue, _ := json.Marshal(jsonData)
	response, err := http.Post(restURL+"/ip", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data))
	}
	fmt.Println("Terminating the application...")
}
func sendPost2() {
	jsonData := map[string]string{"firstname": "Nic", "lastname": "Raboy"}
	jsonValue, _ := json.Marshal(jsonData)
	request, _ := http.NewRequest("POST", restURL+"/ip", bytes.NewBuffer(jsonValue))
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data))
	}
}
