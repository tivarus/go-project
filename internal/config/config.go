package config

import (
	"os"
	"strconv"
)

type Config struct {
	DBHost       string
	DBPort       int
	DBUser       string
	DBPassword   string
	DBName       string
	ServerPort   string
	JWTSecret    string
	HMACSecret   string
	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
	SMTPFrom     string
}

func Load() (*Config, error) {
	port, _ := strconv.Atoi(getEnv("DB_PORT", "5432"))
	smtpPort, _ := strconv.Atoi(getEnv("SMTP_PORT", "587"))

	return &Config{
		DBHost:       getEnv("DB_HOST", "localhost"),
		DBPort:       port,
		DBUser:       getEnv("DB_USER", "postgres"),
		DBPassword:   getEnv("DB_PASSWORD", "1234"),
		DBName:       getEnv("DB_NAME", "banking"),
		ServerPort:   getEnv("SERVER_PORT", ":8080"),
		JWTSecret:    getEnv("JWT_SECRET", "secret"),
		HMACSecret:   getEnv("HMAC_SECRET", "secret"),
		SMTPHost:     getEnv("SMTP_HOST", "smtp.example.com"),
		SMTPPort:     smtpPort,
		SMTPUser:     getEnv("SMTP_USER", "user@example.com"),
		SMTPPassword: getEnv("SMTP_PASSWORD", "password"),
		SMTPFrom:     getEnv("SMTP_FROM", "noreply@example.com"),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
