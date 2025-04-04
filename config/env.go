package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port       string
	DBUser     string
	DBPassword string
	DBAddress  string
	DBName     string
}

func initDatabaseConfig() Config {
	godotenv.Load()

	return Config{
		Port:       getEnv("PORT", "8080"),
		DBUser:     getEnv("DB_USER", ""),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBAddress:  fmt.Sprintf("%s:%s", getEnv("DB_HOST", ""), getEnv("DB_PORT", "")),
		DBName:     getEnv("DB_NAME", "ecom_bot"),
	}
}

type GoogleConfig struct {
	GoogleClientID     string
	GoogleClientSecret string
	GCPProjectID       string
	GCBearerToken	  	string
}

type GeminiConfig struct {
	GeminiAPIKey string
}

type ShippoConfig struct {
	ShippoAPIKey string
}

type JWTConfig struct {
	JWTSecret string
}

type SMTPEmailConfig struct {
	SMTPServer   string
	SMTPPort     string
	SenderEmail  string
	AppPassword  string
}

func initSMTPEmailConfig() SMTPEmailConfig {
	godotenv.Load()
	return SMTPEmailConfig{
		SMTPServer:   getEnv("SMTP_SERVER", ""),
		SMTPPort:     getEnv("SMTP_PORT", ""),
		SenderEmail:  getEnv("SENDER_EMAIL", ""),
		AppPassword:  getEnv("APP_PASSWORD", ""),
	}
}

func initGoogleConfig() GoogleConfig {
	godotenv.Load()
	return GoogleConfig{
		GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
		GCPProjectID:       getEnv("GCP_PROJECT_ID", ""),
		GCBearerToken:      getEnv("GC_BEARER_TOKEN", ""),
	}
}

func initShippoConfig() ShippoConfig {
	godotenv.Load()

	return ShippoConfig{
		ShippoAPIKey: getEnv("SHIPPO_API_KEY", ""),
	}
}

func initGeminiAPIConfig() GeminiConfig {
	godotenv.Load()
	return GeminiConfig{
		GeminiAPIKey: getEnv("GEMINI_API_KEY", ""),
	}
}

func initJWTConfig() JWTConfig {
	godotenv.Load()
	return JWTConfig{
		JWTSecret: getEnv("JWT_SECRET", ""),
	}
}

func getEnv(key, fallback string) string {

	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func getEnvAsInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}
		return i

	}
	return fallback
}

var Envs = initDatabaseConfig()
var GoogleEnvs = initGoogleConfig()
var GeminiEnvs = initGeminiAPIConfig()
var ShippoEnvs = initShippoConfig()
var JWTEnvs = initJWTConfig()
var SMTPEnvs = initSMTPEmailConfig()
