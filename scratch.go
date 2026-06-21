package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"os"
)

func main() {
	cModels := []string{}
	// Fetch Groq Models
	reqGroq, _ := http.NewRequest("GET", "https://api.groq.com/openai/v1/models", nil)
	reqGroq.Header.Set("Authorization", "Bearer " + os.Getenv("GROQ_API_KEY"))
	client := &http.Client{}
	respGroq, err := client.Do(reqGroq)
	if err == nil {
		defer respGroq.Body.Close()
		var resGroq struct {
			Data []struct {
				ID string `json:"id"`
			} `json:"data"`
		}
		if json.NewDecoder(respGroq.Body).Decode(&resGroq) == nil {
			for _, m := range resGroq.Data {
				name := m.ID
				if !strings.Contains(name, "whisper") && !strings.Contains(name, "guard") && !strings.Contains(name, "safeguard") && !strings.Contains(name, "orpheus") {
					cModels = append(cModels, "groq/"+name)
				}
			}
		} else {
            fmt.Println("Failed to decode Groq JSON")
        }
	} else {
        fmt.Println("Groq HTTP Error:", err)
    }
    fmt.Println("Groq Models:", len(cModels))
    for _, m := range cModels {
        fmt.Println(m)
    }
}
