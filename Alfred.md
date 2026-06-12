# Alfred — Agentic Secretary AI
### Project Specification

> **Status:** Planning / Pre-development
> **Last updated:** June 13, 2026
> **Stack:** Golang Backend · LadybugDB · SQLite · Groq API · WAHA (GOWS) · PWA Frontend · Ubuntu 24.04 VPS (4GB RAM / 8GB Storage)

---

## Table of Contents

1. [Vision & Scope](#1-vision--scope)
2. [System Overview](#2-system-overview)
3. [The Alfred Persona](#3-the-alfred-persona)
4. [Ingestion Layer](#4-ingestion-layer)
5. [Conversation Block System](#5-conversation-block-system)
6. [LLM Extraction Pipeline](#6-llm-extraction-pipeline)
7. [Memory Vault](#7-memory-vault)
8. [Agentic Query System](#8-agentic-query-system)
9. [Reminder System](#9-reminder-system)
10. [Background Agents](#10-background-agents)
11. [PWA Interface](#11-pwa-interface)
12. [Infrastructure](#12-infrastructure)
13. [Open Questions](#13-open-questions)
14. [Decision Log](#14-decision-log)

---

## 1. Vision & Scope

Alfred is an **agentic secretary AI** that passively watches your conversations — starting with WhatsApp — and turns raw chat into a **living, temporal memory system**. It is not a chat archive. It is a structured knowledge graph that grows over time and can be queried naturally.

### What Alfred Does
- **Remember** — extract and store structured facts, tasks, events, preferences, people, experiences, and social insights
- **Forget intentionally** — raw chat is ephemeral; only curated memory is permanent
- **Summarize** — compress conversation blocks into semantic summaries
- **Merge duplicates** — resolve conflicts between overlapping pieces of information
- **Detect stale info** — flag or update outdated beliefs/states while maintaining historical lineage
- **Track changes over time** — store not just the current state, but the history of how it got there
- **Remind proactively** — surface upcoming deadlines and obligations without being asked
- **Explain its reasoning** — all agent actions, tool calls, and traversal steps are visible and auditable

### Analogy
A personal knowledge base managed by a secretary who reads every conversation, takes smart notes, and knows how to answer your questions without you having to dig through old messages.

### Current Scope
- **Single user** (multi-user architecture planned post-validation)
- **Two WhatsApp group chats** as the initial data sources
- **Personal VPS** deployment, zero cloud spend

---

## 2. System Overview

```
┌─────────────────────────────────────────────────────────────┐
│                      DATA SOURCES                           │
│           WhatsApp (via WAHA / GOWS webhooks)               │
│               (future: Telegram, Email...)                  │
└────────────────────────────┬────────────────────────────────┘
                             │
            Raw messages (webhook) + Auth Token
                             │ 
                             ▼
┌─────────────────────────────────────────────────────────────┐
│                   INGESTION & BLOCKING                      │
│   Debounced conversation block builder                      │
│   Block status: open → committed → abandoned                │
│   Rolling 30-day raw message buffer (then purged)           │
│   *CRITICAL: Pause cleanup jobs immediately on new webhook  │
└────────────────────────────┬────────────────────────────────┘
                             │
                Committed conversation block
                             │ 
                             ▼
┌─────────────────────────────────────────────────────────────┐
│                  LLM EXTRACTION PIPELINE                    │
│   Groq (primary) → open-weight fallback chain               │
│   Extracts: Tasks, Events, Insights, People,                │
│             Quotes, Relationships, Deadlines                │
│   Flags quote-worthy text (stored verbatim)                 │
└────────────────────────────┬────────────────────────────────┘
                             │
                  Structured memory events
                             │ 
                             ▼
┌─────────────────────────────────────────────────────────────┐
│                     MEMORY VAULT                            │
│   LadybugDB (C++ embedded graph DB via go-ladybug)          │
│   Core node types: Person, Event, Task, Insight,            │
│   ConversationBlock                                         │
│   Human-readable views: Markdown files per node             │
│   Async disk writes offloaded to buffered Go channels       │
└────────────────────────────┬────────────────────────────────┘
                             │
               Tool-callable memory interface
                             │ 
                             ▼
┌─────────────────────────────────────────────────────────────┐
│              AGENTIC TOOLING LAYER (TRAVERSAL)              │
│   Custom tools: search_nodes, get_surrounding_graph,        │
│                 read_nodes, ask_user_for_hint,              │
│                 check_reminders, upsert_reminder            │
│   LLM decides which tools to call; Go backend executes      │
│   Uses dynamic stamina (base 2-3) + hard cap (5-6)          │
└────────────────────────────┬────────────────────────────────┘
                             │
              ┌──────────────┴──────────────┐
              ▼                             ▼
 ┌──────────────────────┐       ┌────────────────────────────┐
 │   CHAT INTERFACE     │       │  AGENT PROCESS VIEW        │
 │   PWA                │       │  (Observability / Audit)   │
 │   Natural language   │       │  Live agent reasoning      │
 │   Q&A with memory    │       │  What was extracted, why   │
 └──────────────────────┘       └────────────────────────────┘
```

### Data Flow Summary
1. WhatsApp messages arrive via WAHA webhook → buffered into open conversation blocks
2. After a silence window, the block is committed → sent to the LLM extraction pipeline
3. The LLM extracts structured nodes (people, tasks, events, insights) → written to LadybugDB
4. The agent also writes reminder rows to SQLite for any deadline-bearing tasks
5. The user queries Alfred via the PWA chat → the agent traverses the graph using tools
6. A dumb cron job scans SQLite and fires push notifications for upcoming reminders
7. Nightwatch (background agent) runs nightly to merge duplicates and flag stale data

---

## 3. The Alfred Persona

Every piece of information processed, summarized, or recalled by the system is filtered through a strictly defined perspective: the **"Loyal, Discreet Secretary" (The Alfred Pennyworth Persona)**.

Instead of falling back on bubbly, over-polite, or preachy default assistant personalities, Alfred operates on three pillars:

### Pillar 1: Strict Loyalty (Ego-Centric Bias)
Alfred works exclusively for **you**. The center of the database universe is always "You."
- If a contact is hostile or demanding toward you, the memory is framed relative to your interests. Instead of *"User and Bahlil had a disagreement,"* it records *"Bahlil was uncooperative regarding your task request."*
- Every task, commitment, preference, and insight is mapped to support your schedule, your peace of mind, and your goals.

### Pillar 2: High Observational Attunement (No "Therapist-Speak")
Alfred documents social situations like a skilled, silent observer. It documents verifiable **behaviors** and **stated truths**, never speculative internal feelings.
- **Bad (Therapist POV):** *"Bunga was feeling anxious and unsupported by Bahlil's lack of validation."* (Hallucinated emotion/motive)
- **Good (Alfred POV):** *"Bunga expressed worry about her DPP preparation. She noted frustration regarding Bahlil's missed presentation deadline."* (Objective, empirical observation of text)

### Pillar 3: Dry, Understated, and Professional Tone
Alfred communicates with professional brevity. This dry tone ensures extracted summaries are compact (reducing database token bloat by up to 70%), highly readable, and free of conversational fluff.
- *You: "Did I promise to do anything for Bahlil?"*
- *Alfred: "Yes. On Friday, you committed to sending him the DPP slides. It is currently uncompleted."*

---

## 4. Ingestion Layer

### Source: WAHA (GOWS Webhooks)
WhatsApp messages are ingested via WAHA with GOWS (Go WebSocket server). WAHA sends incoming messages to the Go backend via webhooks.

**Webhook Security:** The Go backend enforces strict token-based validation (Bearer Header) before processing any WAHA webhook payload, preventing raw internet crawlers or malicious payloads from polluting the memory vault.

**Race Condition Prevention:** Any active background job (Nightwatch cleanup, indexing) must be immediately paused the instant a new webhook message arrives. Ingestion and extraction have strict priority over all background work. See Q1.

### Message Attribution: "Me" vs "Others"
WAHA's webhook payload includes a `fromMe: true | false` boolean on every message — this is the ground-truth signal for distinguishing messages sent by the owner vs messages sent by contacts. It is set at ingestion, not inferred by the LLM.

A fixed **singleton "self" Person node** (`name: "me"`, `is_self: true`) anchors all `fromMe: true` messages. This ensures task and commitment attribution is unambiguous without per-conversation guesswork.

---

## 5. Conversation Block System

The conversation block is the **fundamental unit of processing** — analogous to a Git commit. Alfred never processes individual messages in isolation; it always works on committed blocks.

### Block Lifecycle
- **Buffer:** Messages accumulate while a conversation is active.
- **Debounce:** The block commits after a silence threshold (15-30 minutes of no messages) or when a hard topic shift is detected.
- **Open:** A block that has been started but not yet committed.
- **Committed:** A fully sealed block ready for LLM extraction.
- **Abandoned:** A block that was interrupted and never naturally closed. Defaulted after a max-age threshold. On Day 30, any open block reaching the purge window gets a final LLM sweep before raw deletion (see Section 10).

---

## 6. LLM Extraction Pipeline

### Trigger
When a conversation block is committed, the extraction pipeline reads the transcript and creates, links, or updates nodes in LadybugDB according to the schema.

### Language Boundaries
- **Storage (Database properties / Markdown):** Stored explicitly in **Indonesian**. This matches the vocabulary of real messages, ensuring keyword searches (BM25) and fuzzy lookups match naturally without losing semantic nuance in translation.
- **Reasoning (Agent Logic):** Done strictly in **English**. The LLM uses English for its internal monologue (`inner_thoughts`), tool selection, and structure-parsing, as LLMs have significantly stronger logical capabilities on English-trained data.

### LLM Fallback Chain
Groq is the primary inference provider. If a call fails, the system walks down a prioritized list of fallback model APIs in a try/catch loop until one succeeds.

### Prompt Caching
Alfred's system prompt is long and resent on every API call — including every turn of the agentic traversal loop (up to 5-6 turns per query). Groq's prompt caching halves input token costs on repeated prompt prefixes and does not count cached tokens toward rate limits. Prompt caching must be enabled on all extraction and traversal calls.

---

## 7. Memory Vault

### Database: LadybugDB
Following the October 2025 acquisition of Kùzu Inc. by Apple and the subsequent archiving of the Kuzu repository, this project uses **LadybugDB** — the direct open-source community successor. It runs **in-process** inside the compiled Go binary via CGO bindings (`go-ladybug`), maintaining an ultra-lightweight memory footprint suitable for a 4GB VPS.

### Identity Resolution (Bypassing the `@lid` Bug)
Raw WhatsApp JIDs mutate and rotate dynamically across different clients, making them unreliable as database keys.
- `Person.id` is a **generated stable UUID**.
- Incoming messages are resolved to existing `Person` nodes using fuzzy/semantic matching on `name`, `aliases`, and normalized `phone_number`.

### Node Schema (LadybugDB DDL)

```sql
-- 1. Person: Any individual (you, contacts, third parties)
CREATE NODE TABLE Person (
    id STRING,                  -- Generated stable UUID
    phone_number STRING,        -- Normalized E.164 string (nullable)
    name STRING,                -- Best display name available
    aliases STRING[],           -- Nicknames, informal names ["Ejon", "Bahlil"]
    relationship STRING,        -- "friend", "colleague", "self", etc.
    is_self BOOLEAN,            -- True if this node represents YOU (the owner)
    PRIMARY KEY (id)
);

-- 2. Event: Calendar items, meetups, occurrences
CREATE NODE TABLE Event (
    id STRING,
    title STRING,               -- "Mabar MLBB Friday"
    summary STRING,             -- LLM-generated description
    event_date TIMESTAMP,
    is_confirmed BOOLEAN,       -- True if locked-in, False if tentative
    status STRING,              -- "active", "resolved", "stale"
    PRIMARY KEY (id)
);

-- 3. Task: Actions, commitments, deadlines
CREATE NODE TABLE Task (
    id STRING,
    title STRING,               -- "Invite Bunga to game"
    summary STRING,
    due_date TIMESTAMP,         -- Optional
    priority STRING,            -- "high", "medium", "low"
    status STRING,              -- "active", "completed", "abandoned"
    PRIMARY KEY (id)
);

-- 4. Insight: Emotional contexts, character traits, relationship dynamics, vibes
CREATE NODE TABLE Insight (
    id STRING,
    category STRING,            -- "personality", "relationship_dynamic", "preference", "vibe"
    summary STRING,             -- "Bunga gets deeply anxious about career choices"
    confidence STRING,          -- "high", "medium", "low"
    status STRING,              -- "active", "resolved", "stale"
    last_observed TIMESTAMP,
    PRIMARY KEY (id)
);

-- 5. ConversationBlock: Metadata and narrative summaries of chats
CREATE NODE TABLE ConversationBlock (
    id STRING,
    source STRING,              -- "whatsapp"
    summary STRING,             -- LLM narrative summary
    created_at TIMESTAMP,
    PRIMARY KEY (id)
);
```

### Relationship Schema (LadybugDB DDL)

```sql
CREATE REL TABLE PARTICIPANT_IN (FROM Person TO Event);
CREATE REL TABLE ASSIGNED_TO (FROM Task TO Person);
CREATE REL TABLE REQUESTED_BY (FROM Task TO Person);
CREATE REL TABLE PART_OF (FROM Task TO Event);
CREATE REL TABLE DIR_TOWARDS (FROM Insight TO Person);
CREATE REL TABLE ABOUT (FROM Insight TO Person, FROM Insight TO Event);

-- Polymorphic causality (e.g., failed Task triggered Bunga's anger Insight)
CREATE REL TABLE TRIGGERED_BY (
    FROM Insight TO Task,
    FROM Insight TO Event,
    FROM Event TO Event,
    FROM Task TO Task,
    description STRING,
    last_observed TIMESTAMP
);

-- Links every node back to its originating conversation block
CREATE REL TABLE SOURCED_FROM (
    FROM Person TO ConversationBlock,
    FROM Event TO ConversationBlock,
    FROM Task TO ConversationBlock,
    FROM Insight TO ConversationBlock
);
```

### Markdown Node Files (Human-Readable Layer)
Every node in LadybugDB has a corresponding `.md` file on disk — a human-readable mirror used for manual review, Obsidian viewing, and rapid inspection. Files are written asynchronously via a buffered Go Channel (see Section 12) so disk I/O never blocks database operations.

Each file uses Frontmatter and an inline `change_history` log to preserve temporal context:

```yaml
---
id: "insight_789"
type: "Insight"
category: "relationship_dynamic"
title: "Bunga's tension with Bahlil"
status: "resolved"
created_at: "2026-06-12T20:40:00Z"
last_observed: "2026-06-13T10:00:00Z"

change_history:
  - timestamp: "2026-06-12T20:40:00Z"
    field: "status"
    old_value: "N/A"
    new_value: "active"
    reason: "Created because Bahlil missed the presentation deadline."
  - timestamp: "2026-06-13T10:00:00Z"
    field: "status"
    old_value: "active"
    new_value: "resolved"
    reason: "Bahlil sent the slides; Bunga confirmed they are good."
---

# Narrative Summary (Stored in Indonesian)
Awalnya terjadi ketegangan karena Bahlil lupa ngerjain tugas presentasi DPP...
```

### Storage & Purging Strategy
- **Raw WhatsApp Messages:** Retained for 30 days as a temporary recovery buffer, then purged.
- **LadybugDB:** Permanent storage engine. Nodes are rarely deleted; bad data is the only deletion trigger.
- **Markdown Node Files:** Synced human-readable views; permanent alongside the graph.

#### Pre-Purge "Open Ends" Sweep
On Day 30, before raw messages are deleted, the system runs a targeted sweep — **never** a blind audit of all expiring messages. Only `ConversationBlock` nodes still marked `status: open` or with a pending unassigned task are targeted. For these, the LLM writes a final historical narrative summary, marks the block `status: abandoned`, then deletes the raw text.

---

## 8. Agentic Query System

Rather than writing raw database queries, the agent acts like a human secretary retrieving files — stepping through the graph using a restricted set of simple tools. This is the **"Link-by-Link" Traversal Loop**.

### Traversal Workflow
1. **Search** — locate a starting node using keywords (entry point)
2. **Scan** — look at the surrounding map of connected edges
3. **Read** — open multiple connected nodes simultaneously
4. **Evaluate** — decide if there's enough context to answer; if not, take another step or ask the user

### Stamina System (Fail Fast, Deep Dive)
To prevent infinite loops and token drain:
- **Base Stamina (2-3 turns):** Small action budget. If completely lost, agent fails fast and asks for a hint immediately.
- **"Getting Warmer" Bonus (+1 stamina):** If the agent senses it's close to the answer, it can request an extension: `{"status": "getting warmer", "request_extra_stamina": true}`. Backend grants +1 per valid request.
- **Hard Cap (5-6 turns):** Absolute backend kill switch. Once hit, the backend cuts the loop and forces the agent to ask the user for a hint.

### Edge Ranking (Preventing Hub Choking)
Highly connected nodes (like yourself or core friends) will eventually accumulate thousands of edges. `get_surrounding_graph` never dumps raw data — the backend ranks adjacent nodes by combining **temporal recency** (`last_observed`) and **semantic similarity** to the original query, returning only the top 15.

### Full Agent Toolkit

| Tool | Parameters | Output | Purpose |
|---|---|---|---|
| `search_nodes` | `keywords: string` | List of node IDs & titles | Entry point search using LadybugDB-native HNSW vector index |
| `get_surrounding_graph` | `node_ids: list` | Top 15 ranked adjacent node IDs & edge types | Traversal map from current position |
| `read_nodes` | `node_ids: list` | Array of full node data (JSON/Markdown) | Opens up to 5 nodes simultaneously |
| `ask_user_for_hint` | `question: string` | Text response from user | Pauses loop and requests a clarifying clue |
| `check_reminders` | `task_ref?: string` | List of existing reminder rows | Checks SQLite before inserting to prevent duplicates |
| `upsert_reminder` | `message, deadline, status, task_ref?` | Confirmation | Inserts or updates a reminder row in SQLite |

---

## 9. Reminder System

Alfred is a proactive secretary. The Reminder System enables Alfred to notify the user about upcoming deadlines and obligations without requiring an LLM call at notification time.

### Philosophy
Reminders are **operational/transient data**, not knowledge. Storing them as LadybugDB nodes would create overlap with Task nodes with no benefit. They live in a separate **SQLite file** (`reminders.db`) alongside the main binary. SQLite requires no server process, no configuration, and handles concurrent Go reads/writes safely.

### Schema

```sql
CREATE TABLE reminders (
    id          TEXT PRIMARY KEY,
    message     TEXT NOT NULL,          -- Human-readable reminder text (Indonesian)
    deadline    DATETIME NOT NULL,
    status      TEXT NOT NULL,          -- pending | sent | needs_clarification | dismissed
    task_ref    TEXT,                   -- Optional FK to a LadybugDB Task node ID
    created_at  DATETIME NOT NULL
);
-- Prevents duplicate reminders for the same task+deadline
CREATE UNIQUE INDEX ON reminders(task_ref, deadline);
```

### Who Populates It
The **main agent** is solely responsible for writing reminders — never Nightwatch. This happens in two flows:
- **Extraction pipeline** — when a committed block contains a user-owned task or deadline
- **User query flow** — when the agent surfaces a task during traversal and a reminder is warranted

Before inserting, the agent calls `check_reminders` to prevent duplicates. The `UNIQUE` index is a database-level safety net.

### Status Lifecycle

```
pending → sent                  (cron fires push notif)
pending → dismissed             (user dismisses from PWA)
pending → needs_clarification   (agent unsure what the task actually is)
needs_clarification → pending   (user clarifies, agent updates)
needs_clarification → dismissed
```

### Cron Job (Dumb Scanner — No LLM)
A Go cron job runs on a configurable interval (e.g. every hour):
1. Query SQLite for `status = 'pending'` rows where `deadline` is within the notification window
2. Fire push notification via PWA Push API
3. Update `status` to `sent`

The cron never touches LadybugDB and never calls Groq.

### Immediate Push on `needs_clarification`
Any reminder or node flagged `needs_clarification` is pushed to the user **immediately upon commit** — it does not wait for the cron. This rule applies system-wide: the cron handles scheduled pending reminders only; anything requiring user input is surfaced right away.

### Cascading Deletion
If a Task node is deleted from LadybugDB (bad data), its associated reminder rows cascade-delete via `task_ref`.

---

## 10. Background Agents

### Nightwatch (Database Maintenance Agent)
Nightwatch runs during low-traffic hours (nightly). Its **sole responsibility is database maintenance** — it does not write reminders and does not handle user-facing logic.

**Nightwatch responsibilities:**
- Detect potentially duplicate nodes and push them to the **Memory Review Inbox** (PWA) rather than merging silently. Example: *"I noticed 'Friday gaming' and 'MLBB group session' might be the same event. Tap to merge, swipe to keep separate."*
- Flag or update stale nodes (e.g. tasks that are long overdue with no activity)
- Run the Pre-Purge Open Ends Sweep on Day 30 (see Section 7)

Nightwatch must immediately yield and pause when a new WAHA webhook arrives.

### Reminder Cron
A separate, LLM-free cron job responsible only for scanning SQLite and dispatching push notifications. Described in full in Section 9.

---

## 11. PWA Interface

### Authentication
A JWT credential gate loads before the PWA renders anything. The user provides a username and password; the Go backend validates and returns a JWT which the PWA holds in memory and attaches to all subsequent API requests. No session persistence — re-login on refresh is acceptable. This is a security barrier only, not a user management system.

### Chat Interface
Natural language Q&A with Alfred's memory. The user types a question; the agent runs the traversal loop and responds in Alfred's persona.

### Observability Layer
An interactive debug log inside the chat view (similar to Claude's "thinking" blocks). The user can expand a dropdown during or after a query to watch the agent's exact step-by-step journey:
- *Thought: "I need to look for Bahlil."*
- *Tool Call: search_nodes("Bahlil")*
- *Tool Call: get_surrounding_graph(["rez_123"])*
- *Thought: "I see a task related to Bunga. Let me read that."*
- *Tool Call: read_nodes(["task_456"])*

### Memory Review Inbox
A Tinder-style swipe interface. When Nightwatch finds potentially duplicate or conflicting nodes, it pushes them here rather than merging silently. The user taps to merge or swipes to keep separate. This doubles as a spaced-repetition system for the user's own life events.

### Push Notifications
Delivered via the PWA Push API + Service Workers. Service Workers run in the background even when the browser is closed, enabling reliable delivery on Android. iOS is supported since Safari 16.4 (2023).

---

## 12. Infrastructure

### Stack

| Component | Choice | Reason |
|---|---|---|
| OS | Ubuntu 24.04 | Stable LTS on budget VPS |
| Backend | Golang | ~15MB idle RAM, native goroutines, single static binary, CGO interop for LadybugDB |
| Graph DB | LadybugDB (via `go-ladybug`) | In-process C++ engine, HNSW vector index, full-text search, Kuzu's open-source successor |
| Reminder DB | SQLite (`reminders.db`) | Single file, no server, concurrent-safe, native query support |
| LLM | Groq API (Llama 3 70B / 8B) | Free tier generous enough for personal scale, fast inference |
| Embeddings | Gemini Flash or HuggingFace API | Free external APIs to avoid loading vector models into 4GB RAM |
| WhatsApp | WAHA + GOWS WebSockets | Proven webhook provider |
| Frontend | PWA | Browser-native, no app store, works across all devices via URL, push notif support |

### DevOps & Compilation
`go-ladybug` uses CGO to compile native C++ bindings. Building natively on the 4GB VPS will crash due to insufficient RAM. The build pipeline is:
- Cross-compile locally using multi-stage Docker builds matching the target VPS architecture
- Deploy a pre-built static binary directly to the VPS

### Non-Blocking Async Markdown Writes
To ensure disk I/O never blocks database transactions or webhook responses:
- When a node is updated in LadybugDB, a JSON payload is pushed to a buffered **Go Channel**
- A background goroutine pulls from the channel and writes the updated Markdown file asynchronously

---

## 13. Open Questions

### ⚫ Critical Priority

**Q1: Missing Content Field on Node Types**
Most node types currently only have structured metadata fields with no free-form content field where the actual substance of the memory lives. Every node likely needs a `content STRING` or consistent `summary STRING` field that holds the narrative and context. The current schema is too restrictive without it.

**Q2: Memory Update & Modification Flow**
The spec describes how nodes get created but not how they get updated when new information contradicts or extends existing memory. Undecided: when the extraction pipeline finds information relating to an existing node, does it overwrite, append, or version? Who writes to `change_history` — the agent during extraction or the user manually? What is the exact PWA mechanism for user corrections (Q2)? What happens to both nodes' content when Nightwatch merges two duplicates via the Memory Review Inbox?

### 🔴 High Priority

**Q1: Ingestion Queue Architecture**
When the WAHA webhook fires, background jobs must immediately pause to avoid database lock hazards. A clean transactional lock mechanism needs to be designed in the Go backend — the exact implementation is undecided.

**Q2: Error Correction Flow**
How does the user correct Alfred when it extracts something wrong? The agent may misinterpret a message, create a wrong node, or link things incorrectly. A deliberate correction mechanism — likely via the PWA — needs to be designed so corrections feed back into the agent and update the graph cleanly without leaving stale data.

### 🟡 Low Priority / Future

**Q1: Off-Site Backup Strategy**
LadybugDB and SQLite are both single files on the VPS. If the VPS dies, all memory is lost. A scheduled backup strategy is needed — candidates: rsync to another machine, or push to object storage (Backblaze B2 or Cloudflare R2). Not important for prototype phase.

**Q2: Multi-Source Architecture**
How do future ingestion sources (Telegram, Email) plug in without rewriting the ingestion layer? Likely a source-agnostic message interface that WAHA and future adapters all conform to. Not a near-future priority.

---

## 14. Decision Log

| # | Decision | Rationale | Date |
|---|---|---|---|
| 20 | **No local vector models** | Squeezing vector models into 4GB RAM is a bottleneck. Free lightweight external APIs handle embeddings instead. | Jun 12, 2026 |
| 21 | **Memory Review Inbox** | Low-confidence merges or conflicts are pushed to a user inbox rather than merged silently. Acts as a spaced-repetition system for the user's own life events. | Jun 12, 2026 |
| 22 | **Link-by-Link Traversal Tooling** | No raw Cypher queries. The agent uses 4 strict tools for human-like graph navigation, avoiding syntax errors and hallucinated queries. | Jun 12, 2026 |
| 23 | **Stamina + Hard Cap Rule** | Base stamina 2-3 turns to fail fast. Warm leads grant +1. Hard cap at 5-6 prevents token drain and infinite loops. | Jun 12, 2026 |
| 24 | **Insight Table (not Fact Table)** | "Facts" are too rigid. Insights capture qualitative values like character traits, vibes, shared memories, and emotional dynamics. | Jun 12, 2026 |
| 25 | **Polymorphic Causality Edges** | Polymorphic `TRIGGERED_BY` REL tables allow any node to cause, link to, or influence any other node organically. | Jun 12, 2026 |
| 26 | **Pause on Webhook** | All background jobs must instantly yield when a new WAHA message arrives to avoid data hazards and DB locks. | Jun 12, 2026 |
| 27 | **Edge Ranking in Traversal** | Adjacent nodes ranked by recency and semantic relevance, top 15 only. Prevents context choking on highly connected hub nodes. | Jun 12, 2026 |
| 28 | **Indonesian Nodes, English Brain** | Node content stored in Indonesian to match real search terms. Agent reasoning done in English for stronger logical performance. | Jun 12, 2026 |
| 29 | **UUIDs + Semantic Identity Resolution** | Raw JIDs discarded as keys due to `@lid` mutation risk. Stable UUIDs used; incoming senders matched via fuzzy semantic matching. | Jun 12, 2026 |
| 30 | **Migrate to LadybugDB** | Apple's October 2025 acquisition of Kùzu Inc. and subsequent repo archiving makes Kuzu unmaintained. LadybugDB is the direct open-source successor. | Jun 12, 2026 |
| 31 | **Golang Backend** | ~15MB idle RAM, fast, native goroutines for concurrent webhook handling, seamless CGO interop with LadybugDB. | Jun 12, 2026 |
| 32 | **Cross-Compilation DevOps** | Compiling LadybugDB's C++ bindings natively on the VPS would OOM crash it. Compiled locally via multi-stage Docker, deployed as a static binary. | Jun 12, 2026 |
| 33 | **Non-Blocking Async Markdown Writes** | Disk I/O writes offloaded to background goroutines via Go channels. Keeps database transactions and ingestion at peak speed. | Jun 12, 2026 |
| 34 | **The Alfred Persona** | Prevents token bloat and hallucinated therapist-speak. Loyal, dry, professional, ego-centric secretary voice. | Jun 12, 2026 |
| 35 | **Pre-Purge Open Ends Sweep Only** | Never audit all expiring messages blindly. Only sweep unresolved open blocks on Day 30 to close them gracefully. | Jun 12, 2026 |
| 36 | **Flutter → PWA** | Alfred is currently single-user; Flutter's build pipeline adds complexity with no gain at this stage. PWA covers all required features (push notifications, chat, observability, swipe inbox) and opens across all devices via URL. iOS dev burden avoided. PWA served as a static build from the Go backend. Multi-user architecture planned post-validation. | Jun 13, 2026 |
| 37 | **JWT Credential Gate** | PWA is browser-accessible via URL — a security gate is required. JWT login screen before anything renders. No session persistence. Security barrier only, not a user management system. Chosen over Basic Auth for cleaner future multi-user extensibility. | Jun 13, 2026 |
| 38 | **SQLite for Reminder Storage** | Reminders are operational/transient data, not knowledge. Storing as LadybugDB nodes would overlap with Task nodes. SQLite is a single file, no server process, concurrent-safe. Chosen over JSON file for concurrency safety and native query support. | Jun 13, 2026 |
| 39 | **Reminders Owned by Main Agent, Not Nightwatch** | Nightwatch is database maintenance only. Reminder creation is a side effect of extraction and query flows — the agent has the context to decide what is reminder-worthy. | Jun 13, 2026 |
| 40 | **Dumb Cron for Reminder Dispatch** | The cron that fires push notifications is intentionally LLM-free. Pure deadline scanner: query SQLite, fire push notif, mark sent. LLM only involved upstream (writing reminders) and downstream if clarification needed. | Jun 13, 2026 |
| 41 | **Immediate Push on needs_clarification** | Anything flagged needs_clarification is pushed to the user immediately upon commit, not queued for the cron. Applies system-wide. | Jun 13, 2026 |
| 42 | **Prompt Caching** | System prompt is long and resent on every LLM call including every traversal loop turn. Groq prompt caching halves input token costs on repeated prefixes and cached tokens don't count toward rate limits. Must be enabled on all extraction and traversal calls. | Jun 13, 2026 |