package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	GeminiAPIKey string
	GroqAPIKey   string
}

func LoadConfig() *Config {
	err := godotenv.Overload()
	if err != nil {
		log.Println("Warning: No .env file found, falling back to system environment variables")
	}

	geminiKey := os.Getenv("GEMINI_API_KEY")
	if geminiKey == "" {
		log.Fatal("GEMINI_API_KEY is missing from environment")
	}

	groqKey := os.Getenv("GROQ_API_KEY")
	if groqKey == "" {
		log.Fatal("GROQ_API_KEY is missing from environment")
	}

	return &Config{
		GeminiAPIKey: geminiKey,
		GroqAPIKey:   groqKey,
	}
}
