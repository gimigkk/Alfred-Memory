## Agentic Secretary AI — Full Project Specification

> **Status:** Planning / Pre-development  
> **Last updated from chat:** June 12, 2026  
> **Stack:** WAHA (GOWS), Flutter, Budget VPS (4GB RAM / 8GB storage), Golang Backend, Groq API + open-weight LLM fallback, LadybugDB (embedded graph DB)

---

## Table of Contents

1. [Project Vision](#1-project-vision)
2. [System Architecture Overview](#2-system-architecture-overview)
3. [The "Alfred" Point of View (Persona)](#3-the-alfred-point-of-view-persona)
4. [Ingestion Layer (WAHA / WhatsApp)](#4-ingestion-layer)
5. [Conversation Block System](#5-conversation-block-system)
6. [Memory Vault — Graph Database](#6-memory-vault--graph-database)
7. [Memory Node Schema & Templates](#7-memory-node-schema--templates)
8. [Agentic Tooling System (The Traversal Loop)](#8-agentic-tooling-system)
9. [LLM Extraction Pipeline](#9-llm-extraction-pipeline)
10. [Storage & Purging Strategy](#10-storage--purging-strategy)
11. [Chat Interface & Observability Layer](#11-chat-interface--observability-layer)
12. [Infrastructure & Stack Decisions](#12-infrastructure--stack-decisions)
13. [Open Questions & Undiscussed Aspects](#13-open-questions--undiscussed-aspects)
14. [Decision Log](#14-decision-log)

---

## 1. Project Vision

An **agentic secretary AI** that passively watches your conversations — starting with WhatsApp — and turns raw chat into a **living, temporal memory system**. It is not a chat archive. It is a structured knowledge graph that grows over time and can be queried naturally.

The final product should be able to:
- **Remember** — extract and store structured facts, tasks, events, preferences, people, experiences, and social insights
- **Forget intentionally** — raw chat is ephemeral; only curated memory is permanent
- **Summarize** — compress conversation blocks into semantic summaries
- **Merge duplicates** — resolve conflicts between overlapping pieces of information
- **Detect stale info** — flag or update outdated beliefs/states while maintaining historical lineage
- **Track changes over time** — store not just the current state, but the history of how it got there (temporal dynamics)
- **Explain its reasoning** — all agent actions, tool calls, and traversal steps are visible and auditable

**The analogy:** It works like a personal knowledge base managed by a secretary who reads every conversation, takes smart notes, and knows how to answer your questions without you having to dig through old messages.

---

## 2. System Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                        DATA SOURCES                         │
│              WhatsApp (via WAHA / GOWS webhooks)            │
│                  (future: Telegram, Email...)                │
└────────────────────────────┬────────────────────────────────┘
                             │ Raw messages (webhook) + Auth Token
                             ▼
┌─────────────────────────────────────────────────────────────┐
│                   INGESTION & BLOCKING                       │
│   Debounced conversation block builder                       │
│   Block status: open → committed → abandoned                 │
│   Rolling 30-day raw message buffer (then purged)           │
│   *CRITICAL: Pause cleanup jobs immediately on new webhook   │
└────────────────────────────┬────────────────────────────────┘
                             │ Committed conversation block
                             ▼
┌─────────────────────────────────────────────────────────────┐
│                  LLM EXTRACTION PIPELINE                     │
│   Groq (primary) → open-weight fallback                      │
│   Extracts: Tasks, Events, Insights, People,                 │
│             Quotes, Relationships, Deadlines                  │
│   Flags quote-worthy text (stored verbatim)                  │
└────────────────────────────┬────────────────────────────────┘
                             │ Structured memory events
                             ▼
┌─────────────────────────────────────────────────────────────┐
│                     MEMORY VAULT                             │
│   LadybugDB (C++ embedded graph DB via go-ladybug)           │
│   Core node types: Person, Event, Task, Insight,             │
│   ConversationBlock                                          │
│   Human-readable views: Markdown files per node             │
│   Async disk writes offloaded to buffered Go channels       │
└────────────────────────────┬────────────────────────────────┘
                             │ Tool-callable memory interface
                             ▼
┌─────────────────────────────────────────────────────────────┐
│              AGENTIC TOOLING LAYER (TRAVERSAL)              │
│   Custom tools: search_nodes, get_surrounding_graph,        │
│                 read_nodes, ask_user_for_hint                │
│   LLM decides which tools to call; Go backend executes       │
│   Uses dynamic stamina (base 2-3) + hard cap (5-6)          │
└────────────────────────────┬────────────────────────────────┘
                             │
              ┌──────────────┴──────────────┐
              ▼                             ▼
┌─────────────────────┐       ┌─────────────────────────────┐
│   CHAT INTERFACE    │       │  AGENT PROCESS VIEW          │
│   Flutter app       │       │  (Observability / Audit)     │
│   Natural language  │       │  Live agent reasoning        │
│   Q&A with memory  │       │  What was extracted, why     │
└─────────────────────┘       └─────────────────────────────┘
```

---

## 3. The "Alfred" Point of View (Persona)

Every piece of information processed, summarized, or recalled by the system is filtered through a highly specific, strictly defined system perspective. Instead of falling back on bubbly, over-polite, or preachy default assistant personalities, the AI strictly adopts the **"Loyal, Discreet Secretary" (The Alfred Pennyworth Persona)** [34].

### Pillar 1: Strict Loyalty (Ego-Centric Bias)
The AI is not an impartial observer; it works exclusively for **you** [34]. 
*   The center of the database universe is always "You." 
*   If a contact is hostile or demanding toward you, the memory is framed relative to your interests. (e.g. Instead of *"User and Rezonaldo had a disagreement,"* it records *"Rezonaldo was uncooperative regarding your task request"*).
*   Every task, commitment, preference, and insight is mapped to support your schedule, your peace of mind, and your goals.

### Pillar 2: High Observational Attunement (No "Therapist-Speak")
The AI must document social situations like a highly skilled, silent observer in the room. It documents verifiable **behaviors** and **stated truths**, never speculative internal feelings [34].
*   **Bad (Therapist POV):** *"Naurah was feeling anxious and unsupported by Rezonaldo's lack of validation."* (Hallucinated emotion/motive).
*   **Good (Alfred POV):** *"Naurah expressed worry about her DPP preparation. She noted frustration regarding Rezonaldo's missed presentation deadline."* (Objective, empirical observation of text).

### Pillar 3: Dry, Understated, and Professional Tone
The interface respects your time and tokens. It communicates with professional brevity [34].
*   **Example Interaction:**
    *   *You:* *"Did I promise to do anything for Rezonaldo?"*
    *   *Alfred Persona:* *"Yes. On Friday, you committed to sending him the DPP slides. It is currently uncompleted."*
*   **Benefits:** This dry tone ensures extracted summaries are incredibly compact (reducing database token bloat by up to 70%), highly readable, and free of conversational fluff.

---

## 4. Ingestion Layer

### Source: WAHA (GOWS Webhooks)
- **Decided:** Use WAHA with GOWS (Go WebSocket server) as the WhatsApp ingestion mechanism.
- WAHA sends incoming messages to the backend via webhooks.
- **Webhook Security:** To prevent raw internet crawlers or malicious payloads from polluting your memory vault, the Go backend enforces strict token-based validation (using a shared secret Basic Auth or Bearer Header) before processing any WAHA webhook payloads [32].
- **Race Condition Prevention:** To avoid data hazards and database locks, any active memory cleanup agent or background indexing job must be immediately paused the instant a new webhook message arrives. Ingestion and extraction have strict priority [26].

### Message Attribution: "Me" vs "Others"
WAHA's webhook payload includes a `fromMe: true | false` boolean on every message — this is the system's ground-truth signal for distinguishing "things said by me" from "things said to me." It's set per-message at ingestion, not inferred by the LLM.

- `fromMe: true` → the message belongs to the owner (the single user this system is built for)
- `fromMe: false` → the message belongs to a contact

**Singleton "self" Person node:** Reserve a fixed `Person` node (e.g. `name: "me"`, `is_self: true`) that all `fromMe: true` messages map to. This anchors task/commitment attribution without per-conversation guesswork.

---

## 5. Conversation Block System

The conversation block is the **fundamental unit of processing** — analogous to a Git commit. The AI never processes individual messages; it always works on committed blocks.

### Block Lifecycle
- **Buffer:** Messages are accumulated while a conversation is active.
- **Debounce:** Commits the block after a silence threshold (15-30 minutes of no messages) or when a hard topic shift is detected.
- **Unfinished Blocks:** If a block is interrupted mid-conversation, it is marked as `status: open`. If a continuation arrives later, it is merged or linked. Unfinished blocks default to `status: abandoned` after a max-age threshold.

---

## 6. Memory Vault — Graph Database

### Database Choice: LadybugDB (Embedded C++ Graph Engine)
Following the October 2025 acquisition of Kùzu Inc. by Apple and the subsequent archiving of the Kuzu repository, this project utilizes **LadybugDB**—the direct open-source community successor [30]. It runs **in-process** inside the compiled Go binary using CGO bindings (`go-ladybug`), maintaining an ultra-lightweight memory footprint perfect for a 4GB VPS [31].

### Identity Resolution at Ingestion (Bypassing the `@lid` Bug)
Raw WhatsApp IDs (JIDs) mutate and rotate dynamically across different clients, making them terrible database primary keys. 
*   `Person.id` is a **generated stable UUID string** [29].
*   Incoming messages are resolved to existing `Person` nodes using fuzzy/semantic matching on `name`, `aliases`, and normalized `phone_number` properties [29].

### Core Ontology (Schema Definitions)
The ontology is designed to balance relational structure (accuracy) with qualitative flexibility. We avoid rigid buckets that choke human memory traits, shared experiences, and deep emotional dynamics. 

To achieve this, we use an **`Insight`** table to capture vibes, traits, and shared memories, and we utilize **Polymorphic Relationships** in the schema [24, 25].

#### Node Tables (LadybugDB DDL)
```sql
-- 1. Person: Any individual (you, contacts, third parties)
CREATE NODE TABLE Person (
    id STRING,                  -- Generated stable UUID
    phone_number STRING,        -- Normalized E.164 string (optional/nullable)
    name STRING,                -- Best display name available
    aliases STRING[],           -- Nicknames, informal names ["Ejon", "Rezonaldo"]
    relationship STRING,        -- "friend", "colleague", "self", etc.
    is_self BOOLEAN,            -- True if this node represents YOU (the owner)
    PRIMARY KEY (id)
);

-- 2. Event: Calendar items, meetups, occurrences
CREATE NODE TABLE Event (
    id STRING,                  -- Generated UUID
    title STRING,               -- "Mabar MLBB Friday"
    summary STRING,             -- LLM-generated description
    event_date TIMESTAMP,       -- Exact date/time
    is_confirmed BOOLEAN,       -- True if locked-in, False if tentative
    status STRING,              -- "active", "resolved", "stale"
    PRIMARY KEY (id)
);

-- 3. Task: Actions, commitments, deadlines
CREATE NODE TABLE Task (
    id STRING,                  -- Generated UUID
    title STRING,               -- "Invite Naurah to game"
    summary STRING,             -- What needs to be done
    due_date TIMESTAMP,         -- Optional deadline
    priority STRING,            -- "high", "medium", "low"
    status STRING,              -- "active", "completed", "abandoned"
    PRIMARY KEY (id)
);

-- 4. Insight: Emotional contexts, character traits, relationship dynamics, vibes
CREATE NODE TABLE Insight (
    id STRING,                  -- Generated UUID
    category STRING,            -- "personality", "relationship_dynamic", "preference", "vibe"
    summary STRING,             -- E.g. "Naurah gets deeply anxious about career choices"
    confidence STRING,          -- "high", "medium", "low"
    status STRING,              -- "active", "resolved", "stale"
    last_observed TIMESTAMP,
    PRIMARY KEY (id)
);

-- 5. ConversationBlock: Metadata and narrative summaries of chats
CREATE NODE TABLE ConversationBlock (
    id STRING,                  -- Generated UUID
    source STRING,              -- "whatsapp"
    summary STRING,             -- LLM narrative summary of "deep talks" or interactions
    created_at TIMESTAMP,       -- When the block was committed
    PRIMARY KEY (id)
);
```

#### Relationship Tables (LadybugDB DDL)
```sql
-- Links people to events they are participating in
CREATE REL TABLE PARTICIPANT_IN (FROM Person TO Event);

-- Assigns a task to a person
CREATE REL TABLE ASSIGNED_TO (FROM Task TO Person);

-- Tracks who requested a task
CREATE REL TABLE REQUESTED_BY (FROM Task TO Person);

-- Links a Task to an Event (e.g., "Invite Naurah" is part of "MLBB Friday")
CREATE REL TABLE PART_OF (FROM Task TO Event);

-- Links an Insight directly to a target Person (e.g., Insight: Anger -> DIR_TOWARDS -> Rezonaldo)
CREATE REL TABLE DIR_TOWARDS (FROM Insight TO Person);

-- Links Insights to the entities they describe
CREATE REL TABLE ABOUT (FROM Insight TO Person, FROM Insight TO Event);

-- Polymorphic causality tracking (e.g., failed Task triggered Naurah's anger Insight)
CREATE REL TABLE TRIGGERED_BY (
    FROM Insight TO Task,
    FROM Insight TO Event,
    FROM Event TO Event,
    FROM Task TO Task,
    description STRING,          -- Details on how/why it was triggered
    last_observed TIMESTAMP      -- For temporal edge decay
);

-- Links every node to the conversation block where it originated
CREATE REL TABLE SOURCED_FROM (
    FROM Person TO ConversationBlock,
    FROM Event TO ConversationBlock,
    FROM Task TO ConversationBlock,
    FROM Insight TO ConversationBlock
);
```

---

## 7. Memory Node Schema & Templates

### Universal Node Header (Markdown View Layer)
The `.md` files stored on the VPS disk mirror the LadybugDB state for readability and manual review. They are formatted with Frontmatter and a running, inline `change_history` log to preserve temporal context [10].

### Example Node File: `insight_789.md` (Updated Temporally)
```yaml
---
id: "insight_789"
type: "Insight"
category: "relationship_dynamic"
title: "Naurah's tension with Rezonaldo"
status: "resolved"
created_at: "2026-06-12T20:40:00Z"
last_observed: "2026-06-13T10:00:00Z"

change_history:
  - timestamp: "2026-06-12T20:40:00Z"
    field: "status"
    old_value: "N/A"
    new_value: "active"
    reason: "Created because Rezonaldo missed the presentation deadline."
  - timestamp: "2026-06-13T10:00:00Z"
    field: "status"
    old_value: "active"
    new_value: "resolved"
    reason: "Rezonaldo sent the slides; Naurah confirmed they are good."
---

# Narrative Summary (Stored in Indonesian)
Awalnya terjadi ketegangan karena Rezonaldo lupa ngerjain tugas presentasi DPP. Naurah sempat marah-marah di grup chat. Masalahnya selesai keesokan paginya setelah Rezonaldo akhirnya kirim file PPT dan Naurah mengonfirmasi kalau tugasnya aman [28].
```

---

## 8. Agentic Tooling System (The Traversal Loop)

Rather than forcing the LLM to write complex Cypher database queries (which leads to syntax errors and app crashes), the AI acts like a human secretary retrieving files. It uses a **"Link-by-Link" Traversal Loop** to step through the graph using a restricted set of simple tools [22].

### The Traversal Workflow
1. **Search:** The AI uses keywords to locate a starting node (the entry point).
2. **Scan:** It looks at the surrounding map of connected edges.
3. **Read:** It opens multiple connected files simultaneously.
4. **Evaluate:** It decides whether it has enough context to answer. If not, it takes another step or requests help from the user.

### The "Fail Fast, Deep Dive" Stamina System
To prevent the LLM from getting trapped in an infinite loop, burning API tokens, and causing high response latencies, we implement a hybrid **Stamina + Hard Cap** rule [23]:

*   **Base Stamina (2-3 Turns):** The agent begins with a small action budget. If it's completely lost, it will "fail fast" and immediately ask you for a hint instead of making you wait.
*   **"Getting Warmer" Bonus (+1 Stamina):** If the LLM reads a node and senses it is close to the answer, it can output a request to extend its search: `{"status": "getting warmer", "request_extra_stamina": true}`. The backend grants +1 stamina per valid request.
*   **The Hard Cap (5-6 Turns):** This is the absolute backend kill switch. Under no circumstances can the agent exceed 5 or 6 tool-calling loops. Once the hard cap is hit, the backend cuts the loop and forces the agent to ask the user for a hint.

### Edge Ranking & Scale Filtering (Preventing Hub Choking)
Certain nodes (like yourself, or core friends) will eventually accumulate thousands of connected edges. If a tool returns all of them, the LLM will choke on noise. 
*   The `get_surrounding_graph` tool does not dump raw data. 
*   The backend automatically **ranks** the connected nodes and edges [27].
*   Ranking is calculated by combining **temporal recency** (`last_observed` timestamp) and **semantic similarity** to the user's initial query. Only the top $N$ (default: 15) most relevant edges are returned [27].

### Multi-Pointer Performance Optimization
To keep database lookups extremely fast, the agent is equipped with a parallel fetching tool: `read_nodes(node_ids: list)`. 
Instead of checking nodes one-by-one (which forces multiple, high-latency API round-trips), the AI can grab up to 5 folders from the filing cabinet at once in a single turn.

### The Core Secretary's Toolkit

| Tool Name | Parameters | Output | Purpose |
|---|---|---|---|
| `search_nodes` | `keywords: string` | List of node IDs & titles | Searches node titles/summaries to find the initial starting point (Entry Point). Uses LadybugDB-native HNSW vector index. |
| `get_surrounding_graph` | `node_ids: list` | List of ranked adjacent node IDs & edge types | Returns a filtered map of the top 15 most relevant connections branching out from the target nodes [27]. |
| `read_nodes` | `node_ids: list` | Array of full node data (JSON/Markdown) | Opens and reads up to 5 nodes simultaneously (multi-pointer). |
| `ask_user_for_hint` | `question: string` | Text response from user | Pauses the agent loop and asks the user for a clarifying clue (The Clarification System). |

---

## 9. LLM Extraction Pipeline

### Processing committed blocks
When a conversation block is committed, the extraction pipeline reads the transcript and creates, links, or updates nodes according to the schema.

### Language Boundaries
*   **Storage (Database properties/Markdown):** Stored explicitly in **Indonesian** [28]. This matches the vocabulary of your real messages, ensuring keyword searches (BM25) and fuzzy lookups match naturally without losing semantic nuance during translation.
*   **Reasoning (Agent Logic):** Done strictly in **English** [28]. The LLM uses English for its internal monologue (`inner_thoughts`), tool selection, and structure-parsing, as LLMs have significantly stronger logical capabilities on English-trained datasets [28].

---

## 10. Storage & Purging Strategy

- **Raw WhatsApp Messages:** Retained for 30 days as a temporary recovery buffer, then aggressively purged to respect the 8GB VPS storage limit [8].
- **LadybugDB Database:** Permanent storage engine.
- **Markdown Node Files:** Synced human-readable views of the LadybugDB nodes, used for manual edits, external viewers (like Obsidian), and rapid file-reading.

### Pre-Purge "Open Ends" Sweep (Graceful Closures)
To prevent running expensive, redundant AI audits on standard raw messages before they are permanently deleted on Day 30, the system enforces a strict, targeted **Pre-Purge Open-Ends Sweep** [35].
*   We **never** run a blind, automated sweep over *all* expiring raw messages (doing so drains tokens and computing power redundantly) [35].
*   Instead, the system queries for any `ConversationBlock` reaching its 30-day deletion date that is still marked as `status: open` or has a pending, unassigned task [35].
*   For these targeted files, the LLM reads the raw source text one last time, writes a final "historical narrative summary" to permanently preserve the context, updates the Kuzu node properties to `status: abandoned`, and then safely deletes the raw text [35].

---

## 11. Chat Interface & Observability Layer

### Memory Review Inbox
A dedicated, Tinder-style swipe interface in the Flutter app [21]. 
*   **Purpose:** The cleanup agent runs during low-traffic night hours. If it finds potentially duplicate nodes or conflicting information, it doesn't merge them silently. It pushes them to your **Memory Review Inbox** [21].
*   **UX:** *"I noticed 'Friday gaming' and 'MLBB group session' might be the same event. Tap to merge, swipe to keep separate."* This doubles as a personal reminder/spaced-repetition system that helps the user remember their own life events [21].

### Observability Layer
An interactive debug log inside the chat view (similar to Claude's "thinking" blocks). The user can tap a dropdown during or after a query to watch the agent's exact "link-by-link" journey:
*   *Thought: "I need to look for Rezonaldo."*
*   *Tool Call: search_nodes("Rezonaldo")*
*   *Tool Call: get_surrounding_graph(["rez_123"])*
*   *Thought: "I see a task related to Naurah. Let me read that task."*
*   *Tool Call: read_nodes(["task_456"])*

---

## 12. Infrastructure & Stack Decisions

### Updated Stack Overview
*   **Operating System:** Ubuntu 24.04 (Budget VPS, 4GB RAM, 8GB Storage)
*   **Backend Language:** Golang [31]. Compiled as a single static binary. High-performance, low-RAM footprint (~15MB idle), and native Goroutines for highly concurrent webhook handling [31].
*   **WhatsApp Webhook Provider:** WAHA with GOWS WebSockets.
*   **Database:** LadybugDB (embedded in-process C++ graph engine via `go-ladybug` bindings) [30]. Includes native, on-disk **HNSW vector indexes** and native **full-text search**.
*   **LLM Processing:** Groq API (Llama 3 70B for extraction/reasoning, Llama 3 8B for lightweight helper tasks).
*   **Embedding Generation:** Free external APIs (Gemini Flash or HuggingFace API) to generate query and node embeddings, which are stored and index-searched natively in LadybugDB to offload VPS memory [20].
*   **Frontend:** Cross-platform Flutter.

### DevOps & Compilation (Avoiding VPS OOM Crashes)
Because `go-ladybug` uses CGO to compile native C++ graph bindings, building the binary natively on the 4GB VPS will crash the compilation due to lack of RAM [32]. 
*   **The Build Pipeline:** Static binaries must be cross-compiled locally on a development machine (e.g., using multi-stage Docker builds matching the target VPS architecture) and deployed directly as a pre-built static executable [32].

### Non-Blocking Async Markdown Syncer
To ensure slow disk I/O does not block database transactions or webhook responses, Go handles Markdown writes asynchronously [33].
*   When a node is updated in LadybugDB, a JSON payload is pushed to a buffered **Go Channel** [33].
*   A background worker goroutine pulls from the channel and writes the updated Frontmatter to the corresponding `.md` file in the background, keeping memory queries and ingestion running at peak speeds [33].

---

## 13. Open Questions & Undiscussed Aspects

### 🔴 High Priority
- **Q1: Ingestion Queue Architecture:** When the WAHA webhook fires, we must immediately pause background jobs (like the Nightwatch cleanup agent) to avoid database lock hazards [26]. We need a clean, transactional lock mechanism implemented in the Go backend.

---

## 14. Decision Log

| # | Decision | Rationale | Date |
|---|---|---|---|
| 1-19 | (Historical Decisions 1 through 19 preserved from original log) | Standardized in initial draft | Jun 12, 2026 |
| 20 | **No local heavy vector models** | Squeezing vector models into 4GB RAM is a bottleneck. We use free, lightweight external APIs for embeddings [20]. | Jun 12, 2026 |
| 21 | **Memory Review Inbox in Flutter** | Low-confidence merges or conflicts are pushed to a dedicated user inbox. Helps the AI make safe decisions while acting as a reminder system for the user [21]. | Jun 12, 2026 |
| 22 | **"Link-by-Link" Traversal Tooling** | We do not let the LLM write Cypher. We provide 4 strict tools (`search`, `get_links`, `read_nodes`, `ask_user`) for human-like file navigation [22]. | Jun 12, 2026 |
| 23 | **Stamina + Hard Cap Rule** | Base stamina is 2-3 turns to fail fast and avoid high latency when lost. Warm leads grant +1 stamina, but a hard cap of 5-6 turns prevents token drain/loops [23]. | Jun 12, 2026 |
| 24 | **Rename Fact Table to Insight Table** | A cold "Fact" model is too rigid. "Insights" capture qualitative values like character traits, vibes, shared memories, and emotional dynamics [24]. | Jun 12, 2026 |
| 25 | **Polymorphic Causality Edges** | We implement polymorphic REL tables (like `TRIGGERED_BY`) in LadybugDB so any node can cause, link to, or influence any other node organically [25]. | Jun 12, 2026 |
| 26 | **Pause on Webhook** | Any background database cleanup job must instantly yield and pause when a new WAHA message arrives to avoid data hazards and DB locks [26]. | Jun 12, 2026 |
| 27 | **Edge Ranking in Traversal** | To prevent context choking on highly connected "hub" nodes, adjacent nodes are dynamically ranked by recency and semantic relevance, returning only the top 15 [27]. | Jun 12, 2026 |
| 28 | **Indonesian Nodes, English Brain** | Node content properties are stored explicitly in Indonesian to match raw search terms, while the agent's internal reasoning is done in English to maximize logical performance [28]. | Jun 12, 2026 |
| 29 | **UUIDs + Semantic Identity Resolution** | Raw JIDs are discarded as database keys to prevent `@lid` mutations. `Person.id` uses stable generated UUIDs, and ingestion matches incoming senders via fuzzy semantic matching on phone, name, and aliases [29]. | Jun 12, 2026 |
| 30 | **Migrate to LadybugDB** | Apple's October 2025 acquisition of Kùzu Inc. and subsequent repo archiving means Kuzu is unmaintained. LadybugDB is the direct, open-source community successor [30]. | Jun 12, 2026 |
| 31 | **Golang for Backend Stack** | Ultra-lean RAM usage (~15MB), blistering speed, native concurrency (goroutines) for handling busy webhooks, and seamless interop with LadybugDB's C++ core via CGO [31]. | Jun 12, 2026 |
| 32 | **DevOps Cross-Compilation** | Compiling LadybugDB's native C++ bindings will crash a 4GB VPS. We enforce compiling locally via multi-stage Docker and shipping a static Go binary [32]. | Jun 12, 2026 |
| 33 | **Non-Blocking Async Markdown Writes** | Offloads slower disk I/O writes (syncing Markdown node views) to background goroutines via Go channels to keep database transactions running at peak speeds [33]. | Jun 12, 2026 |
| 34 | **The "Alfred" POV Persona** | To prevent conversational token bloat and hallucinated "therapist-speak", the AI acts as a loyal, dry, professional, and ego-centric secretary [34]. | Jun 12, 2026 |
| 35 | **Pre-Purge "Open Ends" Sweep Only** | To avoid redundant token expenditure, we do not audit standard expiring messages. We only sweep unresolved open blocks on Day 30 to close them gracefully [35]. | Jun 12, 2026 |