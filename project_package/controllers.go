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

	_ = json.NewDecoder(r.Body).Decode(&message)
	message.Id = id
	message.Created = GetTime()
	messages = append(messages, message)
	id += 1

	session := CassandraConnection()
	defer session.Close()

	values := fmt.Sprintf("%v, '%v', '%v', '%v', %v, %v", message.Id, message.Email, message.Title,
		message.Content, message.MagicNumber, message.Created)
	query := fmt.Sprintf("INSERT INTO messages_space.messages_table (id, email, title, content, magic_number," +
		"created) VALUES (%v);", values)
	ExecQuery(query)
	w.WriteHeader(201)
	w.Write([]byte("201 - Message created!"))
}

// Get all messages
func GetAllMessages(w http.ResponseWriter, r *http.Request){
	messages := GetSliceMessages()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

// Get messages by email
func GetMessages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // Get params
	email := params["emailValue"]
	found := GetSliceMessagesEmail(email)
	json.NewEncoder(w).Encode(found)
}

// Send message
func SendMessage(w http.ResponseWriter, r *http.Request) {
	messages := GetSliceMessages()
	checkErr := false
	var sendTo SendTo
	var toDelete []int
	json.NewDecoder(r.Body).Decode(&sendTo)
	for _, item := range messages {
		if item.MagicNumber == sendTo.MagicNumber {
			// send email
			//todo - insert valid credentials
			m := gomail.NewMessage()
			m.SetHeader("From", "author@example.com")
			m.SetHeader("To", item.Email)
			m.SetHeader("Subject", item.Title)
			m.SetBody("text/html", item.Content)
			d := gomail.NewDialer("smtp.example.com", 1111, "user", "password")

			// Send the email to Bob, Cora and Dan.
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

// Insert mock data
func CreateMockData() {
	//
	values := fmt.Sprintf("%v, %v, %v, %v, %v, %v", id, "'jan.kowalski@example.com'", "'Interview'",
		"'simple text'", 101, GetTime())
	query := fmt.Sprintf("INSERT INTO messages_space.messages_table (id, email, title, content, magic_number," +
		"created) VALUES (%v);", values)
	id += 1
	ExecQuery(query)
	//
	values2 := fmt.Sprintf("%v, %v, %v, %v, %v, %v", id, "'jan.kowalski@example.com'", "'Interview 2'",
		"'simple text 2'", 22, GetTime())
	query2 := fmt.Sprintf("INSERT INTO messages_space.messages_table (id, email, title, content, magic_number," +
		"created) VALUES (%v);", values2)
	id += 1
	ExecQuery(query2)
	//
	values3 := fmt.Sprintf("%v, %v, %v, %v, %v, %v", id, "'anna.zajkowska@example.com'", "'Interview 3'",
		"'simple text 3'", 101, GetTime())
	query3 := fmt.Sprintf("INSERT INTO messages_space.messages_table (id, email, title, content, magic_number," +
		"created) VALUES (%v);", values3)
	id += 1
	ExecQuery(query3)
}

