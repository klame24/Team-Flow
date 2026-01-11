package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
	URL      string
}

type JWTConfig struct {
	SecretKey       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

func init() {
	_ = godotenv.Load("env.env")
}

func Load() (*Config, error) {

	return &Config{
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", "8080"),
			ReadTimeout:  getEnvAsDuration("SERVER_READ_TIMEOUT", 10*time.Second),
			WriteTimeout: getEnvAsDuration("SERVER_WRITE_TIMEOUT", 10*time.Second),
			IdleTimeout:  getEnvAsDuration("SERVER_IDLE_TIMEOUT", 60*time.Second),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "root"),
			DBName:   getEnv("DB_NAME", "postgres"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
			URL:      buildDatabaseURL(),
		},
		JWT: JWTConfig{
			SecretKey:       getEnv("JWT_SECRET", "your-secret-key"),
			AccessTokenTTL:  getEnvAsDuration("JWT_ACCESS_TTL", 15*time.Minute),
			RefreshTokenTTL: getEnvAsDuration("JWT_REFRESH_TTL", 30*24*time.Hour),
		},
	}, nil
}

func buildDatabaseURL() string {

	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "root")
	dbname := getEnv("DB_NAME", "postgres")
	sslmode := getEnv("DB_SSL_MODE", "disable")

	return "postgresql://" + user + ":" + password + "@" + host + ":" + port + "/" + dbname + "?sslmode=" + sslmode
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
