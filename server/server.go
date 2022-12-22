package server

import (
	"encoding/json"
	"go_demo/db"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	// import rabbit package
	"go_demo/rabbit"
)

const (
	// RabbitMQURL is the URL of the RabbitMQ server
	RabbitMQURL = "amqp://guest:guest@localhost:5672/"
)

func handleMessage(w http.ResponseWriter, r *http.Request) {
	// Decode the request body into an SMS struct
	var message db.Message
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		// Handle error
		log.Fatal("Error decoding request body: ", err)
	}

	// Borrow a connection from the pool
	conn := db.GetConnection()
	// Return the connection to the pool
	defer db.PutConnection(conn)

	// Use the connection

	// Create a new record in the database
	if err := conn.Create(&message).Error; err != nil {
		// Handle error
		log.Fatal("Error creating record: ", err)
		//panic(err)
	}

	// create a new RabbitMQ connection
	rabbitConn, err := rabbit.NewConnection(RabbitMQURL)
	if err != nil {
		log.Println("Error creating rabbit connection: ", err)
		//panic(err)
	}
	defer rabbitConn.Close()

	// publish the message to the queue
	if err := rabbitConn.Publish("messages", message); err != nil {
		log.Println("Error publishing message: ", err)
		//panic(err)
	}

	log.Println("Message published to queue")

	// Return a response to the client
	w.WriteHeader(http.StatusOK)
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	// write h1 html response to the client with the message "Server is up and running"
	w.Write([]byte("<h1>Server is up and running</h1>"))
}

// create a function that will display the messages in the database in table format
func handleMessages(w http.ResponseWriter, r *http.Request) {
	// Borrow a connection from the pool
	conn := db.GetConnection()

	// Use the connection

	// Create a new record in the database
	var messages []db.Message
	if err := conn.Find(&messages).Error; err != nil {
		// Handle error
		log.Fatal("Error creating record: ", err)
	}

	// Return the connection to the pool
	defer db.PutConnection(conn)

	// Return a response to the client
	w.WriteHeader(http.StatusOK)
	// create a table with the messages
	w.Write([]byte("<table>"))
	for _, message := range messages {
		w.Write([]byte("<tr>"))
		w.Write([]byte("<td>" + message.ShortCode + "</td>"))
		w.Write([]byte("<td>" + message.CreatedAt.String() + "</td>"))
		w.Write([]byte("</tr>"))
	}
	w.Write([]byte("</table>"))

}

func StartServer() {
	r := mux.NewRouter()
	r.HandleFunc("/api/intake", handleMessage).Methods("POST")
	r.HandleFunc("/api/status", handleStatus).Methods("GET")
	r.HandleFunc("/api/messages", handleMessages).Methods("GET")
	log.Println("Server is up and running")
	http.ListenAndServe(":8080", r)
}
