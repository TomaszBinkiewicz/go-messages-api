package project_package

import (
	"fmt"
	"github.com/gocql/gocql"
	"log"
	"strconv"
	"time"
)

// Init id
var id int = 1


// Get current time as integer
func GetTime() int {
	now := time.Now()
	timeStr := fmt.Sprintf("%02v%02v", now.Hour(), now.Minute())
	timeInt, err := strconv.Atoi(timeStr)
	if err != nil {
		log.Fatal(err)
	}
	return timeInt
}

// Create session
func CassandraConnection() *gocql.Session {
	//Init db
	cluster := gocql.NewCluster("db") // Insert cluster IP
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
func ExecQuery(query string) {
	session := CassandraConnection()
	defer session.Close()
	err := session.Query(query).Exec()
	if err != nil {
		log.Println(err)
		return
	}
}

// Get all messages as a slice
func GetSliceMessages() []Message{
	var messages []Message
	var message Message

	session := CassandraConnection()
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
func GetSliceMessagesEmail(email string) []Message{
	var messages []Message
	var message Message

	session := CassandraConnection()
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

