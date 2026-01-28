package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port               string
	BaseURL            string
	DatabaseURL        string
	JWTSecret          string
	JWTExpirationHour  int
	CORSAllowedOrigins string
	FrontendURL        string

	Mail    MailConfig
	Storage StorageConfig
}

type MailConfig struct {
	Host string
	Port int
	User string
	Pass string
	From string
}

type StorageConfig struct {
	// Storage type: "local" or "s3"
	Type string

	// Local storage settings
	LocalUploadDir string

	// S3 storage settings
	S3Region          string
	S3Bucket          string
	S3AccessKeyID     string
	S3SecretAccessKey string
	S3Endpoint        string // Optional: for custom S3-compatible endpoints (MinIO, DigitalOcean Spaces, etc.)
}

func LoadConfig() (*Config, error) {
	// Try to load .env from multiple possible locations
	_ = godotenv.Load() // Current directory
	_ = godotenv.Load(".env")
	_ = godotenv.Load("../.env")
	_ = godotenv.Load("../../.env")

	cfg := &Config{
		Port:               getEnv("PORT", "8080"),
		BaseURL:            getEnv("BASE_URL", "http://localhost:8080"),
		DatabaseURL:        mustEnv("DATABASE_URL"),
		JWTSecret:          getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		JWTExpirationHour:  getEnvAsInt("JWT_EXPIRATION_HOUR", 24),
		CORSAllowedOrigins: getEnv("CORS_ALLOWED_ORIGINS", "*"),
		FrontendURL:        getEnv("FRONTEND_URL", "http://localhost:3000"),
		Mail: MailConfig{
			Host: getEnv("SMTP_HOST", "smtp.gmail.com"),
			Port: getEnvAsInt("SMTP_PORT", 587),
			User: getEnv("SMTP_USER", ""),
			Pass: getEnv("SMTP_PASS", ""),
			From: getEnv("SMTP_FROM", "noreply@example.com"),
		},
		Storage: StorageConfig{
			Type:              getEnv("STORAGE_TYPE", "local"),
			LocalUploadDir:    getEnv("LOCAL_UPLOAD_DIR", "./uploads"),
			S3Region:          getEnv("AWS_S3_REGION", "us-east-1"),
			S3Bucket:          getEnv("AWS_S3_BUCKET", ""),
			S3AccessKeyID:     getEnv("AWS_ACCESS_KEY_ID", ""),
			S3SecretAccessKey: getEnv("AWS_SECRET_ACCESS_KEY", ""),
			S3Endpoint:        getEnv("AWS_S3_ENDPOINT", ""), // For MinIO, DigitalOcean Spaces, etc.
		},
	}
	return cfg, nil
}

func mustEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("%s is required ", key)
	}
	return val
}

func getEnv(key, def string) string {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	return val
}

func getEnvAsInt(key string, def int) int {
	val := os.Getenv(key)
	if val == "" {
		return def
	}

	intVal := 0
	_, err := fmt.Sscanf(val, "%d", &intVal)
	if err != nil {
		return def
	}
	return intVal
}
