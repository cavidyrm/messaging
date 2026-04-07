package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	EventDB  DatabaseConfig
	Kafka    KafkaConfig
	SMS      SMSConfig
	Email    EmailConfig
	Logger   LoggerConfig
}
type ServerConfig struct {
	Port int
}

type DatabaseConfig struct {
	Host         string
	Port         int
	Username     string
	Password     string
	DatabaseName string
	SSLMode      string
}

type EventDBConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
	SSLMode  string
}

type KafkaConfig struct {
	Brokers []string
	Topic   string
	GroupID string
}
type EmailConfig struct {
	Host     string
	APIKey   string
	From     string
	Port     string
	Username string
	Password string
}

type SMSConfig struct {
	Provider string
	APIKey   string
	Endpoint string
}
type LoggerConfig struct {
	Level string
	File  string
}

func Load() (*Config, error) {
	return &Config{
		Server: ServerConfig{
			Port: getEnvInt("SERVER_PORT", 8080),
		},
		Database: DatabaseConfig{
			Host:         getEnv("DATABASE_HOST", "localhost"),
			Port:         getEnvInt("DATABASE_PORT", 5432),
			Username:     getEnv("DATABASE_USER", "user"),
			Password:     getEnv("DATABASE_PASSWORD", "password"),
			DatabaseName: getEnv("DATABASE_NAME", "my_database"),
			SSLMode:      getEnv("DATABASE_SSLMODE", "disable"),
		},
		EventDB: DatabaseConfig{
			Host:         getEnv("EVENT_HOST", "localhost"),
			Port:         getEnvInt("EVENT_PORT", 5433),
			Username:     getEnv("EVENT_USER", "user"),
			Password:     getEnv("EVENT_PASS", "password"),
			DatabaseName: getEnv("EVENT_DATABASE_NAME", "my_database"),
			SSLMode:      getEnv("EVENT_SSLMODE", "disable"),
		},
		Kafka: KafkaConfig{
			Brokers: []string{getEnv("KAFKA_HOST", "localhost:9092")},
			Topic:   getEnv("KAFKA_TOPIC", "messaging"),
			GroupID: getEnv("KAFKA_GROUP_ID", "messaging-group"),
		},
		SMS: SMSConfig{
			Provider: getEnv("SMS_PROVIDER", "kavehnegar"),
			APIKey:   getEnv("SMS_API_KEY", "kavehnegar-api-key"),
			Endpoint: getEnv("SMS_ENDPOINT", "https://api.kavehnegar.com/send"),
		},
		Email: EmailConfig{
			Host:     getEnv("EMAIL_HOST", "arvan"),
			APIKey:   getEnv("EMAIL_API_KEY", "smtp-api-key"),
			From:     getEnv("EMAIL_FROM", "noreply@netpardaz.com"),
			Port:     getEnv("EMAIL_PORT", "https://smtp.com"),
			Username: getEnv("EMAIL_USER", "smtp-user"),
			Password: getEnv("EMAIL_PASS", "smtp-pass"),
		},
		Logger: LoggerConfig{
			Level: getEnv("LOG_LEVEL", "info"),
			File:  getEnv("LOG_FILE", "kafka.log"),
		},
	}, nil
}

func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.Username, d.Password, d.DatabaseName, d.SSLMode)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return defaultValue
}
