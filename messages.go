package main

import (
	"encoding/json"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
	"gopkg.in/gomail.v2"
	"log"
	"net/http"
	"time"
)

// Message struct
type Message struct {
	Id          int       `json:"id"`
	Email       string    `json:"email"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	MagicNumber int       `json:"magic_number"`
	Created     time.Time `json:"created"`
}

// SendTo struct
type SendTo struct {
	MagicNumber int `json:"magic_number"`
}

// Init messages var as slice Message struct

var messages []Message

// Init id

var id int = 1

// Create a new message

func createMessage(w http.ResponseWriter, r *http.Request) {
	var message Message
	_ = json.NewDecoder(r.Body).Decode(&message)
	message.Id = id
	message.Created = time.Now()
	messages = append(messages, message)
	id += 1
	json.NewEncoder(w).Encode(message)
}

// Get all messages

func getAllMessages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

// Get messages by email

func getMessages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // Get params
	var found []Message
	// loop through messages
	for _, item := range messages {
		if item.Email == params["emailValue"] {
			found = append(found, item)
		}
	}
	json.NewEncoder(w).Encode(found)
}

// Send message

func sendMessage(w http.ResponseWriter, r *http.Request) {
	var sendTo SendTo
	var toDelete []int
	json.NewDecoder(r.Body).Decode(&sendTo)
	for _, item := range messages {
		if item.MagicNumber == sendTo.MagicNumber {
			// send email
			m := gomail.NewMessage()
			m.SetHeader("From", "author@example.com")
			m.SetHeader("To", item.Email)
			m.SetHeader("Subject", item.Title)
			m.SetBody("text/html", item.Content)
			d := gomail.NewDialer("smtp.example.com", 1111, "user", "password")

			// Send the email to Bob, Cora and Dan.
			if err := d.DialAndSend(m); err != nil {
				//panic(err) // todo - uncomment for error notification
			}

			// delete from db
			toDelete = append(toDelete, item.Id)
		}
	}
	for _, value := range toDelete {
		deleteMessage(value)
	}
}

// Delete message

func deleteMessage(id int) {
	// loop through messages
	for index, item := range messages {
		if item.Id == id {
			messages = append(messages[:index], messages[index+1:]...)
			break
		}
	}
}

func main() {
	//Init db
	cluster := gocql.NewCluster("172.19.0.2") // Insert cluster IP
	cluster.Consistency = gocql.Quorum
	cluster.ProtoVersion = 4
	cluster.ConnectTimeout = time.Second * 10
	//cluster.Authenticator = gocql.PasswordAuthenticator{Username:"Username", Password:"Password"} // Insert auth credentials
	session, err := cluster.CreateSession()
	if err != nil {
		log.Println(err)
		return
	}
	defer session.Close()

	// Create keyspace
	err = session.Query("CREATE KEYSPACE IF NOT EXISTS messages_space WITH REPLICATION =" +
		"{'class':'SimpleStrategy','replication_factor':1};").Exec()
	if err != nil {
		log.Println(err)
		return
	}

	// Create table
	err = session.Query("CREATE TABLE IF NOT EXISTS messages_space.messages_table" +
		"(id int, email text, title text, content text, magic_number int, created time,PRIMARY KEY (id));").Exec()
	if err != nil {
		log.Println(err)
		return
	}

	// Init router
	r := mux.NewRouter()

	// Mock data - todo - implement DB
	now := time.Now()
	values := fmt.Sprintf("%v, %v, %v, %v, %v, %02v:%02v:%02v", id, "jan.kowalski@example.com", "Interview",
		"simple text", 101, now.Hour(), now.Minute(), now.Second())
	query := fmt.Sprintf("INSERT INTO messages_space.messages_table (id, email, title, content, magic_number, created) VALUES (%v);", values)
	id += 1
	err = session.Query(query).Exec()
	//
	values2 := fmt.Sprintf("%v, %v, %v, %v, %v, %02v:%02v:%02v",	id, "jan.kowalski@example.com", "Interview 2",
		"simple text 2", 101, now.Hour(), now.Minute(), now.Second())
	query2 := fmt.Sprintf("INSERT INTO messages_space.messages_table (id, email, title, content, magic_number, created) VALUES (%v);", values2)
	id += 1
	err = session.Query(query2).Exec()
	//
	values3 := fmt.Sprintf("%v, %v, %v, %v, %v, %02v:%02v:%02v",	id, "anna.zajkowska@example.com", "Interview 3",
		"simple text 3", 101, now.Hour(), now.Minute(), now.Second())
	query3 := fmt.Sprintf("INSERT INTO messages_space.messages_table (id, email, title, content, magic_number, created) VALUES (%v);", values3)
	id += 1
	err = session.Query(query3).Exec()
	if err != nil {
		log.Println(err)
		return
	}

	//messages = append(messages, Message{Id: id, Email: "jan.kowalski@example.com",
	//	Title: "Interview", Content: "simple text", MagicNumber: 101, Created: time.Now()})
	//id += 1
	//messages = append(messages, Message{Id: id, Email: "jan.kowalski@example.com",
	//	Title: "Interview 2", Content: "simple text 2", MagicNumber: 22, Created: time.Now()})
	//id += 1
	//messages = append(messages, Message{Id: id, Email: "anna.zajkowska@example.com",
	//	Title: "Interview 3", Content: "simple text 3", MagicNumber: 101, Created: time.Now()})
	//id += 1

	// Checking for old messages
	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		for range ticker.C {
			now := time.Now()
			var toDelete []int
			for _, item := range messages {
				diff := now.Sub(item.Created)
				if diff.Minutes() > 5 {
					toDelete = append(toDelete, item.Id)
				}
			}
			println("check performed")
			for _, value := range toDelete {
				deleteMessage(value)
				println("message deleted")
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
