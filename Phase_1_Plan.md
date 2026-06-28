# Phase 1: Core Loop (Ingestion & Chat) - Implementation Plan

This document serves as the absolute "bible" for **Phase 1**. Any AI agent working on this phase must refer to this document for structural, architectural, and procedural constraints.

## Goal
Establish the full Core Loop, consisting of the end-to-end ingestion pipeline and the Alfred Chat query flow.
1. **Ingestion**: Alfred will receive WhatsApp messages (mocked via `curl`), run them through the Agentic Extraction pipeline, and correctly log the resulting `Task`, `Event`, and `Insight` nodes into LadybugDB and the `reminders.db` SQLite database.
2. **Alfred Chat**: Implement the PWA chat interface and the underlying `query_rag` loop, allowing Alfred to answer natural language questions and update memory mid-chat.

## 1. Repository Structure & Prompt Storage
To prevent multiple Go projects in the same workspace from confusing the IDE or agents, **we will nuke the `phase0-rag` sandbox completely**. We will initialize a single Go module at the project root (`/Alfred`). The sandbox logic will be migrated into `internal/rag`.

The prompts will be stored as modular markdown files in an `assets/prompts/` directory. These files will be loaded into memory at build-time using Go's `//go:embed` directive. This ensures atomic, crash-proof binaries while keeping the actual text files completely modular.

```text
/Alfred (Root Go Module)
├── cmd/
│   └── alfred/
│       └── main.go                 # Entry point, HTTP server initialization
├── assets/
│   └── prompts/                    # Modular prompt components embedded with //go:embed
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

### LadybugDB Nodes (Mocked in memory)
Based on the Decision Log, all nodes (except `Person`) get a `content`, `history`, `created_at`, `aliases`, and `embedding` field.
* `Person (id STRING, name STRING, phone_number STRING, aliases STRING[], created_at TIMESTAMP, needs_clarification BOOLEAN)`
* `Circle (id STRING, name STRING, aliases STRING[], title STRING, content STRING, verbatim STRING, history STRING[], created_at TIMESTAMP, needs_clarification BOOLEAN, embedding FLOAT[768])` Note: No longer created inline by the ingestion agent. Created via Layer 2 batch promotion.
* `Task (id STRING, content STRING, aliases STRING[], verbatim STRING, group_mentions STRING, status STRING, due_date TIMESTAMP, priority STRING, history STRING[], created_at TIMESTAMP, needs_clarification BOOLEAN, clarification_basis STRING, embedding FLOAT[768])`
* `Event (id STRING, content STRING, aliases STRING[], status STRING, group_mentions STRING, event_date TIMESTAMP, history STRING[], created_at TIMESTAMP, needs_clarification BOOLEAN, clarification_basis STRING, embedding FLOAT[768])`
* `Insight (id STRING, content STRING, category STRING, confidence STRING, aliases STRING[], verbatim STRING, history STRING[], created_at TIMESTAMP, needs_clarification BOOLEAN, clarification_basis STRING, embedding FLOAT[768])`

### LadybugDB Edges
* `ASSIGNED_TO` (FROM Person TO Task)
* `MENTIONED_IN` (FROM Person TO Task/Event)
* `HAS_ROLE` (FROM Person TO Event)
* `MEMBER_OF` (FROM Person TO Circle, role STRING, since TIMESTAMP)
* `PART_OF` (FROM Task TO Event/Circle)
* `DIR_TOWARDS` (FROM Insight TO Person/Circle)
* `LINKS_TO` (Universal generic link, context STRING)
* `CONTRADICTS` (FROM Insight TO Insight, detected_at TIMESTAMP, resolved BOOLEAN)
* `KNOWS` (FROM Person TO Person, descriptor STRING, context STRING, since TIMESTAMP)

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

## 4. Prompt Composition & Modularization (Additive Dynamic ReAct)
To mathematically prevent "Lost in the Middle" instruction decay without completely breaking Google's Context Caching, Alfred uses an **Additive-with-Pruning** dynamic injection architecture within the ReAct loop.

1. **Phase 1 (Discovery):** The Go backend initiates the LLM loop with only `core_persona.md` + `skill_discovery.md`. The LLM's attention is 100% focused on extracting entities and querying RAG.
2. **The Trigger:** The `skill_discovery` prompt explicitly instructs the agent to output `[REQUEST_SCHEMA]` in its thought block when it has finished querying the vault.
3. **The Interceptor (Prune & Append):** Go intercepts the ReAct loop mid-flight. When it detects `[REQUEST_SCHEMA]`, it sweeps the `history` slice and DELETES any previous message tagged with `[SYSTEM_INJECTION_SKILL_COMMIT]`. It then APPENDS a fresh `User` message to the end of the history containing the massive 2,000-token `core_schema.md` constraints.
4. **Phase 2 (Topology):** The LLM receives the prompt with the strict rules placed at the absolute bottom of its context window (maximizing Recency Bias without payload bloat).

Defensive guardrails inside the Go Orchestrator (`internal/agent/orchestrator.go`) intercept the agent's context and strictly enforce a state machine before schema injection is permitted:
- **Gate 1 (Manifest Validation):** Blocks `[REQUEST_SCHEMA]` if `extract_transcript_manifest` has not been called.
- **Gate 2 (Speaker Resolution):** Blocks if there are speakers in the manifest that have not been resolved via `query_rag` or explicitly declared new. The `query_rag` tool enforces a strict 1:1 `target_speakers` array matching the length of the `queries` array.
- **Gate 3 (Temporal Obligations Check):** Blocks if the agent has not called the `query_speaker_obligations` tool for its resolved speakers. This forces the agent to check for existing `needs_clarification: true` nodes to perform temporal updates (`UPDATE_NODE`) instead of creating duplicate nodes.

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

### Sub-Phase 1.1: Modular Prompt Refactoring & Dynamic ReAct (COMPLETED)
- [x] Split `assets/prompts/ingestion_agent.md` into `core_persona.md`, `core_schema.md`, `skill_discovery.md`, `skill_commit.md`, and `skill_chat.md`.
- [x] Integrate the Interceptor pattern into `internal/llm/router.go` to capture the `[REQUEST_SCHEMA]` string in the agent's thought process.
- [x] Implement the `Additive-with-Pruning` history mutation logic in `internal/agent/orchestrator.go` to dynamically inject the `skill_commit.md` rules at the peak of the context window.
- [x] Implement Anti-Premature and Anti-Forgetful state guardrails to prevent discovery-bypassing and schema-skipping.
- [x] Harden prompt constraints: Enforce STRICT DEFAULT on 5W Clarity Checks (removing the 'operationally necessary' loophole) and mandate two unique explicit keywords for Event Inference to prevent RAG-bias hallucination.

### Sub-Phase 1.2: Mock Database Pivot & Schema Guardrails (COMPLETED)
- `[x]` Pause CGO/LadybugDB integration due to VPS compilation blockers. Revert `internal/ladybug/mock.go` to handle fully parsed graph state in memory.
- `[x]` Implement Layer 1 Circle Deferral: remove `Circle` node creation from ingestion agent and add `group_mentions` property to Tasks/Events.
- `[x]` Implement Gate 5 in `internal/agent/tool_handlers.go` to mechanically hard-reject inline Circle creations.
- `[x]` Harden JSON parsing via `DisallowUnknownFields()` and Rule 22 to eliminate property-nesting hallucinations.
- `[x]` Restructure test suite into domain-specific stress point folders (e.g., `01_core_extraction`, `03_advanced_hubbing`).

### Sub-Phase 1.3: Temporal Update Logistics (COMPLETED)
- `[x]` Intercept `UPDATE_NODE` operations within the Go orchestrator.
- `[x]` Execute a `MATCH (n) WHERE n.id = $id RETURN n.content` query. (Superseded by Atomic Cypher)
- `[x]` Format the timestamp and prepend the old content to the `history STRING[]` array. (Superseded by Atomic Cypher)
- `[x]` Execute the final `SET n.content = $new, n.history = $history` Cypher query. (Implemented via Atomic Cypher in `execution.go`)

### Sub-Phase 1.4: SQLite Reminders Integration (COMPLETED)
- `[x]` Pass the SQLite `*sql.DB` connection into the `Orchestrator` struct during initialization.
- `[x]` Add an intercept in the orchestrator loop: if a `Task` mutation is committed and contains a `due_date`, write it to `reminders.db`.
- `[x]` Execute `INSERT OR REPLACE INTO reminders (id, node_id, deadline, is_sent, message)`.

### Sub-Phase 1.5: The Chat Agent (Backend)
- [ ] Create an `/api/chat` POST endpoint in `cmd/alfred/main.go`.
- [ ] Build the `RunChatAgent()` loop in `internal/agent/chat.go`.
- [ ] Implement the `query_node_history(node_id)` tool to fetch the `history STRING[]` array natively.
- [ ] Equip the Chat Agent with its full toolkit and ensure the prompt composition logic stitches `core_persona` + `core_schema` + `skill_chat`.

### Sub-Phase 1.6: Dev Chat Interface (Frontend)
- [ ] Update `public/index.html` to include a mock chat overlay alongside the Graph viewer.
- [ ] Implement HTTP fetching, message rendering, and loading state UI in `public/app.js` to test the backend API.

### Sub-Phase 1.7: Layer 2 Mention Promotion (Cron Job)
- [ ] Create a standalone Go routine or separate binary that runs periodically.
- [ ] Scan `Task` and `Event` nodes for unpromoted `group_mentions` payloads.
- [ ] Call the Gemini API with a specialized clustering prompt to group aliases (e.g., "divisi logistik", "div logistik") into a single canonical `Circle` node.
- [ ] Execute `CREATE_NODE` for the new `Circle` and link the participants via `MEMBER_OF` edges.
