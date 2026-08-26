package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	CognoDBURI      string
	CognoDBUsername string
	CognoDBPassword string
	Port            string
}

func Load() *Config {
	_ = godotenv.Load()
	_ = godotenv.Load("../.env")
	_ = godotenv.Load("../../.env")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	// Docs: URI is bolt+s://..., auth is ("cognodb", DB_PASSWORD). Username defaults to cognodb per docs if not set.
	uri := os.Getenv("COGNODB_URI")
	username := os.Getenv("COGNODB_USERNAME")
	if username == "" {
		username = os.Getenv("COGNODB_USER")
	}
	if username == "" && uri != "" {
		username = "cognodb"
	}
	password := os.Getenv("COGNODB_PASSWORD")
	if password == "" {
		password = os.Getenv("DB_PASSWORD")
	}
	return &Config{
		CognoDBURI:      uri,
		CognoDBUsername: username,
		CognoDBPassword: password,
		Port:            port,
	}
}

func (c *Config) HasDBConfig() bool {
	return c.CognoDBURI != "" && c.CognoDBPassword != ""
}
