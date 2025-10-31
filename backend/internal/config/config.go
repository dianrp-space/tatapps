package config

import (
	"os"
	"strconv"
)

type Config struct {
	// Application
	AppName string
	AppEnv  string
	AppPort string
	AppURL  string

	// Database
	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string
	DBSSLMode  string

	// JWT
	JWTSecret     string
	JWTExpiration int

	// WhatsApp API
	WAApiURL   string
	WAApiKey   string
	WASender   string

	// SMTP Email
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	SMTPFromEmail string
	SMTPFromName  string

	// Frontend
	FrontendURL string
}

func LoadConfig() *Config {
	jwtExp, _ := strconv.Atoi(getEnv("JWT_EXPIRATION", "24"))
	smtpPort, _ := strconv.Atoi(getEnv("SMTP_PORT", "587"))

	return &Config{
		AppName: getEnv("APP_NAME", "TatApps"),
		AppEnv:  getEnv("APP_ENV", "development"),
		AppPort: getEnv("APP_PORT", "8080"),
		AppURL:  getEnv("APP_URL", "http://localhost:8080"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBName:     getEnv("DB_NAME", "tatapps"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		JWTSecret:     getEnv("JWT_SECRET", "your-secret-key"),
		JWTExpiration: jwtExp,

		WAApiURL:   getEnv("WA_API_URL", "https://wa.drpnet.my.id/send-message"),
		WAApiKey:   getEnv("WA_API_KEY", ""),
		WASender:   getEnv("WA_SENDER", ""),

		SMTPHost:      getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:      smtpPort,
		SMTPUsername:  getEnv("SMTP_USERNAME", ""),
		SMTPPassword:  getEnv("SMTP_PASSWORD", ""),
		SMTPFromEmail: getEnv("SMTP_FROM_EMAIL", ""),
		SMTPFromName:  getEnv("SMTP_FROM_NAME", "TatApps"),

		FrontendURL: getEnv("FRONTEND_URL", "http://localhost:5173"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
