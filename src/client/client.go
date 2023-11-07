package client

import (
	"bytes"
	"fmt"
	"net/http"
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
