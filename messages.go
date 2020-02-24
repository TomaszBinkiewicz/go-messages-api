package main

import (
	"encoding/json"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
	"gopkg.in/gomail.v2"
	"log"
	"net/http"
	"strconv"
	"time"
)

// Message struct
type Message struct {
	Id          int    `json:"id"`
	Email       string `json:"email"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	MagicNumber int    `json:"magic_number"`
	Created     int    `json:"created"`
}

// SendTo struct
type SendTo struct {
	MagicNumber int `json:"magic_number"`
}

// Init id
var id int = 1

// Get current time as integer
func getTime() int {
	now := time.Now()
	timeStr := fmt.Sprintf("%02v%02v", now.Hour(), now.Minute())
	timeInt, err := strconv.Atoi(timeStr)
	if err != nil {
		log.Fatal(err)
	}
	return timeInt
}

// Create session
func cassandraConnection() *gocql.Session {
	//Init db
	cluster := gocql.NewCluster("172.18.0.2") // Insert cluster IP
	cluster.Consistency = gocql.Quorum
	cluster.ProtoVersion = 4
	cluster.ConnectTimeout = time.Second * 10
	cluster.Authenticator = gocql.PasswordAuthenticator{Username:"Username", Password:"Password"} // Insert auth credentials
	session, err := cluster.CreateSession()
	if err != nil {
		log.Println(err)
	}
	return session
}

// Execute query
func execQuery(query string) {
	session := cassandraConnection()
	defer session.Close()
	err := session.Query(query).Exec()
	if err != nil {
		log.Println(err)
		return
	}
}

// Get all messages as a slice
func getSliceMessages() []Message{
	var messages []Message
	var message Message

	session := cassandraConnection()
	defer session.Close()
	iter := session.Query("SELECT id, email, title, content, magic_number, created " +
		"FROM messages_space.messages_table;").Iter()
	for iter.Scan(&message.Id, &message.Email, &message.Title, &message.Content, &message.MagicNumber, &message.Created) {
		messages = append(messages, message)
	}
	if err := iter.Close(); err != nil {
		log.Fatal(err)
	}
	return messages
}

// Get messages by email as a slice
func getSliceMessagesEmail(email string) []Message{
	var messages []Message
	var message Message

	session := cassandraConnection()
	defer session.Close()
	query := fmt.Sprintf("SELECT id, email, title, content, magic_number, created FROM " +
		"messages_space.messages_table WHERE email='%v' ALLOW FILTERING;", email)
	iter := session.Query(query).Iter()
	for iter.Scan(&message.Id, &message.Email, &message.Title, &message.Content, &message.MagicNumber, &message.Created) {
		messages = append(messages, message)
	}
	if err := iter.Close(); err != nil {
		log.Fatal(err)
	}
	return messages
}

// Create a new message
func createMessage(w http.ResponseWriter, r *http.Request) {
	messages := getSliceMessages()
	var message Message

	_ = json.NewDecoder(r.Body).Decode(&message)
	message.Id = id
	message.Created = getTime()
	messages = append(messages, message)
	id += 1

	session := cassandraConnection()
	defer session.Close()

	values := fmt.Sprintf("%v, '%v', '%v', '%v', %v, %v", message.Id, message.Email, message.Title,
		message.Content, message.MagicNumber, message.Created)
	query := fmt.Sprintf("INSERT INTO messages_space.messages_table (id, email, title, content, magic_number," +
		"created) VALUES (%v);", values)
	execQuery(query)
	w.WriteHeader(201)
	w.Write([]byte("201 - Message created!"))
}

// Get all messages
func getAllMessages(w http.ResponseWriter, r *http.Request){
	messages := getSliceMessages()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

// Get messages by email
func getMessages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // Get params
	email := params["emailValue"]
	found := getSliceMessagesEmail(email)
	json.NewEncoder(w).Encode(found)
}

// Send message
func sendMessage(w http.ResponseWriter, r *http.Request) {
	messages := getSliceMessages()
	checkErr := false
	var sendTo SendTo
	var toDelete []int
	json.NewDecoder(r.Body).Decode(&sendTo)
	for _, item := range messages {
		if item.MagicNumber == sendTo.MagicNumber {
			// send email todo - insert valid credentials
			m := gomail.NewMessage()
			m.SetHeader("From", "author@example.com")
			m.SetHeader("To", item.Email)
			m.SetHeader("Subject", item.Title)
			m.SetBody("text/html", item.Content)
			d := gomail.NewDialer("smtp.example.com", 1111, "user", "password")

			// Send the email to Bob, Cora and Dan.
			if err := d.DialAndSend(m); err != nil {
				checkErr = true
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("500 - Something bad happened!"))
			}

			// delete from db
			toDelete = append(toDelete, item.Id)
		}
	}
	if len(toDelete) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if checkErr == false {
		for _, value := range toDelete {
			deleteMessage(value)
		}
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("202 - Message(s) sent!"))
	}
}

// Delete message
func deleteMessage(id int) {
	query := fmt.Sprintf("DELETE FROM messages_space.messages_table WHERE id=%v;", id)
	execQuery(query)
}

// Insert mock data
func createMockData() {
	//
	values := fmt.Sprintf("%v, %v, %v, %v, %v, %v", id, "'jan.kowalski@example.com'", "'Interview'",
		"'simple text'", 101, getTime())
	query := fmt.Sprintf("INSERT INTO messages_space.messages_table (id, email, title, content, magic_number," +
		"created) VALUES (%v);", values)
	id += 1
	execQuery(query)
	//
	values2 := fmt.Sprintf("%v, %v, %v, %v, %v, %v", id, "'jan.kowalski@example.com'", "'Interview 2'",
		"'simple text 2'", 22, getTime())
	query2 := fmt.Sprintf("INSERT INTO messages_space.messages_table (id, email, title, content, magic_number," +
		"created) VALUES (%v);", values2)
	id += 1
	execQuery(query2)
	//
	values3 := fmt.Sprintf("%v, %v, %v, %v, %v, %v", id, "'anna.zajkowska@example.com'", "'Interview 3'",
		"'simple text 3'", 101, getTime())
	query3 := fmt.Sprintf("INSERT INTO messages_space.messages_table (id, email, title, content, magic_number," +
		"created) VALUES (%v);", values3)
	id += 1
	execQuery(query3)
}

func main() {
	// Create keyspace
	query := "CREATE KEYSPACE IF NOT EXISTS messages_space WITH REPLICATION = {'class':'SimpleStrategy'," +
		"'replication_factor':1};"
	execQuery(query)

	// Create table
	query = "CREATE TABLE IF NOT EXISTS messages_space.messages_table (id int, email text, title text, content text," +
		"magic_number int, Created int,PRIMARY KEY (id, email));"
	execQuery(query)

	//createMockData() // init with mock data

	// Init router
	r := mux.NewRouter()

	// Checking for old messages
	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		for range ticker.C {
			messages := getSliceMessages()
			now := getTime()
			var toDelete []int
			for _, item := range messages {
				diff := now - item.Created
				if diff > 5 {
					toDelete = append(toDelete, item.Id)
				}
			}
			for _, value := range toDelete {
				deleteMessage(value)
			}
		}
	}()

	// Route handlers / endpoints
	r.HandleFunc("/api/messages", getAllMessages).Methods("GET")
	r.HandleFunc("/api/messages/{emailValue}", getMessages).Methods("GET")
	r.HandleFunc("/api/message", createMessage).Methods("POST")
	r.HandleFunc("/api/send", sendMessage).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", r))
}
