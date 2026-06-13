# Alfred — Agentic Memory System
### Architecture & Design Specification

> **Status:** Planning / Pre-development
> **Last updated:** June 13, 2026 (revalidated)
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
13. [Project Phases](#13-project-phases)
14. [Open Questions](#14-open-questions)
15. [Decision Log](#15-decision-log)

---

## 1. Vision & Scope

### The Real Thesis

Most AI memory systems either ask you to manually log information, or they treat memory as a flat retrieval problem — embed everything, search by similarity, call it done. Alfred is neither.

**The claim Alfred is built to prove:** an agent can maintain a coherent, temporally-aware, self-correcting knowledge graph from noisy, informal, multilingual conversational data — without any user curation. Not just storing what you said, but knowing *when you said it*, *whether it's still true*, *how confident to be*, and *how things connect to each other*.

The secretary framing is the test case, not the product. WhatsApp group chats are the data source because real data is messy — informal language, implicit references, contradictions, half-finished thoughts, Indonesian slang mid-sentence. If the memory model holds up there, it holds up anywhere.

### The Concrete Problem

The specific failure mode Alfred is built to solve has two distinct forms:

**Missing** — a message comes in while you're busy, you skim it, you don't register that it contained an obligation directed at you. The information never entered your head. This is the harder case: Alfred must catch what you didn't.

**Forgetting** — you saw it, you knew you had to do something, it got buried, and it fell out of working memory before you acted. Alfred must surface it before the deadline passes.

Both failures are common in group chats where obligations are often implicit, informally phrased, or directed at you without your name being mentioned. Alfred's proactive push — surfacing `needs_clarification` flags and deadline reminders without being asked — is the most important feature in the system. The knowledge graph and traversal system exist to make those pushes accurate.

### Design Philosophy: Honest Gap Over Confident Wrong

Alfred is explicitly designed to keep the user in the loop rather than quietly automate their judgment away. When the extraction pipeline cannot attribute a task with sufficient confidence, it creates a partial node flagged `needs_clarification` and pushes it to the user immediately. It never guesses.

This is intentional. A real secretary surfaces ambiguous obligations and asks — they don't silently assign them to the wrong person. The volume of `needs_clarification` flags is not a failure metric; it is an accurate reflection of how ambiguous group chat obligations actually are. The user decides. Alfred catches and presents.

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
Not just a smart notes app. Alfred is closer to a personal epistemics engine — something that knows not just *what* you know, but *when you learned it*, *whether it's still true*, and *how confident to be*. The secretary is the interface. The knowledge graph is the point.

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
│   ConversationBlock, Circle                                 │
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

Blocks are strictly **per-chat**. `chat_id` on ConversationBlock is the enforced boundary. Messages from different group chats never share a block.

---

## 6. LLM Extraction Pipeline

### Trigger
When a conversation block is committed, the extraction pipeline reads the transcript and creates, links, or updates nodes in LadybugDB according to the schema.

### Cross-Block Event Deduplication
The same event (e.g. SoTQ) may be referenced across many blocks over time. Phase 2 first searches by alias match. If no alias matches but the block is temporally proximate to a known event and context is consistent, the agent attempts to link. If confidence is insufficient, it flags `needs_clarification` and asks the user rather than silently creating a duplicate node.

### Language Boundaries
- **Storage (Database properties):** Stored explicitly in **Indonesian**. This matches the vocabulary of real messages, ensuring keyword searches (BM25) and fuzzy lookups match naturally without losing semantic nuance in translation.
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

> **Content model:** Every node type except Person carries a `content STRING` field — Alfred's dry, first-person narrative of what this node represents, written in Indonesian. This is the primary field the linking agent reads in Phase 2. Structured metadata fields (status, dates, priority) exist alongside it for querying and Nightwatch logic. Person is a pure identity anchor — its substance lives in the surrounding graph.

> **Aliases:** Every node type carries `aliases STRING[]`. Aliases are informal names, abbreviations, or references this node is known by in chat ("rapat tadi", "sotq", "event kemarin"). The Phase 2 linking agent searches aliases, not just canonical names, to satisfy the explicit keyword rule.

> **History & Updates:** Nodes are modified in place. `content` always reflects current truth. When `content` is overwritten, the old value is prepended to `history STRING[]` as a human-readable entry (`"YYYY-MM-DD HH:MM - [narrative]"`), newest first. The new `content` must briefly acknowledge the prior state so the agent is hinted that something changed without needing to read `history`. `history` is only traversed when a query explicitly requires it.

> **needs_clarification:** Every node type carries `needs_clarification BOOLEAN`. When set true, the node was created with one or more fields that could not be grounded in source text or vault data. An immediate push notification is sent to the user. The node is visible in the Memory Review Inbox until resolved. Alfred never guesses to fill a field — a confident wrong node is worse than an honest gap.

```sql
-- 1. Person: Identity anchor. Substance lives in surrounding graph, not on the node itself.
CREATE NODE TABLE Person (
    id STRING,                      -- Generated stable UUID
    name STRING,                    -- Best display name available
    aliases STRING[],               -- Nicknames, informal names ["qil", "pit", "rapit"]
    phone_number STRING,            -- Normalized E.164 (nullable)
    is_self BOOLEAN,                -- True if this node represents the owner
    needs_clarification BOOLEAN,    -- True if identity could not be fully resolved
    PRIMARY KEY (id)
);

-- 2. Event: Any occurrence — meetings, sessions, social events, deadlines
CREATE NODE TABLE Event (
    id STRING,
    name STRING,                    -- Canonical name: "SoTQ IEEE 2026"
    aliases STRING[],               -- ["sotq", "acara tadi", "event kemarin"]
    content STRING,                 -- Alfred's narrative of current state, written in Indonesian. Must acknowledge prior state if overwriting.
    history STRING[],               -- Past content values, newest first. Format: "YYYY-MM-DD HH:MM - [narrative]"
    event_date TIMESTAMP,           -- Nullable if date unconfirmed
    status STRING,                  -- "planned" | "active" | "completed" | "cancelled" | "stale"
    created_at TIMESTAMP,
    needs_clarification BOOLEAN,
    PRIMARY KEY (id)
);

-- 3. Task: Actions, commitments, obligations with a named owner
CREATE NODE TABLE Task (
    id STRING,
    name STRING,                    -- Canonical short label
    aliases STRING[],               -- Alternative references in chat
    content STRING,                 -- Current state: what, who, why. Must acknowledge prior state if overwriting.
    history STRING[],               -- Past content values, newest first. Format: "YYYY-MM-DD HH:MM - [narrative]"
    verbatim STRING,                -- Exact source text if wording itself is meaningful (nullable)
    due_date TIMESTAMP,             -- Nullable
    priority STRING,                -- "high" | "medium" | "low"
    status STRING,                  -- "active" | "completed" | "abandoned" | "stale"
    created_at TIMESTAMP,
    needs_clarification BOOLEAN,
    PRIMARY KEY (id)
);

-- 4. Insight: Behavioral patterns, relationship dynamics, character observations
CREATE NODE TABLE Insight (
    id STRING,
    name STRING,                    -- Short label
    aliases STRING[],               -- Alternative phrasings observed in chat
    content STRING,                 -- Current observational narrative, grounded in source text. Must acknowledge prior state if overwriting.
    history STRING[],               -- Past content values, newest first. Format: "YYYY-MM-DD HH:MM - [narrative]"
    verbatim STRING,                -- Exact statement that triggered this insight (nullable)
    category STRING,                -- "personality" | "relationship_dynamic" | "preference" | "pattern"
    confidence STRING,              -- "high" | "medium" | "low"
    status STRING,                  -- "active" | "stale" | "contradicted"
    last_observed TIMESTAMP,
    created_at TIMESTAMP,
    needs_clarification BOOLEAN,
    PRIMARY KEY (id)
);

-- 5. ConversationBlock: The atomic unit of processing. Holds the raw transcript.
CREATE NODE TABLE ConversationBlock (
    id STRING,
    source STRING,                  -- "whatsapp"
    chat_id STRING,                 -- Which group/chat this block came from
    raw_transcript STRING,          -- Full message log including WAHA quote payloads, preserved as-is
    content STRING,                 -- Alfred's narrative summary, written after extraction
    status STRING,                  -- "open" | "committed" | "abandoned"
    created_at TIMESTAMP,
    PRIMARY KEY (id)
);

-- 6. Circle: A named group of people with shared context
-- Replaces the need for a separate Organization node type.
-- "BPH IEEE", "Divisi C&M", a manager's subordinates — all Circles.
-- Speaker-scoped aliases ("anak gua") are resolved by the agent via vault context, not schema.
CREATE NODE TABLE Circle (
    id STRING,
    name STRING,                    -- Canonical group name
    aliases STRING[],               -- Informal references to this group
    content STRING,                 -- Alfred's description of this group and its purpose
    history STRING[],               -- Past content values, newest first. Format: "YYYY-MM-DD HH:MM - [narrative]"
    created_at TIMESTAMP,
    needs_clarification BOOLEAN,
    PRIMARY KEY (id)
);
```

### Relationship Schema (LadybugDB DDL)

```sql
-- Person participated in an Event. Role captures their function (nullable).
-- Grounded in data: Rafid had "sambutan" role at SoTQ — participation alone is insufficient.
CREATE REL TABLE PARTICIPANT_IN (
    FROM Person TO Event,
    role STRING                     -- "sambutan" | "peserta" | "panitia" | "timekeeper" | null
);

-- Task is assigned to a Person (the doer)
CREATE REL TABLE ASSIGNED_TO (FROM Task TO Person);

-- Task was requested by a Person (the requester, may differ from assignee)
CREATE REL TABLE REQUESTED_BY (FROM Task TO Person);

-- Task belongs to an Event or Circle scope
CREATE REL TABLE PART_OF (
    FROM Task TO Event,
    FROM Task TO Circle             -- Task scoped to a Circle (e.g. a division's responsibility)
);

-- Insight is directed at a Person or Circle
CREATE REL TABLE DIR_TOWARDS (
    FROM Insight TO Person,
    FROM Insight TO Circle          -- Insight about a group dynamic
);

-- Causal relationship: one node directly caused another
-- Grounded in data: SoTQ Event caused multiple Tasks; metkuan deadline caused urgency Tasks
CREATE REL TABLE CAUSED_BY (
    FROM Task TO Event,             -- Task arose from an Event
    FROM Task TO Task,              -- Task arose from another Task (explicit only, not inferred)
    FROM Event TO Event,            -- Event caused another Event (e.g. cancellation → reschedule)
    created_at TIMESTAMP
);

-- Insight is supported by observed evidence
-- An Insight's confidence is only as good as its evidence chain
CREATE REL TABLE EVIDENCED_BY (
    FROM Insight TO Task,           -- Insight supported by a Task observation
    FROM Insight TO Event,          -- Insight supported by an Event observation
    FROM Insight TO ConversationBlock, -- Insight supported by a specific block
    observed_at TIMESTAMP
);

-- Nightwatch conflict detection: two Insights make contradicting claims
-- Never created by extraction pipeline — only by Nightwatch
CREATE REL TABLE CONTRADICTS (
    FROM Insight TO Insight,
    detected_at TIMESTAMP,
    resolved BOOLEAN                -- False until user resolves via Memory Review Inbox
);

-- Person→Person relationship. Structural descriptor, not behavioral (behavioral lives in Insights).
-- Descriptor uses canonical vocabulary from skill.md but is not a hard enum.
CREATE REL TABLE KNOWS (
    FROM Person TO Person,
    descriptor STRING,              -- "teman dekat" | "rekan" | "senior" | "junior" | "kenalan"
    context STRING,                 -- Where this relationship exists
    since TIMESTAMP
);

-- Person belongs to a Circle, with an optional role
-- Role is what enables speaker-scoped alias resolution ("anak gua" → kadiv's subordinates)
CREATE REL TABLE MEMBER_OF (
    FROM Person TO Circle,
    role STRING,                    -- "kadiv" | "manager" | "staff" | null
    since TIMESTAMP
);

-- Every node traces back to its originating ConversationBlock
CREATE REL TABLE SOURCED_FROM (
    FROM Person TO ConversationBlock,
    FROM Event TO ConversationBlock,
    FROM Task TO ConversationBlock,
    FROM Insight TO ConversationBlock,
    FROM Circle TO ConversationBlock
);
```

### Storage & Purging Strategy
- **Raw WhatsApp Messages:** Retained for 30 days as a temporary recovery buffer, then purged.
- **LadybugDB:** Permanent storage engine. Nodes are rarely deleted; bad data is the only deletion trigger.

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

### Agent Write Access During Query Flow
The agent can write during query flow via `update_node` and `upsert_reminder`. If the user instructs a change ("mark that done", "that was Rafid's task not mine"), the agent executes it directly using tooling — no separate update pipeline. Updates follow the same `content` → `history` mechanism as extraction updates.

### Full Agent Toolkit

| Tool | Parameters | Output | Purpose |
|---|---|---|---|
| `search_nodes` | `keywords: string` | List of node IDs & titles | Entry point search using LadybugDB-native HNSW vector index |
| `get_surrounding_graph` | `node_ids: list` | Top 15 ranked adjacent node IDs & edge types | Traversal map from current position |
| `read_nodes` | `node_ids: list` | Array of full node data | Opens up to 5 nodes simultaneously |
| `ask_user_for_hint` | `question: string` | Text response from user | Pauses loop and requests a clarifying clue |
| `check_reminders` | `task_ref?: string` | List of existing reminder rows | Checks SQLite before inserting to prevent duplicates |
| `upsert_reminder` | `message, deadline, status, task_ref?` | Confirmation | Inserts or updates a reminder row in SQLite |
| `update_node` | `node_id: string, fields: map` | Confirmation | Updates fields on an existing node (status, content, due_date, etc). Moves old content to history. Used during query flow when user instructs a correction or status change. |

> **Note:** Tool list is incomplete. Additional tools will be identified during implementation — particularly around node creation during query flow, edge manipulation, and reminder deletion.

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

### Deadline Changes
When a task's `due_date` changes, the agent calls `upsert_reminder` with the updated deadline. The unique index on `(task_ref, deadline)` handles this as a delete + insert. The old reminder row is replaced.

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

---

## 13. Project Phases

### Phase 1 — Core Loop
**Goal:** Ingest real WhatsApp data and query it. Nothing else matters until this works. If the memory model holds up here, the thesis is proven.

**Pipelines:**
- Ingestion & Extraction (Phase 1 + Phase 2)
- Alfred Chat / Query Flow

**Done when:** Alfred can watch a live group chat, extract nodes correctly, and answer natural language questions about what it saw.

---

### Phase 2 — Proactive Secretary
**Goal:** Alfred becomes useful day-to-day. Catches obligations and surfaces them before they're missed.

**Pipelines:**
- Reminder Creation
- Reminder Dispatch (cron)
- needs_clarification Push

**Done when:** Alfred proactively pushes deadline reminders and surfaces ambiguous obligations without being asked.

---

### Phase 3 — Graph Maintenance
**Goal:** The graph stays healthy over time as data accumulates.

**Pipelines:**
- Nightwatch (dedup, stale flagging, pre-purge sweep)
- Nightwatch Merge (user-driven via Memory Review Inbox)
- Embedding Generation

**Done when:** Nightwatch runs nightly without incident, the Memory Review Inbox surfaces real duplicates, and search quality holds up as the graph grows.

---

### Phase 4 — Polish & Expansion
**Goal:** Hardening, extensibility, and personalisation post-validation.

- Error correction flow (PWA-driven)
- Off-site backup strategy
- Multi-source ingestion (Telegram, Email)
- User-taught extraction rules (self-updating skill.md)

---

## 14. Open Questions

### ⚫ Critical Priority

**Q1: Agentic Pipeline Specification**
The turn-by-turn logic of every pipeline needs to be fully specced — what the agent sees at each step, what it decides, what tools it calls, and under what conditions it terminates or escalates. This will surface missing tools and edge cases that can't be caught from the schema alone. Pipelines to spec:

1. **Ingestion & Extraction** — webhook → block builder → commit → Phase 1 extract → Phase 2 link → write to graph
2. **Alfred Chat / Query Flow** — user message → traversal loop → optional write → response
3. **Reminder Creation** — extraction or query flow → check → upsert SQLite
4. **Reminder Dispatch** — cron → scan SQLite → push notification
5. **Nightwatch** — nightly → dedup detection → stale flagging → pre-purge sweep
6. **needs_clarification Push** — flagged node created → format → push notification → surface in Memory Review Inbox
7. **Nightwatch Merge** — user swipes merge → surviving node picked → edges re-pointed → duplicate deleted
8. **Embedding Generation** — node created/updated → call external API → store vector

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

**Q3: User-Taught Extraction Rules (Self-Updating skill.md)**
The user could indirectly prompt-engineer Alfred's extraction criteria by telling it what to save or ignore. Alfred would then rewrite its own extraction skill.md to reflect those preferences (e.g. "stop saving random event announcements unless I'm directly involved"). Scope-creep risk for now — noted as a future personalisation feature post-validation.

---

## 15. Decision Log

| # | Decision | Rationale | Date |
|---|---|---|---|
| 1 | **No local vector models** | Squeezing vector models into 4GB RAM is a bottleneck. Free lightweight external APIs handle embeddings instead. | Jun 12, 2026 |
| 2 | **Memory Review Inbox** | Low-confidence merges or conflicts are pushed to a user inbox rather than merged silently. Acts as a spaced-repetition system for the user's own life events. | Jun 12, 2026 |
| 3 | **Link-by-Link Traversal Tooling** | No raw Cypher queries. The agent uses a restricted set of tools for human-like graph navigation, avoiding syntax errors and hallucinated queries. Tool list is expected to grow during implementation. | Jun 12, 2026 |
| 4 | **Stamina + Hard Cap Rule** | Base stamina 2-3 turns to fail fast. Warm leads grant +1. Hard cap at 5-6 prevents token drain and infinite loops. | Jun 12, 2026 |
| 5 | **Insight Table (not Fact Table)** | "Facts" are too rigid. Insights capture qualitative values like character traits, vibes, shared memories, and emotional dynamics. | Jun 12, 2026 |
| 6 | ~~**Polymorphic Causality Edges**~~ | ~~Polymorphic `TRIGGERED_BY` REL tables allow any node to cause, link to, or influence any other node organically.~~ Superseded by Decision 35. | Jun 12, 2026 |
| 7 | **Pause on Webhook** | All background jobs must instantly yield when a new WAHA message arrives to avoid data hazards and DB locks. | Jun 12, 2026 |
| 8 | **Edge Ranking in Traversal** | Adjacent nodes ranked by recency and semantic relevance, top 15 only. Prevents context choking on highly connected hub nodes. | Jun 12, 2026 |
| 9 | **Indonesian Nodes, English Brain** | Node content stored in Indonesian to match real search terms. Agent reasoning done in English for stronger logical performance. | Jun 12, 2026 |
| 10 | **UUIDs + Semantic Identity Resolution** | Raw JIDs discarded as keys due to `@lid` mutation risk. Stable UUIDs used; incoming senders matched via fuzzy semantic matching. | Jun 12, 2026 |
| 11 | **Migrate to LadybugDB** | Apple's October 2025 acquisition of Kùzu Inc. and subsequent repo archiving makes Kuzu unmaintained. LadybugDB is the direct open-source successor. | Jun 12, 2026 |
| 12 | **Golang Backend** | ~15MB idle RAM, fast, native goroutines for concurrent webhook handling, seamless CGO interop with LadybugDB. | Jun 12, 2026 |
| 13 | **Cross-Compilation DevOps** | Compiling LadybugDB's C++ bindings natively on the VPS would OOM crash it. Compiled locally via multi-stage Docker, deployed as a static binary. | Jun 12, 2026 |
| 14 | ~~**Non-Blocking Async Markdown Writes**~~ | ~~Disk I/O writes offloaded to background goroutines via Go channels. Keeps database transactions and ingestion at peak speed.~~ Superseded by Decision 43 — markdown files cut entirely. | Jun 12, 2026 |
| 15 | **The Alfred Persona** | Prevents token bloat and hallucinated therapist-speak. Loyal, dry, professional, ego-centric secretary voice. | Jun 12, 2026 |
| 16 | **Pre-Purge Open Ends Sweep Only** | Never audit all expiring messages blindly. Only sweep unresolved open blocks on Day 30 to close them gracefully. | Jun 12, 2026 |
| 17 | **Flutter → PWA** | Alfred is currently single-user; Flutter's build pipeline adds complexity with no gain at this stage. PWA covers all required features (push notifications, chat, observability, swipe inbox) and opens across all devices via URL. iOS dev burden avoided. PWA served as a static build from the Go backend. Multi-user architecture planned post-validation. | Jun 13, 2026 |
| 18 | **JWT Credential Gate** | PWA is browser-accessible via URL — a security gate is required. JWT login screen before anything renders. No session persistence. Security barrier only, not a user management system. Chosen over Basic Auth for cleaner future multi-user extensibility. | Jun 13, 2026 |
| 19 | **SQLite for Reminder Storage** | Reminders are operational/transient data, not knowledge. Storing as LadybugDB nodes would overlap with Task nodes. SQLite is a single file, no server process, concurrent-safe. Chosen over JSON file for concurrency safety and native query support. | Jun 13, 2026 |
| 20 | **Reminders Owned by Main Agent, Not Nightwatch** | Nightwatch is database maintenance only. Reminder creation is a side effect of extraction and query flows — the agent has the context to decide what is reminder-worthy. | Jun 13, 2026 |
| 21 | **Dumb Cron for Reminder Dispatch** | The cron that fires push notifications is intentionally LLM-free. Pure deadline scanner: query SQLite, fire push notif, mark sent. LLM only involved upstream (writing reminders) and downstream if clarification needed. | Jun 13, 2026 |
| 22 | **Immediate Push on needs_clarification** | Anything flagged needs_clarification is pushed to the user immediately upon commit, not queued for the cron. Applies system-wide. | Jun 13, 2026 |
| 23 | **Prompt Caching** | System prompt is long and resent on every LLM call including every traversal loop turn. Groq prompt caching halves input token costs on repeated prefixes and cached tokens don't count toward rate limits. Must be enabled on all extraction and traversal calls. | Jun 13, 2026 |
| 24 | **Two-Phase Extraction: Extract then Link** | Extraction (reading the conversation) and linking (reading the graph) are cognitively distinct tasks. Phase 1 reads the committed block + prior block summary, outputs candidate nodes with no graph access. Phase 2 traverses the vault to dedup and link. Conflating both in one prompt degrades quality on both. | Jun 13, 2026 |
| 25 | **Explicit Keyword Linking Only** | Nodes are only linked if the source text explicitly mentions or references the target node by name or alias. No implicit linking based on the LLM's background knowledge (e.g. two people both being from IEEE is not sufficient to link a node to the IEEE organisation). Every edge must be justifiable from the node's own content. | Jun 13, 2026 |
| 26 | **Aliases as First-Class Concept Across All Node Types** | Aliases are not just a Person concern. Events, Tasks, and other nodes are referred to informally in Indonesian chat ("rapat tadi", "meeting kemarin"). The linking agent in Phase 2 must search against aliases, not just canonical names, otherwise the explicit keyword rule breaks on informal language. Schema redesign required. | Jun 13, 2026 |
| 27 | **Ambiguous Extractions Always Fail to needs_clarification** | When the extraction pipeline cannot resolve context (e.g. a payment task with no stated reason), it creates a partial node flagged needs_clarification and triggers an immediate push to the user. It never guesses or silently creates incomplete nodes. | Jun 13, 2026 |
| 28 | **WAHA Quote Payload Preserved in ConversationBlock** | WhatsApp reply-to metadata (the quoted message) is load-bearing for correct extraction. "2" as a reply to "maap yah gereja" is a non-commitment; without the quote context it could be misread as volunteering. The quote payload from WAHA must be preserved in the ConversationBlock transcript, not stripped before LLM processing. | Jun 13, 2026 |
| 29 | **Relevance Filter: BPH Perspective** | Alfred extracts for the owner as an IEEE BPH member. Multi-node commits are allowed but each node is evaluated independently against this relevance filter. Organisational events and tasks are likely relevant even when the owner is not explicitly named, given the group chat context. | Jun 13, 2026 |
| 30 | **Phase 2 Linking: Search Aliases + Explicit Match** | The Phase 2 linking agent searches the vault for candidate nodes sharing explicit keywords (including aliases) with the new node. A link is created only if the new node's content directly references the candidate. Semantic similarity alone is insufficient — the match must be grounded in the text. | Jun 13, 2026 |
| 31 | **Person is a Pure Identity Anchor** | Person node holds only identity fields: name, aliases, phone_number, is_self, needs_clarification. No content body, no bio, no cached summary. A person's substance is their surrounding subgraph. Denormalized caches (bio fields) break graph atomicity and create dual-source-of-truth problems. | Jun 13, 2026 |
| 32 | **Unified content STRING replaces scattered summary fields** | All node types except Person carry a single `content STRING` field — Alfred's dry Indonesian narrative of what the node represents. Replaces the inconsistent `summary` fields on Event, Task, Insight, ConversationBlock. ConversationBlock additionally carries `raw_transcript STRING` for the unprocessed message log. | Jun 13, 2026 |
| 33 | **needs_clarification is first-class on all node types** | Any field that cannot be grounded in source text or vault data is left null and `needs_clarification` set true. Never inferred or guessed. Triggers immediate push to user and surfaces in Memory Review Inbox. A confident wrong node is worse than an honest gap. | Jun 13, 2026 |
| 34 | **aliases STRING[] on all node types** | Informal references ("rapat tadi", "sotq", "event kemarin") must be searchable by the Phase 2 linking agent. Aliases are not just a Person concern — Events, Tasks, and Insights are referenced informally in Indonesian chat. | Jun 13, 2026 |
| 35 | **TRIGGERED_BY split into CAUSED_BY, EVIDENCED_BY, CONTRADICTS** | TRIGGERED_BY was semantically overloaded across causality, evidence, and contradiction. Split into three typed REL tables so the pointer fan-out strategy can follow specific edge types deliberately. CONTRADICTS is Nightwatch-only — never written by the extraction pipeline. | Jun 13, 2026 |
| 36 | **PARTICIPANT_IN carries role STRING** | Participation alone is insufficient. Rafid had "sambutan" role at SoTQ — this is meaningfully different from a regular attendee. Role is nullable for general participation. | Jun 13, 2026 |
| 37 | **KNOWS REL table added for Person→Person** | Structural relationship descriptor between people. Captures role/closeness ("teman dekat", "senior", "rekan") and context ("IEEE BPH", "kuliah"). Behavioral dynamics live in Insights, not on this edge. Descriptor uses canonical vocabulary from skill.md but is not a hard enum. | Jun 13, 2026 |
| 38 | **Event status replaces is_confirmed BOOLEAN** | Boolean is too coarse. Real events are planned, active, completed, cancelled, or stale — not just confirmed/unconfirmed. Status enum: "planned" \| "active" \| "completed" \| "cancelled" \| "stale". | Jun 13, 2026 |
| 39 | **ConversationBlock stores raw_transcript** | Raw message log including WAHA quote payloads must be preserved on the node, not stripped before LLM processing. Quote payloads are load-bearing for correct reply-chain interpretation. content field holds Alfred's post-extraction summary. | Jun 13, 2026 |
| 40 | **Circle node type added** | Named groups of people with shared context ("BPH IEEE", "Divisi C&M", a manager's subordinates). Replaces the need for a separate Organization node. Speaker-scoped aliases ("anak gua") are resolved by the agent via vault context (MEMBER_OF role field), not by schema encoding. | Jun 13, 2026 |
| 41 | **verbatim STRING field on Task and Insight** | Quotes live inside nodes, not as separate node types. A Quote is never a navigation destination — it's supporting material. verbatim is populated only when exact wording carries meaning that Alfred's paraphrase would lose (explicit commitments, certificate numbers, strong direct statements). Nullable on all nodes that carry it. | Jun 13, 2026 |
| 42 | **MEMBER_OF REL table added** | Person→Circle relationship with role STRING. Role ("kadiv", "manager", "staff") is what enables the agent to resolve speaker-scoped aliases — if Syazana is kadiv of a Circle, "anak gua" from Syazana resolves to that Circle's members. | Jun 13, 2026 |
| 43 | **Markdown Node Files Cut** | Markdown files were a derived mirror of LadybugDB — not used by the agent, not used for search, only for human inspection. The PWA observability layer and graph itself serve that purpose. Pointless complexity removed entirely. | Jun 13, 2026 |
| 44 | **Nodes Modified In Place, Not Versioned** | No duplicate nodes, no version chains. When new information updates a node, `content` is overwritten with the current truth. Old content is prepended to `history STRING[]` before overwrite. Graph stays clean; history is accessible without traversal. | Jun 13, 2026 |
| 45 | **`history STRING[]` — Human-Readable Changelog** | Each history entry is a self-contained string: `"YYYY-MM-DD HH:MM - [narrative]"`. Newest entries prepended (index 0 = most recent). The LLM reads it like a changelog — no parsing, no parallel arrays, no extra schema. `history` is only read when the query explicitly requires it. | Jun 13, 2026 |
| 46 | **`content` Must Self-Signal Change** | When overwriting `content`, the extraction pipeline (via skill.md template) must acknowledge the prior state in the new narrative. E.g. "Awalnya tugas ini milik Bahlil, sekarang dialihkan ke kamu." The agent reading `content` alone knows a change occurred without touching `history`. | Jun 13, 2026 |
| 47 | **`created_at TIMESTAMP` Added to All Non-Person Node Types** | Enables temporal anchor queries ("what was the first task at X") without traversing history. `created_at` is set once at node creation and never updated. ConversationBlock already had it. Person excluded — identity anchor, no content lifecycle. | Jun 13, 2026 |