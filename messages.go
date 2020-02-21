package main

import (
	"encoding/json"
	_ "github.com/gocql/gocql"
	"github.com/gorilla/mux"
	"gopkg.in/gomail.v2"
	"log"
	"math/rand"
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

// Create a new message

func createMessage(w http.ResponseWriter, r *http.Request) {
	var message Message
	_ = json.NewDecoder(r.Body).Decode(&message)
	message.Id = rand.Intn(10000000) // Mock ID - not safe - todo
	message.Created = time.Now()
	messages = append(messages, message)
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
	// Init router
	r := mux.NewRouter()

	// Mock data - todo - implement DB
	messages = append(messages, Message{Id: 1, Email: "jan.kowalski@example.com",
		Title: "Interview", Content: "simple text", MagicNumber: 101, Created: time.Now()})
	messages = append(messages, Message{Id: 2, Email: "jan.kowalski@example.com",
		Title: "Interview 2", Content: "simple text 2", MagicNumber: 22, Created: time.Now()})
	messages = append(messages, Message{Id: 3, Email: "anna.zajkowska@example.com",
		Title: "Interview 3", Content: "simple text 3", MagicNumber: 101, Created: time.Now()})

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
