package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	GeminiAPIKey string
	GroqAPIKey   string
	OwnerID      string
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

	ownerID := os.Getenv("OWNER_ID")
	if ownerID == "" {
		ownerID = "USER_GILANG"
	}

	return &Config{
		GeminiAPIKey: geminiKey,
		GroqAPIKey:   groqKey,
		OwnerID:      ownerID,
	}
}
