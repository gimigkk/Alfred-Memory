# Phase 1: Core Ingestion Pipeline - Implementation Plan

This document serves as the absolute "bible" for **Phase 1**. Any AI agent working on this phase must refer to this document for structural, architectural, and procedural constraints.

## Goal
Establish the end-to-end ingestion pipeline. Alfred will receive WhatsApp messages (mocked via `curl`), run them through the Two-Phase Extraction pipeline, and correctly log the resulting `Task`, `Event`, and `Insight` nodes into LadybugDB and the `reminders.db` SQLite database.

## 1. Repository Structure & Prompt Storage
To prevent multiple Go projects in the same workspace from confusing the IDE or agents, **we will nuke the `phase0-rag` sandbox completely**. We will initialize a single Go module at the project root (`/Alfred`). The sandbox logic will be migrated into `internal/rag`.

The prompts will be stored as markdown files in an `assets/prompts/` directory using Go's `//go:embed` directive.

```text
/Alfred (Root Go Module)
├── cmd/
│   └── alfred/
│       └── main.go                 # Entry point, HTTP server initialization
├── assets/
│   └── prompts/                    # Embedded with //go:embed
│       ├── extraction_skill.md     # Persona & Phase 1 rules
│       └── linking_skill.md        # Identity Linking rules
├── internal/
│   ├── agent/                      # Groq client, Extraction/Linking Orchestrator
│   ├── config/                     # Environment loading (.env)
│   ├── db/                         # LadybugDB & SQLite managers + Schema Setup
│   ├── embed/                      # Gemini Client for vector generation
│   ├── rag/                        # Migrated hybrid search math (Phase 0 logic)
│   └── waha/                       # Webhook handlers and JSON payload models
├── .env                            # API keys
└── go.mod                          # Single project-wide Go module
```

---

## 2. Database Schema (DDL)

### LadybugDB Nodes
Based on the Decision Log, all nodes (except `Person`) get a `content`, `history`, `created_at`, `aliases`, and `embedding` field.
* `Person (id STRING, name STRING, phone_number STRING, aliases STRING[], created_at TIMESTAMP, needs_clarification BOOLEAN)`
* `Circle (id STRING, name STRING, aliases STRING[], content STRING, history STRING[], created_at TIMESTAMP, needs_clarification BOOLEAN, embedding FLOAT[768])`
* `Task (id STRING, content STRING, aliases STRING[], verbatim STRING, status STRING, due_date TIMESTAMP, history STRING[], created_at TIMESTAMP, needs_clarification BOOLEAN, embedding FLOAT[768])`
* `Event (id STRING, content STRING, aliases STRING[], status STRING, start_date TIMESTAMP, history STRING[], created_at TIMESTAMP, needs_clarification BOOLEAN, embedding FLOAT[768])`
* `Insight (id STRING, content STRING, aliases STRING[], verbatim STRING, history STRING[], created_at TIMESTAMP, needs_clarification BOOLEAN, embedding FLOAT[768])`
* `ConversationBlock (id STRING, chat_id STRING, raw_transcript STRING, created_at TIMESTAMP)`

### LadybugDB Edges
* `PARTICIPANT_IN (FROM Person TO Event, role STRING)`
* `MEMBER_OF (FROM Person TO Circle, role STRING)`
* `KNOWS (FROM Person TO Person, descriptor STRING, context STRING)`
* `CAUSED_BY (FROM Task TO ConversationBlock, context STRING)`
* `EVIDENCED_BY (FROM Insight TO ConversationBlock, context STRING)`
* `LINKS_TO (FROM Task TO Event, context STRING)` *(Generic link)*

### SQLite Schema (`reminders.db`)
* `Reminders (id TEXT PRIMARY KEY, node_id TEXT, deadline DATETIME, is_sent BOOLEAN, message TEXT)`

---

## 3. Execution Pipeline (The Webhook Flow)
When a mock WAHA `curl` hits `/api/webhook`:
1. **Debounce:** The system aggregates messages into a `ConversationBlock`.
2. **Phase 1 (Blind Extract):** Sends the `raw_transcript` + `extraction_skill.md` to Groq. Groq outputs candidate JSON nodes (no graph access).
3. **RAG Intermediary:** The Go backend automatically embeds the candidates and runs `query_rag` against LadybugDB to fetch relevant graph context (like existing Persons or Events).
4. **Phase 2 (Link):** Sends the Candidates + Graph Context + `linking_skill.md` to Groq. Groq outputs final `CREATE_NODE` or `UPDATE_NODE` mutations.
5. **Commit:** The Go backend executes the Cypher mutations in LadybugDB.

## 4. Verification Plan
1. Initialize the new root module and delete `phase0-rag`.
2. Send a mocked WAHA JSON payload via `curl` containing a conversation where "Bahlil" asks the user for a "DPP design by Friday".
3. Query LadybugDB to ensure a `Person` node for Bahlil, a `Task` node for the design, and a `ConversationBlock` node were created and correctly linked via `CAUSED_BY`.
4. Check the SQLite database to ensure the Friday deadline was logged in `reminders.db`.
