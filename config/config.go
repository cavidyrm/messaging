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
	Provider string
	APIKey   string
	From     string
	Endpoint string
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
			Username:     getEnv("DATABASE_USER", "postgres"),
			Password:     getEnv("DATABASE_PASSWORD", "postgres"),
			DatabaseName: getEnv("DATABASE_NAME", "postgres"),
			SSLMode:      getEnv("DATABASE_SSLMODE", "disable"),
		},
		EventDB: DatabaseConfig{
			Host:         getEnv("EVENT_HOST", "localhost"),
			Port:         getEnvInt("EVENT_PORT", 9090),
			Username:     getEnv("EVENT_USER", "postgres"),
			Password:     getEnv("EVENT_PASS", "postgres"),
			DatabaseName: getEnv("EVENT_DATABASE_NAME", "postgres"),
			SSLMode:      getEnv("EVENT_SSLMODE", "disable"),
		},
		Kafka: KafkaConfig{
			Brokers: []string{getEnv("KAFKA_HOST", "localhost:9092")},
			Topic:   getEnv("KAFKA_TOPIC", "messaging"),
			GroupID: getEnv("KAFKA_GROUP_ID", "messaging-service"),
		},
		SMS: SMSConfig{
			Provider: getEnv("SMS_PROVIDER", "kavehnegar"),
			APIKey:   getEnv("SMS_API_KEY", "kavehnegar-api-key"),
			Endpoint: getEnv("SMS_ENDPOINT", "https://api.kavehnegar.com/send"),
		},
		Email: EmailConfig{
			Provider: getEnv("EMAIL_PROVIDER", "arvan"),
			APIKey:   getEnv("EMAIL_API_KEY", "smtp-api-key"),
			From:     getEnv("EMAIL_FROM", "noreply@netpardaz.com"),
			Endpoint: getEnv("EMAIL_ENDPOINT", "https://smtp.com"),
		},
		Logger: LoggerConfig{
			Level: getEnv("LOG_LEVEL", "info"),
			File:  getEnv("LOG_FILE", "messaging.log"),
		},
	}, nil
}

func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
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
