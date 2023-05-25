package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

var otherServers = make(map[string]string)
var db *sql.DB
var err error

func main() {

	//open data base
	db, err = sql.Open("mysql", "root:@tcp(localhost:3306)/server3")
	if err != nil {
		fmt.Println("Failed to connect to MySQL:", err)
		return
	}
	defer db.Close()

	//createTable()

	// My Servers
	otherServers["S2"] = "http://localhost:8082"
	otherServers["S1"] = "http://localhost:8081"

	for {
		http.HandleFunc("/", handleRequest)

		fmt.Println("Waiting...")

		http.ListenAndServe(":8083", nil)
	}

}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	// Read the data from the client
	data, err := ioutil.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	header := r.Header.Get("request-from")
	fmt.Print("From: ", header)
	fmt.Printf(" >>  %s\n", data)

	//Add data to database
	insert(db, string(data))

	if header == "Client" {
		// Forward the data to other servers
		err = forwardDataToServers(data)
		if err != nil {
			http.Error(w, "Error forwarding data to servers", http.StatusInternalServerError)
			return
		} else {
			fmt.Println("Data forwarded successfully.")
		}
		// Send a response to the client
		fmt.Fprint(w, "Data received and forwarded successfully")

	} else {
		// from server
		fmt.Fprint(w, "Data received :)")
	}

}

func forwardDataToServers(data []byte) error {

	//make HTTP requests to other servers and send the data
	for _, serverAddress := range otherServers {

		req, err := http.NewRequest("POST", serverAddress, bytes.NewBuffer(data))
		req.Header.Set("request-from", "Server 3")

		if err != nil {
			return err
		}

		// Send the request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error sending request forwardDataToServers:", err)
			return err
		}
		defer resp.Body.Close()

	}

	return nil
}

func insert(db *sql.DB, message string) {

	insertQuery := "INSERT INTO chat (message) VALUES (?)"
	// Prepare the SQL statement

	stmt, err := db.Prepare(insertQuery)
	if err != nil {
		fmt.Println("Failed to prepare statement:", err)
		return
	}
	defer stmt.Close()

	// Execute the statement with the data
	_, err = stmt.Exec(message)
	if err != nil {
		fmt.Println("Failed to insert data:", err)
		return
	}
	fmt.Println("Data inserted in database successfully!")

}

func createTable() {
	createTableQuery := "CREATE TABLE chat (id INT AUTO_INCREMENT PRIMARY KEY,message VARCHAR(500));"
	_, err = db.Exec(createTableQuery)
	if err != nil {
		fmt.Println("Failed to create table:", err)
		return
	}
	fmt.Println("Table created successfully!")

}
