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
	tt := []struct {
		email    string
		validate bool
	}{
		{
			"some.email@example.com",
			false,
		},
		{
			"some12.!#$%&'*+/=?^_`{|}~-EMAIL@example.com",
			false,
		},
		{
			"some.email.example.com",
			true,
		},
		{
			"some.email@",
			true,
		},
		{
			"@example.com",
			true,
		},
	}
	for _, tc := range tt {
		validation := ValidateEmail(tc.email)
		if validation == tc.validate {
			t.Errorf("Function ValidateEmail should return %v, but returned %v.", tc.validate, validation)
		}
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

func PrepareTestKesyspace() *gocql.Session {
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

func AddTestingData() {
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

func DropTestKeyspace() {
	query := "DROP KEYSPACE IF EXISTS test_space;"
	ExecQuery(query)
	KeyspaceInitialized = false
	id = 1
}

func TestGetSliceMessages(t *testing.T) {
	session := PrepareTestKesyspace()
	AddTestingData()
	defer session.Close()
	// proper email address
	numberOfMails := len(GetSliceMessages())
	if numberOfMails != 3 {
		t.Errorf("Function GetSliceMessages should return 3 messages, but returned %v.", numberOfMails)
	}
	DropTestKeyspace()
}

func TestGetSliceMessagesEmail(t *testing.T) {
	session := PrepareTestKesyspace()
	AddTestingData()
	defer session.Close()
	tt := []struct {
		email    string
		quantity int
	}{
		{
			"jan.kowalski@example.com",
			2,
		},
		{
			"anna.zajkowska@example.com",
			1,
		},
		{
			"no.such.email@example.com",
			0,
		},
	}
	for _, tc := range tt {
		// jan.kowalski@example.com
		numberOfMails := len(GetSliceMessagesEmail(tc.email))
		if numberOfMails != tc.quantity {
			t.Errorf("Function GetSliceMessages should return 3 messages, but returned %v.", numberOfMails)
		}
	}
	DropTestKeyspace()
}

func TestDeleteOldMessages(t *testing.T) {
	session := PrepareTestKesyspace()
	AddTestingData()
	defer session.Close()
	DeleteOldMessages()
	time.Sleep(65 * time.Second)
	numberOfMails := len(GetSliceMessages())
	if numberOfMails != 1 {
		t.Errorf("Function GetSliceMessages should return 1 messages, but returned %v.", numberOfMails)
	}
	DropTestKeyspace()
}
