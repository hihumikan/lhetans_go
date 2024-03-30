package main

import (
	"bytes"
	"fmt"
	"net/http"
)

func main() {
	url := "https://example.com/api/endpoint"

	jsonStr := []byte(`{"key1":"value1","key2":"value2"}`)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("POST request failed with status code:", resp.StatusCode)
		return
	}

	var responseBody bytes.Buffer
	_, err = responseBody.ReadFrom(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	fmt.Println("POST request was successful!")
	fmt.Println("Response:")
	fmt.Println(responseBody.String())
}
