package config

import (
	"os"
	"github.com/joho/godotenv"
)

type Config struct {
	AppPort       string
	DBHost        string
	DBPort        string
	DBUser        string
	DBPass        string
	DBName        string
	JWTSecretKey  string
	MQURL         string
	MQVHost       string
}

func LoadConfig() *Config {
	_ = godotenv.Load()
	return &Config{
		AppPort:      getEnv("APP_PORT", "10001"),
		DBHost:       getEnv("DB_HOST", "0.0.0.0"),
		DBPort:       getEnv("DB_PORT", "5432"),
		DBUser:       getEnv("DB_USER", "postgres"),
		DBPass:       getEnv("DB_PASS", "postgres"),
		DBName:       getEnv("DB_NAME", "community"),
		JWTSecretKey: getEnv("JWT_SECRET_ACCESS_KEY", "esto es una key re segura"),
		MQURL:        getEnv("MQ_URL", "amqp://guest:guest@localhost:5672/"),
		MQVHost:      getEnv("MQ_VHOST", "/"),
	}
}

func getEnv(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return def
}
