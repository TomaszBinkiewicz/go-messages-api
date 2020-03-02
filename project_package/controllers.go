package project_package

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gopkg.in/gomail.v2"
	"net/http"
)

// Create a new message
func CreateMessage(w http.ResponseWriter, r *http.Request) {
	messages := GetSliceMessages()
	var message Message

	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil || ValidateEmail(message.Email) == false {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - bad data"))
		return
	}
	message.Id = id
	message.Created = GetTime()
	messages = append(messages, message)
	id += 1

	session := CassandraConnection()
	defer session.Close()

	values := fmt.Sprintf("%v, '%v', '%v', '%v', %v, %v", message.Id, message.Email, message.Title,
		message.Content, message.MagicNumber, message.Created)
	query := fmt.Sprintf("INSERT INTO messages_space.messages_table (id, email, title, content, magic_number,"+
		"created) VALUES (%v);", values)
	ExecQuery(query)
	w.WriteHeader(201)
	w.Write([]byte("201 - Message created!"))
}

// Get all messages
func GetAllMessages(w http.ResponseWriter, r *http.Request) {
	messages := GetSliceMessages()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

// Get messages by email
func GetMessages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // Get params
	email := params["emailValue"]
	if ValidateEmail(email) == false {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - invalid email address"))
		return
	}
	found := GetSliceMessagesEmail(email)
	json.NewEncoder(w).Encode(found)
}

// Send message
func SendMessage(w http.ResponseWriter, r *http.Request) {
	messages := GetSliceMessages()
	checkErr := false
	var sendTo SendTo
	var toDelete []int
	err := json.NewDecoder(r.Body).Decode(&sendTo)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - bad data"))
		return
	}
	for _, item := range messages {
		if item.MagicNumber == sendTo.MagicNumber {
			// send email
			m := gomail.NewMessage()
			m.SetHeader("From", emailConfig.author)
			m.SetHeader("To", item.Email)
			m.SetHeader("Subject", item.Title)
			m.SetBody("text/html", item.Content)
			d := gomail.NewDialer("smtp.example.com", 1111, emailConfig.username, emailConfig.password)

			if err := d.DialAndSend(m); err != nil {
				// todo - uncomment for errors handling
				//checkErr = true
				//w.WriteHeader(http.StatusInternalServerError)
				//w.Write([]byte("500 - Something bad happened!"))
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
			DeleteMessage(value)
		}
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("202 - Message(s) sent!"))
	}
}

// Delete message
func DeleteMessage(id int) {
	query := fmt.Sprintf("DELETE FROM messages_space.messages_table WHERE id=%v;", id)
	ExecQuery(query)
}
