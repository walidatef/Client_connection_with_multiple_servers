package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

var servers = make(map[string]string)

func main() {
	// Data to be sent
	servers["1"] = "http://localhost:8081"
	servers["2"] = "http://localhost:8082"
	servers["3"] = "http://localhost:8083"

	scanner := bufio.NewScanner(os.Stdin)

	for {

		fmt.Print("Enter your data:\n ")
		scanner.Scan()
		data := scanner.Bytes()

		fmt.Print("Choose your server you need send data:\n 1. server 1 \n 2. server 2 \n 3. server 3 \n")
		fmt.Print(">>")

		scanner.Scan()
		server := scanner.Text()

		for server != "1" && server != "2" && server != "3" {
			fmt.Println("Server not found in list :(")

			fmt.Print(">>")
			fmt.Scanln(&server)
		}

		// Create an HTTP request
		req, err := http.NewRequest("POST", servers[server], bytes.NewBuffer(data))
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}

		// custom header
		req.Header.Set("Request-from", "Client")
		// Send the request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error sending request:", err)
			return
		}
		defer resp.Body.Close()

		// Read the response
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response:", err)
			return
		}

		// Print the response
		fmt.Println("Server response:", string(body))
	}
}
