package project_package

import (
	"os"
)

type CassandraConfigStruct struct {
	host     string
	Keyspace string
	username string
	password string
}

type EmailConfigStruct struct {
	author   string
	username string
	password string
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

var CassandraConfig = CassandraConfigStruct{
	host:     getEnv("CASSANDRA_HOST", "db"), // database IP for testing
	Keyspace: getEnv("CASSANDRA_KEYSPACE", "messages_space"), // test_space for testing
	username: getEnv("CASSANDRA_USER", "Username"),
	password: getEnv("CASSANDRA_PASSWD", "Password"),
}

var EmailConfig = EmailConfigStruct{
	author:   getEnv("CASSANDRA_HOST", "author@example.com"),
	username: getEnv("CASSANDRA_USER", "user"),
	password: getEnv("CASSANDRA_PASSWD", "password"),
}
