package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	ServerPort     string
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	JWTSecret      string
	AWSRegion      string
	AWSAccessKeyID string
	AWSSecretKey   string
	AWSS3Bucket    string
}

func LoadConfig() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("No .env file found, relying on environment variables")
	}
	config := &Config{
		ServerPort:     getEnv("SERVER_PORT", "8080"),
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnv("DB_PORT", "5432"),
		DBUser:         getEnv("DB_USER", "postgres"),
		DBPassword:     getEnv("DB_PASSWORD", "uk888888"),
		DBName:         getEnv("DB_NAME", "social_media"),
		JWTSecret:      getEnv("JWT_SECRET", "HelloThisIsMyJWTSecret"),
		AWSRegion:      getEnv("AWS_REGION", "us-east-1"),
		AWSAccessKeyID: getEnv("AWS_ACCESS_KEY_ID", ""),
		AWSSecretKey:   getEnv("AWS_SECRET_ACCESS_KEY", ""),
		AWSS3Bucket:    getEnv("AWS_S3_BUCKET", ""),
	}

	return config
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
