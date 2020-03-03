package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"messages/project_package"
	"net/http"
)

func main() {
	// Create keyspace
	keyspace := project_package.CassandraConfig.Keyspace
	query := fmt.Sprintf("CREATE KEYSPACE IF NOT EXISTS %v WITH REPLICATION = {'class':'SimpleStrategy', " +
		"'replication_factor':1};", keyspace)
	project_package.ExecQuery(query)
	project_package.KeyspaceInitialized = true

	// Create table
	query = fmt.Sprintf("CREATE TABLE IF NOT EXISTS %v.messages_table (id int, email text, title text, " +
		"content text, magic_number int, Created int,PRIMARY KEY ((id), email));", keyspace)
	project_package.ExecQuery(query)

	query = fmt.Sprintf("CREATE INDEX IF NOT EXISTS ON %v.messages_table (email);", keyspace)
	project_package.ExecQuery(query)

	// Checking for old messages
	project_package.DeleteOldMessages()

	// Init router
	r := mux.NewRouter()

	// Route handlers / endpoints
	r.HandleFunc("/api/messages", project_package.GetAllMessages).Methods("GET")
	r.HandleFunc("/api/messages/{emailValue}", project_package.GetMessages).Methods("GET")
	r.HandleFunc("/api/message", project_package.CreateMessage).Methods("POST")
	r.HandleFunc("/api/send", project_package.SendMessage).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", r))
}
