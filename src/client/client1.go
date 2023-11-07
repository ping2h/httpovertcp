package client

// import (
// 	"bytes"
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"os"
// )

// func Client() {
// 	// Define the URL of the server where you want to send the POST request
// 	url := "http://localhost:8080/upload" // Replace with the actual URL

// 	// Define the path to the HTML or plain text file
// 	filePath := "src/client/example.html" // Replace with the actual file path

// 	// Open and read the file content
// 	file, err := os.Open(filePath)
// 	if err != nil {
// 		fmt.Println("Error opening the file:", err)
// 		return
// 	}
// 	defer file.Close()

// 	var requestBody bytes.Buffer
// 	_, err = io.Copy(&requestBody, file)
// 	if err != nil {
// 		fmt.Println("Error reading the file:", err)
// 		return
// 	}

// 	// Create a request with the "POST" method and request payload
// 	req, err := http.NewRequest("POST", url, &requestBody)
// 	if err != nil {
// 		fmt.Println("Error creating request:", err)
// 		return
// 	}

// 	// Set the request headers if needed (e.g., Content-Type)
// 	req.Header.Set("Content-Type", "text/html") // Adjust the content type if you are sending plain text

// 	// Create an HTTP client
// 	client := &http.Client{}

// 	// Send the request
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		fmt.Println("Error sending POST request:", err)
// 		return
// 	}
// 	defer resp.Body.Close()

// 	// Check the response status code
// 	if resp.Status != "200 OK" {
// 		fmt.Printf("Received non-200 status code: %s\n", resp.Status)
// 		return
// 	}

// 	// Read and print the response body
// 	responseBody, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		fmt.Println("Error reading response body:", err)
// 		return
// 	}

// 	fmt.Printf("Response body:\n%s\n", string(responseBody))
// }
