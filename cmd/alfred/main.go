package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gimigkk/Alfred-Memory/internal/agent"
	"github.com/gimigkk/Alfred-Memory/internal/config"
	"github.com/gimigkk/Alfred-Memory/internal/db"
	"github.com/gimigkk/Alfred-Memory/internal/embed"
	"github.com/gimigkk/Alfred-Memory/internal/llm"
	"github.com/gimigkk/Alfred-Memory/internal/waha"
)

func main() {
	// Remove timestamp from standard logger
	log.SetFlags(0)

	log.Println("Alfred Core Ingestion Pipeline Starting...")

	// 1. Load config & API clients
	cfg := config.LoadConfig()
	geminiEmbed := embed.NewGeminiClient(cfg.GeminiAPIKey)
	llmRouter := llm.NewRouterClient(cfg.GeminiAPIKey, cfg.GroqAPIKey)

	// 2. Initialize DBs
	dbDir := "./.lbug"
	_ = os.MkdirAll(dbDir, 0755)

	lbugClient, err := db.NewClient(dbDir)
	if err != nil {
		log.Fatalf("Failed to initialize LadybugDB: %v", err)
	}
	defer lbugClient.Close()

	conn, err := lbugClient.GetConnection()
	if err != nil {
		log.Fatalf("Failed to get connection: %v", err)
	}
	defer conn.Close()

	log.Println("Initializing LadybugDB Schema...")
	if err := db.InitLadybugSchema(conn); err != nil {
		log.Printf("Schema init warning: %v", err)
	}

	sqliteDB, err := db.InitSQLite("./reminders.db")
	if err != nil {
		log.Fatalf("Failed to init SQLite: %v", err)
	}
	defer sqliteDB.Close()

	// 3. Setup Orchestrator
	orchestrator := agent.NewOrchestrator(llmRouter, geminiEmbed, conn)

	// Register webhook callback
	waha.OnBlockCommitted = func(block *waha.ConversationBlock) {
		go func() {
			_, err := orchestrator.RunAgenticIngestion(block.ID, block.FormatTranscript(), false)
			if err != nil {
				log.Printf("❌ Failed to process block %s: %v", block.ID, err)
			}
		}()
	}

	// 4. Start Server
	http.HandleFunc("/api/webhook", waha.WebhookHandler)

	http.HandleFunc("/api/vault", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		
		// 1. Fetch all nodes
		resNodes, err := conn.Query("MATCH (n) RETURN n.id, label(n), n.content, n.properties")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		var nodes []map[string]any
		for resNodes.HasNext() {
			row := resNodes.GetNext()
			nodeID := row[0].(string)
			nodeType := row[1].(string)
			content := row[2].(string)

			var properties map[string]any
			if len(row) > 3 && row[3] != nil {
				if props, ok := row[3].(map[string]any); ok {
					properties = props
				}
			}
			if properties == nil {
				properties = make(map[string]any)
			}

			// Ensure content is populated in properties if not already present
			if _, ok := properties["content"]; !ok && content != "" && nodeType != "Person" {
				properties["content"] = content
			}

			// Use name or content as label
			label := content
			if nodeType == "Person" {
				if name, ok := properties["name"].(string); ok && name != "" {
					label = name
				} else {
					label = nodeID
				}
			}
			if len(label) > 40 {
				label = label[:37] + "..."
			}

			nodes = append(nodes, map[string]any{
				"id":         nodeID,
				"label":      label,
				"group":      nodeType,
				"properties": properties,
			})
		}
		resNodes.Close()

		// 2. Fetch all edges
		resEdges, err := conn.Query("MATCH (a)-[r]->(b) RETURN a.id, b.id, label(r)")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		var edges []map[string]string
		for resEdges.HasNext() {
			row := resEdges.GetNext()
			from := row[0].(string)
			to := row[1].(string)
			relType := row[2].(string)
			edges = append(edges, map[string]string{
				"id":    fmt.Sprintf("%s_%s_%s", from, relType, to),
				"from":  from,
				"to":    to,
				"label": relType,
			})
		}
		resEdges.Close()

		json.NewEncoder(w).Encode(map[string]any{
			"nodes": nodes,
			"edges": edges,
		})
	})

	http.Handle("/", http.FileServer(http.Dir("./public")))

	port := "8080"
	log.Printf("WAHA Webhook receiver listening on port %s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
