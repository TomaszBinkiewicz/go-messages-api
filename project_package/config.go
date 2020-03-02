package project_package

import (
	"os"
)

type CassandraConfig struct {
	host     string
	username string
	password string
}

type EmailConfig struct {
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

var cassandraConfig = CassandraConfig{
	host:     getEnv("CASSANDRA_HOST", "db"),
	username: getEnv("CASSANDRA_USER", "Username"),
	password: getEnv("CASSANDRA_PASSWD", "Password"),
}

var emailConfig = EmailConfig{
	author:   getEnv("CASSANDRA_HOST", "author@example.com"),
	username: getEnv("CASSANDRA_USER", "user"),
	password: getEnv("CASSANDRA_PASSWD", "password"),
}
