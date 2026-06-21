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
│   ├── agent/                      # LLM Router, ReAct Orchestrator & Native Go Validators
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
* `Circle (id STRING, name STRING, aliases STRING[], content STRING, verbatim STRING, history STRING[], created_at TIMESTAMP, needs_clarification BOOLEAN, embedding FLOAT[768])`
* `Task (id STRING, content STRING, aliases STRING[], verbatim STRING, status STRING, due_date TIMESTAMP, priority STRING, history STRING[], created_at TIMESTAMP, needs_clarification BOOLEAN, clarification_basis STRING, embedding FLOAT[768])`
* `Event (id STRING, content STRING, aliases STRING[], status STRING, event_date TIMESTAMP, history STRING[], created_at TIMESTAMP, needs_clarification BOOLEAN, clarification_basis STRING, embedding FLOAT[768])`
* `Insight (id STRING, content STRING, category STRING, confidence STRING, aliases STRING[], verbatim STRING, history STRING[], created_at TIMESTAMP, needs_clarification BOOLEAN, clarification_basis STRING, embedding FLOAT[768])`

### LadybugDB Edges
* `PARTICIPANT_IN` (FROM Person TO Event, role STRING)
* `MEMBER_OF` (FROM Person TO Circle, role STRING)
* `KNOWS` (FROM Person TO Person, descriptor STRING, context STRING)
* `LINKS_TO` (FROM Task TO Event, context STRING) *(Generic link)*

### SQLite Schema (`reminders.db`)
* `Reminders (id TEXT PRIMARY KEY, node_id TEXT, deadline DATETIME, is_sent BOOLEAN, message TEXT)`

---

## 3. Execution Pipeline (The Webhook Flow)
When a mock WAHA `curl` hits `/api/webhook`:
1. **Debounce:** The system aggregates messages into a `ConversationBlock`.
2. **Orchestrator Init:** The Go backend initializes the `llmRouter` (Gemini-primary) and prepares the tool definitions (`extract_transcript_manifest`, `query_rag`, `commit_mutations`).
3. **Agentic ReAct Loop:** The LLM reads the transcript, pulls a line-by-line manifest, and enters a thought loop. It autonomously queries the `query_rag` tool to resolve participant identities and detect existing events before drafting any DB operations.
4. **Structural Validation:** When the agent calls `commit_mutations`, the Go backend intercepts the JSON payload. A suite of multi-pass structural validators scrubs the mutations to ensure edge directionality, user resolution, and invariant compliance. Rejected edges are stripped.
5. **Commit:** The Go backend executes the surviving Cypher mutations in LadybugDB. Raw transcript provenance is handled via `verbatim` properties on nodes and `evidence_refs` (quotes) stored on edge structures, eliminating the need for a separate ConversationBlock graph node.

## 4. Verification Plan
1. Send mocked WAHA JSON payloads via `curl` to ensure basic HTTP endpoint integration.
2. Run the Deterministic Evaluation Harness (`go run cmd/eval/main.go`). This harness runs the agent against a fixed transcript fixture `N` times in a `dryRun` sandbox.
3. Validate the pass-rate table to ensure `10/10` adherence to Hard Invariants (no directional violations, identity integrity) and observe Soft Completeness metrics.
