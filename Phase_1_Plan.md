# Phase 1: Core Loop (Ingestion & Chat) - Implementation Plan

This document serves as the absolute "bible" for **Phase 1**. Any AI agent working on this phase must refer to this document for structural, architectural, and procedural constraints.

## Goal
Establish the full Core Loop, consisting of the end-to-end ingestion pipeline and the Alfred Chat query flow.
1. **Ingestion**: Alfred will receive WhatsApp messages (mocked via `curl`), run them through the Agentic Extraction pipeline, and correctly log the resulting `Task`, `Event`, and `Insight` nodes into LadybugDB and the `reminders.db` SQLite database.
2. **Alfred Chat**: Implement the PWA chat interface and the underlying `query_rag` loop, allowing Alfred to answer natural language questions and update memory mid-chat.

## 1. Repository Structure & Prompt Storage
To prevent multiple Go projects in the same workspace from confusing the IDE or agents, **we will nuke the `phase0-rag` sandbox completely**. We will initialize a single Go module at the project root (`/Alfred`). The sandbox logic will be migrated into `internal/rag`.

The prompts will be stored as modular markdown files in an `assets/prompts/` directory. These files will be loaded into memory at build time using Go's `//go:embed` directive to ensure the binary remains portable, atomic, and thread-safe.

```text
/Alfred (Root Go Module)
├── cmd/
│   └── alfred/
│       └── main.go                 # Entry point, HTTP server initialization
├── assets/
│   └── prompts/                    # Embedded with //go:embed
│       ├── core_persona.md         # Alfred's personality and tone constraints
│       ├── core_schema.md          # Topology, nodes, and edges constraints
│       ├── skill_ingestion.md      # Rules for webhook manifest extraction
│       └── skill_chat.md           # Rules for natural language Q&A and history fetching
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

6. **Temporal Maintenance:** The Go backend natively intercepts `UPDATE_NODE` operations. Before updating the node, it executes `MATCH (n) WHERE n.id = $id RETURN n.content`. It formats this old state as `"YYYY-MM-DD HH:MM - [old_content]"` and prepends it to the `history STRING[]` array, then executes the `SET` for the new content. This offloads history maintenance from the LLM and prevents data loss.
7. **SQLite Reminder Synchronization:** If the agent commits a `Task` mutation with a `due_date`, the orchestrator automatically executes `INSERT OR REPLACE INTO reminders` against the SQLite DB so the cron job can pick it up.

## 4. Prompt Composition & Modularization
To preserve Gemini Context Caching while keeping prompts highly relevant, the Go backend will compose prompts dynamically at the *pipeline level* (not mid-loop).
* **Ingestion Pipeline:** Go concatenates `core_persona.md` + `core_schema.md` + `skill_ingestion.md` once, before the agent starts.
* **Chat Pipeline:** Go concatenates `core_persona.md` + `core_schema.md` + `skill_chat.md` once, before the agent starts.
This ensures the prefix remains absolutely static during the 5-6 `query_rag` calls within the agent loop, slashing token costs while maximizing attention accuracy.

## 5. Execution Pipeline (The Chat Flow)
The second half of the Core Loop. The user interacts via a PWA frontend connected to an `/api/chat` Go endpoint. This flow is entirely non-linear.
1. **Agent Toolkit**: The Chat Agent operates dynamically using the following tools:
   * `query_rag`: Returns `Nodes` (id, type, content) and `Edges`. History is intentionally omitted to save tokens.
   * `query_node_history(node_id)`: Fetches the full `history STRING[]` changelog for a specific node. **Agent Autonomy Rule**: The LLM is explicitly instructed to call this autonomously if it reads a node's self-signaling `content` and determines its internal reasoning requires precise temporal context to resolve ambiguity. It does not wait for a user's prompt to use this.
   * `create_node`, `update_node`, `delete_node`: For mid-chat graph mutations.
   * `upsert_reminder`, `check_reminders`: For managing SQLite triggers.
   * `ask_user_for_hint`: Yields to the user.

## 6. Verification Plan
1. Send mocked WAHA JSON payloads via `curl` to ensure basic HTTP endpoint integration.
2. Run the Deterministic Evaluation Harness (`go run cmd/eval/main.go`). This harness must execute *actual* Cypher queries against LadybugDB (replacing `mock.go`).
3. Validate the pass-rate table to ensure `10/10` adherence to Hard Invariants.
4. Test the PWA Chat interface by asking complex temporal questions that force the agent to autonomously trigger `query_node_history`.

## 7. Execution Roadmap (The Sub-Phases)
To maintain quality across this massive architectural shift, Phase 1 execution is strictly divided into the following sub-phases. We will not move to the next sub-phase until the current one is tested and verified.

### Sub-Phase 1.1: Modular Prompt Refactoring
- [ ] Split `assets/prompts/ingestion_agent.md` into `core_persona.md`, `core_schema.md`, `skill_ingestion.md`, and `skill_chat.md`.
- [ ] Refactor the Go backend to load these files into memory via `//go:embed` instead of one massive monolithic file.
- [ ] Update the `Orchestrator` to deterministically concatenate the Ingestion prompt before initiating the LLM call.

### Sub-Phase 1.2: Real Database Commits (LadybugDB)
- [ ] Excisé `AddMockNode` and `AddMockEdge` from `internal/ladybug/mock.go`.
- [ ] Implement parameterized/sanitized Cypher string generation in `orchestrator.go` for `CREATE_NODE` and `MATCH... CREATE` edges.
- [ ] Execute real mutations via `o.DBConn.Query()`.
- [ ] Run `cmd/eval/main.go` to verify the DB engine accepts the real queries without syntax failure.

### Sub-Phase 1.3: Temporal Update Logistics
- [ ] Intercept `UPDATE_NODE` operations within the Go orchestrator.
- [ ] Execute a `MATCH (n) WHERE n.id = $id RETURN n.content` query.
- [ ] Format the timestamp and prepend the old content to the `history STRING[]` array.
- [ ] Execute the final `SET n.content = $new, n.history = $history` Cypher query.

### Sub-Phase 1.4: SQLite Reminders Integration
- [ ] Pass the SQLite `*sql.DB` connection into the `Orchestrator` struct during initialization.
- [ ] Add an intercept in the orchestrator loop: if a `Task` mutation is committed and contains a `due_date`, write it to `reminders.db`.
- [ ] Execute `INSERT OR REPLACE INTO reminders (id, node_id, deadline, is_sent, message)`.

### Sub-Phase 1.5: The Chat Agent (Backend)
- [ ] Create an `/api/chat` POST endpoint in `cmd/alfred/main.go`.
- [ ] Build the `RunChatAgent()` loop in `internal/agent/chat.go`.
- [ ] Implement the `query_node_history(node_id)` tool to fetch the `history STRING[]` array natively.
- [ ] Equip the Chat Agent with its full toolkit and ensure the prompt composition logic stitches `core_persona` + `core_schema` + `skill_chat`.

### Sub-Phase 1.6: The PWA Chat Interface (Frontend)
- [ ] Update `public/index.html` to include a chat overlay alongside the Graph viewer.
- [ ] Implement HTTP fetching, message rendering, and loading state UI in `public/app.js`.
- [ ] Apply high-fidelity, dynamic CSS styling (glassmorphism, smooth transitions) to ensure a premium user experience.
