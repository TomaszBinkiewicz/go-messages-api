package main

import (
	"github.com/gorilla/mux"
	"log"
	"messages/project_package"
	"net/http"
	"time"
)

func main() {
	// Create keyspace
	query := "CREATE KEYSPACE IF NOT EXISTS messages_space WITH REPLICATION = {'class':'SimpleStrategy'," +
		"'replication_factor':1};"
	project_package.ExecQuery(query)

	// Create table
	query = "CREATE TABLE IF NOT EXISTS messages_space.messages_table (id int, email text, title text, content text," +
		"magic_number int, Created int,PRIMARY KEY ((id), email));"
	project_package.ExecQuery(query)

	query = "CREATE INDEX IF NOT EXISTS ON messages_space.messages_table (email);"
	project_package.ExecQuery(query)

	project_package.CreateMockData() // init with mock data

	// Checking for old messages
	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		for range ticker.C {
			messages := project_package.GetSliceMessages()
			now := project_package.GetTime()
			var toDelete []int
			for _, item := range messages {
				diff := now - item.Created
				if diff > 5 {
					toDelete = append(toDelete, item.Id)
				}
			}
			for _, value := range toDelete {
				project_package.DeleteMessage(value)
			}
		}
	}()

	// Init router
	r := mux.NewRouter()

	// Route handlers / endpoints
	r.HandleFunc("/api/messages", project_package.GetAllMessages).Methods("GET")
	r.HandleFunc("/api/messages/{emailValue}", project_package.GetMessages).Methods("GET")
	r.HandleFunc("/api/message", project_package.CreateMessage).Methods("POST")
	r.HandleFunc("/api/send", project_package.SendMessage).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", r))
}
