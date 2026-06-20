package waha

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

// WahaPayload represents the incoming JSON payload from WAHA
type WahaPayload struct {
	Event   string `json:"event"`
	Payload struct {
		ID          string `json:"id"`
		From        string `json:"from"`
		To          string `json:"to"`
		Participant string `json:"participant"`
		FromMe      bool   `json:"fromMe"`
		Body        string `json:"body"`
		Timestamp   int64  `json:"timestamp"`
	} `json:"payload"`
}

type ConversationBlock struct {
	ID        string
	ChatID    string
	Messages  []WahaPayload
	CreatedAt time.Time
}

var (
	blocks   = make(map[string]*ConversationBlock)
	blocksMu sync.Mutex
)

// OnBlockCommitted is a callback that the main application registers
// to handle a block once the debounce timer expires.
var OnBlockCommitted func(block *ConversationBlock)

// WebhookHandler handles the incoming WAHA webhooks
func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload WahaPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if payload.Event != "message" && payload.Event != "message.any" {
		w.WriteHeader(http.StatusOK)
		return
	}

	chatID := payload.Payload.From
	if payload.Payload.FromMe {
		chatID = payload.Payload.To
	}

	blocksMu.Lock()
	block, exists := blocks[chatID]
	if !exists {
		block = &ConversationBlock{
			ID:        fmt.Sprintf("block_%d", time.Now().UnixNano()),
			ChatID:    chatID,
			Messages:  []WahaPayload{},
			CreatedAt: time.Now(),
		}
		blocks[chatID] = block

		// Start debounce timer
		go func(cid string) {
			// In production this would be 15 minutes, but for testing we use 2 seconds
			time.Sleep(2 * time.Second)
			blocksMu.Lock()
			committedBlock := blocks[cid]
			delete(blocks, cid)
			blocksMu.Unlock()

			if committedBlock != nil && OnBlockCommitted != nil {
				log.Printf("🚀 [Step 3] Pushing compressed block %s to Ingestion Pipeline...", committedBlock.ID)
				OnBlockCommitted(committedBlock)
			}
		}(chatID)
	}

	block.Messages = append(block.Messages, payload)
	log.Printf("📥 [Step 1] Message %s queued for chat %s. Current queue size: %d", payload.Payload.ID, chatID, len(block.Messages))
	blocksMu.Unlock()

	w.WriteHeader(http.StatusOK)
}

// FormatTranscript formats the block's messages into a raw transcript string for the LLM
func (b *ConversationBlock) FormatTranscript() string {
	log.Printf("🗜️ [Step 2] Compressing %d queued messages into a single transcript block...", len(b.Messages))
	var transcript string
	for _, m := range b.Messages {
		sender := m.Payload.From
		if m.Payload.FromMe {
			sender = "THE USER"
		} else if m.Payload.Participant != "" {
			sender = m.Payload.Participant
		}
		
		// Strip WAHA suffixes like @c.us or @g.us
		if len(sender) > 5 && (sender[len(sender)-5:] == "@c.us" || sender[len(sender)-5:] == "@g.us") {
			sender = sender[:len(sender)-5]
		}
		
		if strings.Contains(strings.ToLower(sender), "gilang") {
			sender = "THE USER"
		}
		
		body := m.Payload.Body
		re := regexp.MustCompile(`(?i)(@)?(m3-117_)?gilang( muhamad w)?`)
		body = re.ReplaceAllString(body, "@THE USER")
		
		t := time.Unix(m.Payload.Timestamp, 0).Format("2006-01-02 15:04:05")
		transcript += fmt.Sprintf("[%s][%s]: %s\n", t, sender, body)
	}
	return transcript
}
