package project_package

import (
	"fmt"
	"github.com/gocql/gocql"
	"reflect"
	"testing"
	"time"
)

func TestGetTime(t *testing.T) {
	var intTime interface{} = GetTime()
	_, ok := intTime.(int)
	if !ok {
		t.Errorf("GetTime return value type is %v, want int", reflect.TypeOf(intTime))
	}
}

func TestValidateEmail(t *testing.T) {
	// proper email address
	validation := ValidateEmail("some.email@example.com")
	if validation == false {
		t.Errorf("Function ValidateEmail should return true, but returned false.")
	}
	// proper email address
	validation = ValidateEmail("some12.!#$%&'*+/=?^_`{|}~-EMAIL@example.com")
	if validation == false {
		t.Errorf("Function ValidateEmail should return true, but returned false.")
	}
	// incorrect email address
	validation = ValidateEmail("some.email.example.com")
	if validation == true {
		t.Errorf("Function ValidateEmail should return false, but returned true.")
	}
	// incorrect email address
	validation = ValidateEmail("some.email@")
	if validation == true {
		t.Errorf("Function ValidateEmail should return false, but returned true.")
	}
	// proper email address
	validation = ValidateEmail("@example.com")
	if validation == true {
		t.Errorf("Function ValidateEmail should return false, but returned true.")
	}
}

func TestCassandraConnection(t *testing.T) {
	var session interface{} = CassandraConnection()
	temp, ok := session.(*gocql.Session)
	if !ok {
		t.Errorf("GetTime return value type is %v, want int", reflect.TypeOf(session))
	}
	temp.Close()
}

func prepareTestKesysoace() *gocql.Session {
	// prepare mock data
	session := CassandraConnection()

	query := "DROP KEYSPACE IF EXISTS test_space;"
	ExecQuery(query)

	// Create keyspace
	query = "CREATE KEYSPACE test_space WITH REPLICATION = {'class':'SimpleStrategy', 'replication_factor':1};"
	ExecQuery(query)
	KeyspaceInitialized = true

	// Create table
	query = "CREATE TABLE test_space.messages_table (id int, email text, title text, content text, magic_number int, " +
		"Created int,PRIMARY KEY ((id), email));"
	ExecQuery(query)

	query = "CREATE INDEX IF NOT EXISTS ON test_space.messages_table (email);"
	ExecQuery(query)

	return session
}

func addTestingData() {
	// Add testing data
	values := fmt.Sprintf("%v, %v, %v, %v, %v, %v", id, "'jan.kowalski@example.com'", "'Interview'",
		"'simple text'", 101, GetTime()-7)
	query := fmt.Sprintf("INSERT INTO test_space.messages_table (id, email, title, content, magic_number, "+
		"created) VALUES (%v);", values)
	id += 1
	ExecQuery(query)

	values = fmt.Sprintf("%v, %v, %v, %v, %v, %v", id, "'jan.kowalski@example.com'", "'Interview 2'",
		"'simple text 2'", 22, GetTime())
	query = fmt.Sprintf("INSERT INTO test_space.messages_table (id, email, title, content, magic_number,"+
		"created) VALUES (%v);", values)
	id += 1
	ExecQuery(query)

	values = fmt.Sprintf("%v, %v, %v, %v, %v, %v", id, "'anna.zajkowska@example.com'", "'Interview 3'",
		"'simple text 3'", 101, GetTime()-6)
	query = fmt.Sprintf("INSERT INTO test_space.messages_table (id, email, title, content, magic_number,"+
		"created) VALUES (%v);", values)
	id += 1
	ExecQuery(query)
}

func dropTestKeyspace() {
	query := "DROP KEYSPACE IF EXISTS test_space;"
	ExecQuery(query)
	KeyspaceInitialized = false
}

func TestGetSliceMessages(t *testing.T) {
	session := prepareTestKesysoace()
	addTestingData()
	defer session.Close()
	// proper email address
	numberOfMails := len(GetSliceMessages())
	if numberOfMails != 3 {
		t.Errorf("Function GetSliceMessages should return 3 messages, but returned %v.", numberOfMails)
	}
	dropTestKeyspace()
}

func TestGetSliceMessagesEmail(t *testing.T) {
	session := prepareTestKesysoace()
	addTestingData()
	defer session.Close()
	// jan.kowalski@example.com
	numberOfMails := len(GetSliceMessagesEmail("jan.kowalski@example.com"))
	if numberOfMails != 2 {
		t.Errorf("Function GetSliceMessages should return 3 messages, but returned %v.", numberOfMails)
	}
	// anna.zajkowska@example.com
	numberOfMails = len(GetSliceMessagesEmail("anna.zajkowska@example.com"))
	if numberOfMails != 1 {
		t.Errorf("Function GetSliceMessages should return 2 messages, but returned %v.", numberOfMails)
	}
	// no.such.email@example.com
	numberOfMails = len(GetSliceMessagesEmail("no.such.email@example.com"))
	if numberOfMails != 0 {
		t.Errorf("Function GetSliceMessages should return 0 messages, but returned %v.", numberOfMails)
	}
	dropTestKeyspace()
}

func TestDeleteOldMessages(t *testing.T) {
	session := prepareTestKesysoace()
	addTestingData()
	defer session.Close()
	DeleteOldMessages()
	time.Sleep(65 * time.Second)
	numberOfMails := len(GetSliceMessages())
	if numberOfMails != 1 {
		t.Errorf("Function GetSliceMessages should return 1 messages, but returned %v.", numberOfMails)
	}
	dropTestKeyspace()
}
