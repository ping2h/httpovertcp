package client

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
)

func Client() {
	// URL where you want to send the POST request
	url := "http://localhost:8080/upload"

	// JSON payload that you want to send in the request body
	jsonPayload := []byte(`{"key": "value"}`)

	// Create a new POST request with the JSON payload
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		panic(err)
	}

	// Set the request headers, such as Content-Type
	req.Header.Set("Content-Type", "application/json")

	// Create an HTTP client and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode == http.StatusOK {
		fmt.Println("Request was successful")
	} else {
		fmt.Println("Request failed with status code:", resp.StatusCode)
	}

	// You can read the response body here if needed
	// responseBody, err := ioutil.ReadAll(resp.Body)
	// fmt.Println("Response Body:", string(responseBody))
}

func ClientPost() {
	// Open the file you want to upload
	file, err := os.Open("/home/dellzp/tmp/dslab1/src/client/example.html")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Read the file data into a buffer
	fileData := new(bytes.Buffer)
	_, err = fileData.ReadFrom(file)
	if err != nil {
		panic(err)
	}

	// Make a POST request to the server
	url := "http://localhost:8080/upload"
	req, err := http.NewRequest("POST", url, fileData)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "text/plain")
	req.Header.Set("Content-Disposition", "inline; filename="+"example.html")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	// Check the response
	if resp.StatusCode == http.StatusOK {
		fmt.Println("200 ok")
	} else {
		// Handle an unsuccessful upload
	}

}
