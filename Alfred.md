# Alfred - Agentic Memory System
### Architecture & Design Specification

> **Status:** Planning / Pre-development
> **Last updated:** June 14, 2026 (revalidated + GraphRAG architecture)
> **Stack:** Golang Backend · LadybugDB · SQLite · Gemini API (Router) · WAHA (GOWS) · PWA Frontend · Ubuntu 24.04 VPS (4GB RAM / 8GB Storage)

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

Most AI memory systems either ask you to manually log information, or they treat memory as a flat retrieval problem - embed everything, search by similarity, call it done. Alfred is neither.

**The claim Alfred is built to prove:** an agent can maintain a coherent, temporally-aware, self-correcting knowledge graph from noisy, informal, multilingual conversational data - without any user curation. Not just storing what you said, but knowing *when you said it*, *whether it's still true*, *how confident to be*, and *how things connect to each other*.

The secretary framing is the test case, not the product. WhatsApp group chats are the data source because real data is messy - informal language, implicit references, contradictions, half-finished thoughts, Indonesian slang mid-sentence. If the memory model holds up there, it holds up anywhere.

### The Concrete Problem

The specific failure mode Alfred is built to solve has two distinct forms:

**Missing** - a message comes in while you're busy, you skim it, you don't register that it contained an obligation directed at you. The information never entered your head. This is the harder case: Alfred must catch what you didn't.

**Forgetting** - you saw it, you knew you had to do something, it got buried, and it fell out of working memory before you acted. Alfred must surface it before the deadline passes.

Both failures are common in group chats where obligations are often implicit, informally phrased, or directed at you without your name being mentioned. Alfred's proactive push - surfacing `needs_clarification` flags and deadline reminders without being asked - is the most important feature in the system. The knowledge graph and traversal system exist to make those pushes accurate.

### Design Philosophy: Honest Gap Over Confident Wrong

Alfred is explicitly designed to keep the user in the loop rather than quietly automate their judgment away. When the extraction pipeline cannot attribute a task with sufficient confidence, it creates a partial node flagged `needs_clarification` and pushes it to the user immediately. It never guesses.

This is intentional. A real secretary surfaces ambiguous obligations and asks - they don't silently assign them to the wrong person. The volume of `needs_clarification` flags is not a failure metric; it is an accurate reflection of how ambiguous group chat obligations actually are. The user decides. Alfred catches and presents.

### What Alfred Does
- **Remember** - extract and store structured facts, tasks, events, preferences, people, experiences, and social insights
- **Forget intentionally** - raw chat is ephemeral; only curated memory is permanent
- **Summarize** - compress conversation blocks into semantic summaries
- **Merge duplicates** - resolve conflicts between overlapping pieces of information
- **Detect stale info** - flag or update outdated beliefs/states while maintaining historical lineage
- **Track changes over time** - store not just the current state, but the history of how it got there
- **Remind proactively** - surface upcoming deadlines and obligations without being asked
- **Explain its reasoning** - all agent actions, tool calls, and traversal steps are visible and auditable

### Analogy
Not just a smart notes app. Alfred is closer to a personal epistemics engine - something that knows not just *what* you know, but *when you learned it*, *whether it's still true*, and *how confident to be*. The secretary is the interface. The knowledge graph is the point.

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
│   Gemini (primary) → Flash/Pro fallback chain               │
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
│   Core node types: Person, Event, Task, Insight, Circle     │
└────────────────────────────┬────────────────────────────────┘
                             │
               Tool-callable memory interface
                             │ 
                             ▼
┌─────────────────────────────────────────────────────────────┐
│            AGENTIC QUERY LAYER (HYBRID GRAPHRAG)            │
│   query_rag(query, top_k?, hops?) -> subgraph               │
│     stages: embed, vector search (HNSW), graph              │
│     expand (1-2 hops), RRF fusion, PageRank                 │
│   Other tools: ask_user_for_hint, check_reminders,          │
│                 upsert_reminder, create_node,               │
│                 update_node, delete_node                    │
│   LLM decides which tools to call; Go executes              │
│   No stamina/hard cap - query_rag auto-runs at              │
│   query start; agent re-queries or writes as                │
│   needed mid-reasoning                                      │
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
5. The user queries Alfred via the PWA chat → `query_rag` runs (Hybrid GraphRAG) and the agent reasons over the returned subgraph, re-querying or writing as needed
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

**Race Condition Prevention:** Any active background job (Nightwatch cleanup, indexing) must be immediately paused the instant a new webhook message arrives. Ingestion and extraction have strict priority over all background work. See Open Questions - Ingestion Queue Architecture.

### Message Attribution: "Me" vs "Others"
WAHA's webhook payload includes a `fromMe: true | false` boolean on every message - this is the ground-truth signal for distinguishing messages sent by the owner vs messages sent by contacts. It is set at ingestion, not inferred by the LLM.

A fixed **singleton "self" Person node** (`name: "me"`, `is_self: true`) anchors all `fromMe: true` messages. This ensures task and commitment attribution is unambiguous without per-conversation guesswork.

---

## 5. Conversation Block System

The conversation block is the **fundamental unit of processing** - analogous to a Git commit. Alfred never processes individual messages in isolation; it always works on committed blocks.

### Block Lifecycle
- **Buffer:** Messages accumulate while a conversation is active.
- **Debounce:** The block commits after a silence threshold (15-30 minutes of no messages) or when a hard topic shift is detected.
- **Open:** A block that has been started but not yet committed.
- **Committed:** A fully sealed block ready for LLM extraction.
- **Abandoned:** A block that was interrupted and never naturally closed. Defaulted after a max-age threshold. On Day 30, any open block reaching the purge window gets a final LLM sweep before raw deletion (see Section 10).

Blocks are strictly **per-chat**. The source `chat_id` is the enforced boundary. Messages from different group chats never share a block. Note: Conversation Blocks are purely transient memory queues; they do *not* become nodes in the graph (replaced by node-level `verbatim` and `evidence_refs`).

---

## 6. LLM Extraction Pipeline

### Dynamic ReAct Architecture & Prompt Decoupling
To prevent context window degradation and recency bias, the pipeline uses a **Dynamic ReAct Architecture**. The monolithic prompt is decoupled into distinct skills (`skill_discovery.md`, `skill_commit.md`).
- **Discovery Phase:** The agent starts with only the discovery instructions. It must use `extract_transcript_manifest` to sequentially enumerate every line, and `query_rag` to fetch required context from the vault.
- **Stateful Middleware Interceptor:** Once the agent resolves entities, it outputs the `[REQUEST_SCHEMA]` token. The Go Orchestrator intercepts this, prunes all intermediate "thought" clutter from the conversation array, and dynamically injects the rigid graph schema (`skill_commit.md`) at the peak of the LLM's context window.

### Agentic Ingestion Loop
The loop is driven by the Go Orchestrator and utilizes three primary tools:
1. **`extract_transcript_manifest`**: The LLM must call this first. If lines are skipped, the orchestrator hard-rejects the extraction.
2. **`query_rag`**: Fetches subgraph context from the vault to resolve aliases and check for duplicate events. The interceptor blocks any schema requests until the vault is queried.
3. **`commit_mutations`**: Outputs the final JSON manifest of `CREATE_NODE` and `UPDATE_NODE` instructions. **Strictly transactional**: if validation fails, the entire batch is rejected.

### Schema-Based Graph Construction & Clarity Constraints
To prevent LLM "tunnel vision" and graph hallucination, the ingestion agent operates under strict structural constraints:
- **Schema-Defined Topology:** Graph topologies are explicitly defined within the schema constraints of the prompt. The LLM mathematically maps outgoing edges (e.g., `PART_OF` for Tasks, `ASSIGNED_TO`/`MENTIONED_IN` for Persons) directly into the mutation payloads.
- **The Clarity Checklist:** The "Default to Uncertainty" rule requires all new entities to start with `needs_clarification: true`. To toggle this to `false`, the LLM is forced to output a strict 5-point checklist (Who, What, When, Where, Why) directly inside the `clarification_basis` JSON property. There are **no "operationally necessary" exceptions**—if any of the 5Ws are implied, the clarification flag must trigger.

### Go Structural Validation Layer
Before any mutation touches LadybugDB, it is intercepted by a multi-pass Go validation layer within the Orchestrator. 
- **Directional Enforcement:** Edges must strictly originate from `Person` or `Task` and point outward. Inverse structures are hard-rejected.
- **Null Hypothesis (Events):** Forces the agent to prove an event match via at least **two explicit, unique keywords** rather than hallucinating semantic links.
- **User Resolution:** Ensures "The User" is natively resolved via `query_rag`.
- **Clarity Guard Override:** If the `clarification_basis` mentions that a field is "unknown", the Go orchestrator forcibly overrides `needs_clarification` to true, regardless of what the LLM claimed.
Any mutation failing structural invariants logs a `[REJECTED EDGE]` and is stripped out, protecting the graph topology from hallucination.

### Deterministic Evaluation Harness
To prevent regression, the pipeline includes an automated evaluation harness (`cmd/eval/main.go`). It runs the agentic loop `N` times against fixed transcript fixtures in a `dryRun` state (bypassing DB commits). It produces a strict pass-rate table grading two categories:
- **Hard Invariants:** Must be 100% (e.g., directional integrity, no fabricated user nodes).
- **Soft Completeness:** Grades variance in reasoning (e.g., task creation rate, speaker coverage).

### Cross-Block Event Deduplication
The agent enforces strict explicit keyword matching. If no alias matches but the block is temporally proximate to a known event, the agent still flags `needs_clarification` and pushes to the Memory Review Inbox rather than silently merging.

### Language Boundaries
- **Storage (Database properties):** Stored explicitly in **Indonesian**. 
- **Reasoning (Agent Logic):** Done strictly in **English**. The LLM uses English for its internal monologue, tool selection, and structure-parsing.

### LLM Fallback Chain
Gemini API is the primary inference provider (via `llmRouter`). If a call hits rate limits (HTTP 429), the system walks down a prioritized list of fallback models (e.g., Gemini Pro Preview -> Flash -> Custom Tools) in a try/catch loop until one succeeds, ensuring zero dropped blocks.

---

## 7. Memory Vault

### Database: LadybugDB
Following the October 2025 acquisition of Kùzu Inc. by Apple and the subsequent archiving of the Kuzu repository, this project uses **LadybugDB** - the direct open-source community successor. It runs **in-process** inside the compiled Go binary via CGO bindings (`go-ladybug`), maintaining an ultra-lightweight memory footprint suitable for a 4GB VPS.

> [!WARNING]
> **TEMPORARY MOCK PIVOT:** The actual C++ LadybugDB implementation has been temporarily swapped out for a pure **in-memory Go mock** (`internal/ladybug/mock.go`) to bypass CGO compilation issues encountered during early pipeline testing. 
> 
> All graph entities (Bahlil, BEM, etc.) are hardcoded directly into the `mockNodes` array in memory. **Every time the server restarts, the database starts completely fresh and automatically seeded.** There is NO on-disk persistence (the `.lbug` directory is unused), and no external seeder scripts or `os.RemoveAll()` logic should be used.
> 
> *For details, see `docs/architecture/database_mock_layer.md`.*

### Identity Resolution (Bypassing the `@lid` Bug)
Raw WhatsApp JIDs mutate and rotate dynamically across different clients, making them unreliable as database keys.
- `Person.id` is a **generated stable UUID**.
- Incoming messages are resolved to existing `Person` nodes using fuzzy/semantic matching on `name`, `aliases`, and normalized `phone_number`.

### Node Schema (LadybugDB DDL)

> **Content model:** Every node type except Person carries a `content STRING` field - Alfred's dry narrative of what this node represents, written in Indonesian. This is the primary field the linking agent reads in Phase 2. Structured metadata fields (status, dates, priority) exist alongside it for querying and Nightwatch logic. Person is a pure identity anchor - its substance lives in the surrounding graph.

> **Aliases:** Every node type carries `aliases STRING[]`. Aliases are informal names, abbreviations, or references this node is known by in chat ("rapat tadi", "sotq", "event kemarin"). The Phase 2 linking agent searches aliases, not just canonical names, to satisfy the explicit keyword rule.

> **History & Updates:** Nodes are modified in place. `content` always reflects current truth. When `content` is overwritten, the old value is prepended to `history STRING[]` as a human-readable entry (`"YYYY-MM-DD HH:MM - [narrative]"`), newest first. The new `content` must briefly acknowledge the prior state so the agent is hinted that something changed without needing to read `history`. `history` is only traversed when a query explicitly requires it.

> **needs_clarification:** Every node type carries `needs_clarification BOOLEAN`. When set true, the node was created with one or more fields that could not be grounded in source text or vault data. An immediate push notification is sent to the user. The node is visible in the Memory Review Inbox until resolved. Alfred never guesses to fill a field - a confident wrong node is worse than an honest gap.

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

-- 2. Event: Any occurrence - meetings, sessions, social events, deadlines
CREATE NODE TABLE Event (
    id STRING,
    name STRING,                    -- Canonical name: "SoTQ IEEE 2026"
    aliases STRING[],               -- ["sotq", "acara tadi", "event kemarin"]
    content STRING,                 -- Alfred's narrative of current state, written in Indonesian. Must acknowledge prior state if overwriting.
    history STRING[],               -- Past content values, newest first. Format: "YYYY-MM-DD HH:MM - [narrative]"
    verbatim STRING,                -- Exact quote referencing the event (nullable)
    group_mentions STRING,          -- JSON Array of group mentions captured in this event
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
    group_mentions STRING,          -- JSON Array of group mentions captured in this task
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



-- 6. Circle: A named group of people with shared context
-- Replaces the need for a separate Organization node type.
-- "BPH IEEE", "Divisi C&M", a manager's subordinates - all Circles.
-- CRITICAL NOTE: The ingestion LLM is mechanically forbidden from creating Circle nodes directly.
-- Instead, it captures organizational references inside the `group_mentions` JSON array on Task and Event nodes.
-- A Layer 2 "Promoter" batch job asynchronously clusters these mentions and creates permanent Circle nodes.
CREATE NODE TABLE Circle (
    id STRING,
    name STRING,                    -- Canonical group name
    aliases STRING[],               -- Informal references to this group
    content STRING,                 -- Alfred's description of this group and its purpose
    history STRING[],               -- Past content values, newest first. Format: "YYYY-MM-DD HH:MM - [narrative]"
    verbatim STRING,                -- Exact quote referencing the group (nullable)
    created_at TIMESTAMP,
    needs_clarification BOOLEAN,
    PRIMARY KEY (id)
);
```

### Relationship Schema (LadybugDB DDL)

```sql
-- Task is assigned to a Person (the doer)
CREATE REL TABLE ASSIGNED_TO (
    FROM Person TO Task,
    evidence_refs STRING
);

-- Person is mentioned in a Task or Event (passive participant or subject)
CREATE REL TABLE MENTIONED_IN (
    FROM Person TO Task,
    FROM Person TO Event,
    evidence_refs STRING
);

-- Person has a titled role in an Event
CREATE REL TABLE HAS_ROLE (
    FROM Person TO Event,
    evidence_refs STRING
);

-- Person belongs to a Circle, with an optional role
-- Role is what enables speaker-scoped alias resolution ("anak gua" → kadiv's subordinates)
CREATE REL TABLE MEMBER_OF (
    FROM Person TO Circle,
    role STRING,                    -- "kadiv" | "manager" | "staff" | null
    since TIMESTAMP,
    evidence_refs STRING
);

-- Task belongs to an Event or Circle scope
CREATE REL TABLE PART_OF (
    FROM Task TO Event,
    FROM Task TO Circle,            -- Task scoped to a Circle (e.g. a division's responsibility)
    evidence_refs STRING
);

-- Insight is directed at a Person or Circle
CREATE REL TABLE DIR_TOWARDS (
    FROM Insight TO Person,
    FROM Insight TO Circle,         -- Insight about a group dynamic
    evidence_refs STRING
);


-- Universal generic link for causal relationships (replaces CAUSED_BY and EVIDENCED_BY)
CREATE REL TABLE LINKS_TO (
    FROM Task TO Event,
    FROM Task TO Insight,
    FROM Task TO Task,
    FROM Event TO Insight,
    FROM Event TO Event,
    FROM Insight TO Insight,
    context STRING,
    evidence_refs STRING
);

-- Nightwatch conflict detection: two Insights make contradicting claims
-- Never created by extraction pipeline - only by Nightwatch
CREATE REL TABLE CONTRADICTS (
    FROM Insight TO Insight,
    detected_at TIMESTAMP,
    resolved BOOLEAN                -- False until user resolves via Memory Review Inbox
);

-- Person→Person relationship. Structural descriptor, not behavioral (behavioral lives in Insights).
-- Descriptor uses canonical vocabulary from core_schema.md but is not a hard enum.
CREATE REL TABLE KNOWS (
    FROM Person TO Person,
    descriptor STRING,              -- "teman dekat" | "rekan" | "senior" | "junior" | "kenalan"
    context STRING,                 -- Where this relationship exists
    since TIMESTAMP,
    evidence_refs STRING
);

### Storage & Purging Strategy
- **Raw WhatsApp Messages:** Retained for 30 days as a temporary recovery buffer, then purged.
- **LadybugDB:** Permanent storage engine. Nodes are rarely deleted; bad data is the only deletion trigger.

---

## 8. Agentic Query System

### Retrieval: Hybrid GraphRAG
Alfred uses a **Hybrid GraphRAG** retrieval pattern — vector search and graph traversal run in parallel, results are fused via RRF, and PageRank weights the final ranking. This replaces the original hand-rolled stamina-gated traversal loop, which was a worse, token-heavier reimplementation of the same idea.

The pattern is ported from [`Volland/ladybug-rag`](https://github.com/Volland/ladybug-rag) (Python reference, ~300 lines) into Go as part of Phase 0. All four stages run natively inside LadybugDB via Cypher — no external services, no Python sidecar.

**Four retrieval stages:**

1. **Embed the query** — HTTP call to the pinned embedding API (same model used at write time)
2. **Vector search** — `QUERY_VECTOR_INDEX` Cypher query returns top-K semantically similar nodes (HNSW)
3. **Graph expand** — second Cypher query follows edges 1-2 hops from vector hits, using LadybugDB's native graph traversal. The traversed edges (type, direction, properties) are carried through, not discarded.
4. **RRF fusion + PageRank ranking** — ~20 lines of Go math merges both ranked lists; PageRank weights structurally important nodes higher regardless of recency. If an edge connects two nodes where one endpoint falls outside the final ranked set, the edge is dropped along with it rather than left dangling.

The pipeline is exposed to the agent as a callable tool (`query_rag`) with configurable `top_k` and `hops` parameters. It runs automatically on query start with default params, but the agent can call it again mid-reasoning with different params if the initial context is insufficient — e.g. a narrower query string, deeper hop expansion, or larger `top_k`. No stamina system, no hard cap, no multi-turn traversal loop.

### query_rag Output Shape
`query_rag` returns a **subgraph**, not a flat node list: `{ nodes: [...], edges: [...] }`.
- Every node carries at minimum `id`, `node_type`, and its schema fields (`content`, `name`, `aliases`, status fields, etc.) — `id` and `node_type` are the handles every write tool keys off.
- Every edge carries `from_id`, `to_id`, `rel_type`, and direction, plus any edge-only properties (`HAS_ROLE.evidence_refs`, `KNOWS.descriptor`/`context`, `MEMBER_OF.role`, `LINKS_TO.context`, etc.) — this is metadata that lives only on the relationship and has nowhere else to surface.

This preserves the structural context that the old manual traversal loop made explicit step-by-step, and gives any future relationship-editing tool the edge identity it needs to act.

### Agent Write Access During Query Flow
The agent can write during query flow. If the user instructs a change ("mark that done", "that was Rafid's task not mine"), the agent executes it directly via tools — no separate update pipeline. Updates follow the same `content` → `history` mechanism as extraction updates.

When re-querying via `query_rag` still can't surface enough context, the agent falls back to `ask_user_for_hint`.

### Alias & Full-Text Search
BM25 full-text search via LadybugDB's native extension is used for alias matching in Phase 2 linking. Replaces the original custom keyword matching logic.

### Agent Toolkit

| Tool | Parameters | Output | Purpose |
|---|---|---|---|
| `query_rag` | `query: string, top_k?: int, hops?: int` | Subgraph: `{ nodes: [...], edges: [...] }` | Runs the full Hybrid GraphRAG pipeline. Auto-called at query start; callable again mid-reasoning. |
| `create_node` | `node_type: string, fields: map` | Confirmation | Creates a new memory mid-chat. Fields include `content`, `aliases`, and `edges`. |
| `update_node` | `node_id: string, fields: map` | Confirmation | Mutates a node. **Fields include `add_edges` and `remove_edges`**. The LLM must provide a new `content` narrative whenever modifying relationships, preventing graph rot. Old content goes to history. |
| `delete_node` | `node_id: string` | Confirmation | Hard deletes a node. Go automatically handles cascade deletes of connected edges and reminders. |
| `upsert_reminder` | `message, deadline, status, task_ref?` | Confirmation | Inserts or updates a deadline reminder in SQLite. |
| `check_reminders` | `task_ref?: string` | List of rows | Checks SQLite for existing reminders to prevent duplicates. |
| `ask_user_for_hint` | `question: string` | Text response | Requests a clarifying clue from the user when graph context is insufficient. |

> **Note:** Standalone `create_edge` and `delete_edge` tools are explicitly banned (Decision 59). Edge mutations must happen via `update_node` to ensure the node's text `content` remains synchronized with its structural relationships.

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
The **main agent** is solely responsible for writing reminders - never Nightwatch. This happens in two flows:
- **Extraction pipeline** - when a committed block contains a user-owned task or deadline
- **User query flow** - when the agent surfaces a task during traversal and a reminder is warranted

Before inserting, the agent calls `check_reminders` to prevent duplicates. The `UNIQUE` index is a database-level safety net.

### Status Lifecycle

```
pending → sent                  (cron fires push notif)
pending → dismissed             (user dismisses from PWA)
pending → dismissed             (Go auto-dismisses when Task status → completed/abandoned/stale)
pending → needs_clarification   (agent unsure what the task actually is)
needs_clarification → pending   (user clarifies, agent updates)
needs_clarification → dismissed
```

### Cron Job (Dumb Scanner - No LLM)
A Go cron job runs on a configurable interval (e.g. every hour):
1. Query SQLite for `status = 'pending'` rows where `deadline` is within the notification window
2. Fire push notification via PWA Push API
3. Update `status` to `sent`

The cron never touches LadybugDB and never calls the LLM.

### Immediate Push on `needs_clarification`
Any reminder or node flagged `needs_clarification` is pushed to the user **immediately upon commit** - it does not wait for the cron. This rule applies system-wide: the cron handles scheduled pending reminders only; anything requiring user input is surfaced right away.

### Cascading Deletion
If a Task node is deleted from LadybugDB (bad data), its associated reminder rows cascade-delete via `task_ref`.

### Auto-Dismissal on Status Change
If the agent calls `update_node` to change a Task's status to `completed`, `abandoned`, or `stale`, the Go backend automatically intercepts this and marks any associated pending reminder in SQLite as `dismissed`. The agent does not need to manually manage the reminder state in these cases.

### Deadline Changes
When a task's `due_date` changes, the agent calls `upsert_reminder` with the updated deadline. The unique index on `(task_ref, deadline)` handles this as a delete + insert. The old reminder row is replaced.

---

## 10. Background Agents

### Nightwatch (Database Maintenance Agent)
Nightwatch runs during low-traffic hours (nightly). Its **sole responsibility is database maintenance** - it does not write reminders and does not handle user-facing logic.

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
A JWT credential gate loads before the PWA renders anything. The user provides a username and password; the Go backend validates and returns a JWT which the PWA holds in memory and attaches to all subsequent API requests. No session persistence - re-login on refresh is acceptable. This is a security barrier only, not a user management system.

### Chat Interface
Natural language Q&A with Alfred's memory. The user types a question; the agent runs the Hybrid GraphRAG query flow (`query_rag` auto-called, agent reasons over the returned subgraph, re-queries or writes if needed) and responds in Alfred's persona.

### Observability Layer
An interactive debug log inside the chat view (similar to Claude's "thinking" blocks). The user can expand a dropdown during or after a query to watch the agent's exact step-by-step journey:
- *Tool Call: query_rag("Bahlil tugas Bunga")* → subgraph with `Task[bahlil_task]`, `Person[bahlil]`, `Person[bunga]`, edge `ASSIGNED_TO(bahlil_task -> bahlil)`
- *Thought: "This task is about Bunga's event, currently assigned to Bahlil."*
- *Tool Call: update_node(task_id="bahlil_task", fields={status: "completed"})*

### Memory Review Inbox
A Tinder-style swipe interface. When Nightwatch finds potentially duplicate or conflicting nodes, it pushes them here rather than merging silently. The user taps to merge or swipes to keep separate. This doubles as a spaced-repetition system for the user's own life events.

### Push Notifications
Delivered via the PWA Push API + Service Workers. Service Workers run in the background even when the browser is closed, enabling reliable delivery on Android. iOS is supported since Safari 16.4 (2023).

---

## 12. Infrastructure

### Stack

| Component | Choice | Links | Reason |
|---|---|---|---|
| OS | Ubuntu 24.04 | — | Stable LTS on budget VPS |
| Backend | Golang | [go.dev](https://go.dev) | ~15MB idle RAM, native goroutines, single static binary, CGO interop for LadybugDB |
| Graph DB | LadybugDB via `go-ladybug` ⚠️ | [ladybug](https://github.com/LadybugDB/ladybug) · [go-ladybug](https://github.com/LadybugDB/go-ladybug) · [docs](https://docs.ladybugdb.com) | In-process C++ engine, HNSW vector index, BM25 full-text search, PageRank, Kuzu's open-source successor. **⚠️ go-ladybug is low-activity (15 stars, 3 forks) — validate in Phase 0 before depending on it.** |
| GraphRAG | Port of `ladybug-rag` to Go | [ladybug-rag (Python ref)](https://github.com/Volland/ladybug-rag) | Hybrid retrieval: HNSW vector search + Cypher graph expand + RRF fusion + PageRank. ~100 lines of Go. Ported in Phase 0. |
| Reminder DB | SQLite (`reminders.db`) via `mattn/go-sqlite3` | [mattn/go-sqlite3](https://github.com/mattn/go-sqlite3) | Single file, no server, concurrent-safe, native query support |
| LLM | Gemini API (`gemini-3-pro-preview` / `flash`) | [aistudio.google.com](https://aistudio.google.com) | Strong multi-turn tool calling, robust fallback chain handling |
| Embeddings | Google Gemini API (`gemini-embedding-2`) | — | External API to avoid loading vector models into 4GB RAM. **Must pin one model — embedding model used at write time must match query time exactly.** |
| STT (Onboarding) | Gemini Audio API | [docs](https://console.groq.com/docs/speech-text) | Already in stack. One HTTP call. Late-V1 feature. |
| WhatsApp | WAHA + GOWS WebSockets | — | Proven webhook provider |
| Push Notifications | `SherClockHolmes/webpush-go` | [webpush-go](https://github.com/SherClockHolmes/webpush-go) | VAPID-based Web Push, drop-in for PWA push notif delivery |
| Auth | `golang-jwt/jwt` | [golang-jwt](https://github.com/golang-jwt/jwt) | Standard Go JWT library for credential gate |
| Frontend | PWA | — | Browser-native, no app store, works across all devices via URL, push notif support |

### DevOps & Compilation
`go-ladybug` uses CGO to compile native C++ bindings. Building natively on the 4GB VPS will crash due to insufficient RAM. The build pipeline is:
- Cross-compile locally using multi-stage Docker builds matching the target VPS architecture
- Deploy a pre-built static binary directly to the VPS

---

## 13. Project Phases

### Phase 0 - Go Hybrid GraphRAG Port ✅ (COMPLETED)
**Goal:** Validate the retrieval layer and `go-ladybug` bindings before building Alfred on top of them. Port [`ladybug-rag`](https://github.com/Volland/ladybug-rag) from Python to Go as a standalone package.

**Why this first:** The entire Alfred stack depends on `go-ladybug` — an official but low-activity binding (15 stars, 3 forks) that had known issues in early releases. If it's broken or missing features, finding that out now costs days. Finding it out in Phase 2 costs months. Phase 0 is a contained ~100-line project that validates the foundation and produces a reusable artifact.

**Scope:**
- LadybugDB connection + schema setup in Go via `go-ladybug`
- Store a document node with an embedding vector via external API
- Build HNSW vector index
- Hybrid retrieval: vector search (`QUERY_VECTOR_INDEX`) + Cypher graph expand (1-2 hops) + RRF fusion (~20 lines of Go math) + PageRank weighting
- Pass ranked context to LLM, get answer back
- **Concurrent connection lifecycle stress test** - `go-ladybug` issue #7 documented a real Cgo/GC race (finalizer destroys a `QueryResult` while another goroutine is still reading a derived `FlatTuple` via Cgo → SIGSEGV), fixed by requiring explicit `Close()` rather than relying on the GC. Phase 0 must simulate Alfred's actual access pattern - multiple goroutines (webhook handler, background job, query flow) hitting the same `.lbug` file concurrently, each with disciplined explicit `Close()` on every `Connection`/`QueryResult` - to confirm this doesn't resurface under load.
- Nothing else — no WhatsApp, no extraction pipeline, no Alfred persona

**Reference:** Python implementation at [`Volland/ladybug-rag`](https://github.com/Volland/ladybug-rag) (~300 lines). Read this first, port to Go, validate it produces equivalent results. Note the reference also includes entity extraction, Louvain community detection, and KNN `SIMILAR` edges - Alfred's port deliberately excludes these (extraction is handled by Alfred's own pipeline; community detection has no current consumer and is deferred until Nightwatch needs it).

**Done when:** A Go binary can ingest a handful of test documents, build the graph, run a hybrid query, return a correctly-ranked context to an LLM, and survive the concurrent connection stress test without crashing. Package is clean enough to extract as a standalone library later.

**Output:** A reusable Go hybrid GraphRAG package that Alfred's query system is built on top of. Long-term, publishable as a standalone open-source library — the first decent Go implementation of hybrid GraphRAG on LadybugDB.

---

### Phase 1 - Core Loop
**Goal:** Ingest real WhatsApp data and query it. Nothing else matters until this works. If the memory model holds up here, the thesis is proven.

**Pipelines:**
- Ingestion & Extraction (Phase 1 + Phase 2)
- Alfred Chat / Query Flow

**Done when:** Alfred can watch a live group chat, extract nodes correctly, and answer natural language questions about what it saw.

---

### Phase 2 - Proactive Secretary
**Goal:** Alfred becomes useful day-to-day. Catches obligations and surfaces them before they're missed.

**Pipelines:**
- Reminder Creation
- Reminder Dispatch (cron)
- needs_clarification Push

**Done when:** Alfred proactively pushes deadline reminders and surfaces ambiguous obligations without being asked.

---

### Phase 3 - Graph Maintenance
**Goal:** The graph stays healthy over time as data accumulates.

**Pipelines:**
- Nightwatch (dedup, stale flagging, pre-purge sweep)
- Nightwatch Merge (user-driven via Memory Review Inbox)
- Embedding Generation

**Done when:** Nightwatch runs nightly without incident, the Memory Review Inbox surfaces real duplicates, and search quality holds up as the graph grows.

---

### Phase 4 - Polish & Expansion
**Goal:** Hardening, extensibility, and personalisation post-validation.

- Error correction flow (PWA-driven)
- Off-site backup strategy
- Multi-source ingestion (Telegram, Email)
- User-taught extraction rules (self-updating prompts)

---

## 14. Open Questions

### ⚫ Critical Priority

*(None currently)*

### 🔴 High Priority

**Q1: Ingestion Queue Architecture**
When the WAHA webhook fires, background jobs must immediately pause to avoid database lock hazards. A clean transactional lock mechanism needs to be designed in the Go backend - the exact implementation is undecided.

**Q2: Error Correction Flow**
How does the user correct Alfred when it extracts something wrong? The agent may misinterpret a message, create a wrong node, or link things incorrectly. A deliberate correction mechanism - likely via the PWA - needs to be designed so corrections feed back into the agent and update the graph cleanly without leaving stale data.

### 🟡 Low Priority / Future

**Q1: Off-Site Backup Strategy**
LadybugDB and SQLite are both single files on the VPS. If the VPS dies, all memory is lost. A scheduled backup strategy is needed - candidates: rsync to another machine, or push to object storage (Backblaze B2 or Cloudflare R2). Not important for prototype phase.

**Q2: Multi-Source Architecture**
How do future ingestion sources (Telegram, Email) plug in without rewriting the ingestion layer? Likely a source-agnostic message interface that WAHA and future adapters all conform to. Not a near-future priority.

**Q3: User-Taught Extraction Rules (Self-Updating Prompts)**
The user could indirectly prompt-engineer Alfred's extraction criteria by telling it what to save or ignore. Alfred would then rewrite its own extraction prompts to reflect those preferences (e.g. "stop saving random event announcements unless I'm directly involved"). Scope-creep risk for now - noted as a future personalisation feature post-validation.

---

## 15. Decision Log

| # | Decision | Rationale | Date |
|---|---|---|---|
| 1 | **No local vector models** | 4GB RAM constraint. External embedding APIs used instead. | Jun 12, 2026 |
| 2 | **Memory Review Inbox** | Low-confidence merges/conflicts pushed to user inbox rather than auto-merged. Also acts as spaced repetition for life events. | Jun 12, 2026 |
| 5 | **Insight Table (not Fact Table)** | "Facts" are too rigid. Insights capture qualitative values: character traits, relationship dynamics, preferences, behavioral patterns. | Jun 12, 2026 |
| 7 | **Pause on Webhook** | All background jobs must instantly yield on new WAHA webhook to avoid DB lock hazards. | Jun 12, 2026 |
| 9 | **Indonesian Nodes, English Brain** | Node content in Indonesian matches real search terms for BM25. Agent reasoning in English for stronger logical performance. | Jun 12, 2026 |
| 10 | **UUIDs + Semantic Identity Resolution** | Raw JIDs discarded due to `@lid` mutation risk. Stable UUIDs used; senders matched via fuzzy semantic matching on name/aliases/phone. | Jun 12, 2026 |
| 11 | **Migrate to LadybugDB** | Apple's October 2025 acquisition of Kùzu Inc. and repo archiving makes Kuzu unmaintained. LadybugDB is the direct open-source successor. | Jun 12, 2026 |
| 12 | **Golang Backend** | ~15MB idle RAM, native goroutines, single static binary, CGO interop for LadybugDB. | Jun 12, 2026 |
| 13 | **Cross-Compilation DevOps** | Building LadybugDB's C++ bindings on the 4GB VPS would OOM. Cross-compiled locally via multi-stage Docker, deployed as static binary. | Jun 12, 2026 |
| 15 | **The Alfred Persona** | Loyal, dry, ego-centric secretary voice. Prevents token bloat and hallucinated therapist-speak. | Jun 12, 2026 |
| 16 | **Pre-Purge Open Ends Sweep Only** | On Day 30, only sweep unresolved open blocks — never a blind audit of all expiring messages. | Jun 12, 2026 |
| 17 | **Flutter → PWA** | Single-user now; Flutter's build pipeline adds complexity with no gain. PWA covers all required features and avoids iOS dev burden. Multi-user post-validation. | Jun 13, 2026 |
| 18 | **JWT Credential Gate** | Security barrier before PWA renders. No session persistence. Chosen over Basic Auth for cleaner future multi-user extensibility. | Jun 13, 2026 |
| 19 | **SQLite for Reminder Storage** | Reminders are operational/transient, not knowledge. Single file, no server, concurrent-safe. Kept separate from LadybugDB Task nodes. | Jun 13, 2026 |
| 20 | **Reminders Owned by Main Agent, Not Nightwatch** | Nightwatch is DB maintenance only. Reminder creation requires the agent context that exists in extraction and query flows. | Jun 13, 2026 |
| 21 | **Dumb Cron for Reminder Dispatch** | Intentionally LLM-free. Pure deadline scanner: query SQLite → fire push notif → mark sent. | Jun 13, 2026 |
| 22 | **Immediate Push on needs_clarification** | Flagged nodes push immediately upon commit, not queued for cron. Applies system-wide. | Jun 13, 2026 |
| 23 | **Prompt Caching** | System prompt is resent on every LLM call. Gemini context caching reduces input token costs on repeated prefixes; cached tokens don't count toward rate limits. | Jun 13, 2026 |
| 24 | **Two-Phase Extraction: Extract then Link** | Phase 1 reads the block, outputs candidate nodes with no graph access. Phase 2 traverses the vault to dedup and link. Conflating both degrades quality on both. | Jun 13, 2026 |
| 25 | **Explicit Keyword Linking Only** | Nodes linked only if source text explicitly references the target by name or alias. No implicit linking from LLM background knowledge. Every edge must be justifiable from the node's own content. | Jun 13, 2026 |
| 26 | **Aliases as First-Class Concept Across All Node Types** | Aliases are not just a Person concern. Events, Tasks, and other nodes are referred to informally in Indonesian chat ("rapat tadi", "meeting kemarin"). The linking agent in Phase 2 must search against aliases, not just canonical names, otherwise the explicit keyword rule breaks on informal language. Schema redesign required. | Jun 13, 2026 |
| 27 | **Ambiguous Extractions Always Fail to needs_clarification** | Extraction pipeline never guesses when context is unresolvable. Creates a partial node flagged `needs_clarification` and triggers immediate push. | Jun 13, 2026 |
| 28 | **WAHA Quote Payload Preserved in ConversationBlock** | Reply-to metadata is load-bearing for extraction. "2" as a reply to "maap yah gereja" is a non-commitment — without quote context it reads as volunteering. Quote payload must survive to LLM processing. | Jun 13, 2026 |
| 29 | **Relevance Filter: BPH Perspective** | Extracts for the owner as IEEE BPH member. Each candidate node evaluated independently. Org events/tasks are relevant even when owner isn't explicitly named. | Jun 13, 2026 |
| 30 | **Phase 2 Linking: Search Aliases + Explicit Match** | Link created only if the new node's content directly references the candidate by name or alias. Semantic similarity alone is insufficient. | Jun 13, 2026 |
| 31 | **Person is a Pure Identity Anchor** | Person holds only identity fields. No content body, no cached bio. Substance lives in the surrounding subgraph. Denormalized caches create dual-source-of-truth problems. | Jun 13, 2026 |
| 32 | **Unified content STRING** | All node types except Person carry a single `content STRING` — Alfred's Indonesian narrative of current state. Replaces scattered `summary` fields. ConversationBlock also carries `raw_transcript`. | Jun 13, 2026 |
| 33 | **needs_clarification is first-class on all node types** | Any ungroundable field is left null with `needs_clarification: true`. Never inferred. Triggers immediate push and surfaces in Memory Review Inbox. | Jun 13, 2026 |
| 34 | **aliases STRING[] on all node types** | Informal references must be searchable by Phase 2 linking agent. Not just a Person concern — Events, Tasks, and Insights are referenced informally in Indonesian chat. | Jun 13, 2026 |
| 35 | **TRIGGERED_BY split into CAUSED_BY, EVIDENCED_BY, CONTRADICTS** | TRIGGERED_BY was semantically overloaded. Three typed REL tables allow deliberate edge-type traversal. CONTRADICTS is Nightwatch-only. | Jun 13, 2026 |
| 36 | **PARTICIPANT_IN carries role STRING** | Participation alone is insufficient — Rafid's "sambutan" role at SoTQ is meaningfully different from general attendance. Nullable for regular participants. | Jun 13, 2026 |
| 37 | **KNOWS REL for Person→Person** | Structural descriptor ("teman dekat", "senior") and context ("IEEE BPH"). Behavioral dynamics live in Insights, not on this edge. | Jun 13, 2026 |
| 38 | **Event status enum replaces is_confirmed BOOLEAN** | Boolean is too coarse. Status: `"planned" \| "active" \| "completed" \| "cancelled" \| "stale"`. | Jun 13, 2026 |
| 39 | **Raw Transcripts as Ephemeral** | Raw message logs are no longer stored as ConversationBlock nodes to avoid graph pollution. `content` holds Alfred's post-extraction summary. | Jun 13, 2026 |
| 40 | **Circle node type added** | Named groups with shared context. Replaces separate Organization node. Speaker-scoped aliases ("anak gua") resolved by agent via MEMBER_OF role, not schema encoding. | Jun 13, 2026 |
| 41 | **verbatim STRING on Task and Insight** | Exact wording stored only when paraphrase would lose meaning (explicit commitments, certificate numbers, strong direct statements). Nullable. | Jun 13, 2026 |
| 42 | **MEMBER_OF REL added** | Person→Circle with role STRING. Role ("kadiv", "staff") enables speaker-scoped alias resolution by the agent. | Jun 13, 2026 |
| 43 | **Hybrid GraphRAG + query_rag Tool** | Stamina-gated traversal loop replaced with: embed → HNSW vector search → Cypher graph expand (1-2 hops) → RRF fusion → PageRank ranking. Pipeline exposed to the agent as `query_rag(query, top_k?, hops?)` — auto-called at query start, re-callable mid-reasoning with custom params if initial context is insufficient. Supersedes Decisions 3, 4, and 8. | Jun 14, 2026 |
| 44 | **LadybugDB Native BM25 for Alias Search** | Phase 2 alias matching uses LadybugDB's native BM25 extension instead of custom keyword logic. | Jun 14, 2026 |
| 45 | **LadybugDB PageRank for Node Ranking** | PageRank weights structurally important nodes higher regardless of recency. Replaces custom edge ranking. | Jun 14, 2026 |
| 46 | **Embedding Model Must Be Pinned** | Model at write time must match model at query time. Switching requires full re-embed. Pin in config before first write. | Jun 14, 2026 |
| 47 | **Phase 0: Port ladybug-rag to Go First** | go-ladybug is low-activity (15 stars, 3 forks) with known early binding issues. Validate foundation in ~100-line project before Alfred depends on it. | Jun 14, 2026 |
| 48 | **Attribution Uncertainty vs Duplicate Detection: Different Urgency** | "Task might be yours" = immediate push. "These two nodes might be the same event" = Nightwatch queue. Separated in Memory Review Inbox UX with different push weights. | Jun 14, 2026 |
| 49 | **Trust Calibration as Explicit First-Weeks Goal** | Log how often the user agrees with Alfred's flags vs dismisses them. Real metric for whether the certainty bar is calibrated correctly. | Jun 14, 2026 |
| 50 | **Onboarding: Conversational Seeding, Late V1** | Cold start via conversational flow in Alfred's persona (not a form). Seeded nodes carry `source: "user_declared"` — Nightwatch never auto-merges them without confirmation. STT via Gemini Audio API. Late-V1 feature. | Jun 14, 2026 |
| 51 | **Markdown Node Files Cut** | Derived LadybugDB mirror not used by agent or search — only for human inspection. PWA observability layer serves that purpose. | Jun 13, 2026 |
| 52 | **Nodes Modified In Place, Not Versioned** | `content` overwritten with current truth; old value prepended to `history STRING[]`. No duplicate nodes or version chains. | Jun 13, 2026 |
| 53 | **`history STRING[]` - Human-Readable Changelog** | Entries are self-contained strings: `"YYYY-MM-DD HH:MM - [narrative]"`, newest first. Read by LLM as a changelog without parsing. Only accessed when the query explicitly requires it. | Jun 13, 2026 |
| 54 | **`content` Must Self-Signal Change** | New `content` must acknowledge the prior state (e.g. "Awalnya tugas ini milik Bahlil, sekarang dialihkan ke kamu") so the agent knows something changed without reading `history`. | Jun 13, 2026 |
| 55 | **`created_at TIMESTAMP` on All Non-Person Node Types** | Enables temporal anchor queries without traversing history. Set once at creation, never updated. Person excluded — identity anchor, no content lifecycle. | Jun 13, 2026 |
| 56 | **`query_rag` Returns a Subgraph, Not a Flat Node List** | Graph-expand stage already traverses edges to find neighbors - discarding them would mean re-deriving the same data later for relationship-editing tools. Output is `{ nodes: [...with id, node_type], edges: [...with from_id, to_id, rel_type, properties] }`. Preserves edge-only metadata (`PARTICIPANT_IN.role`, `KNOWS.descriptor`, etc.) and gives write tools stable handles. | Jun 15, 2026 |
| 57 | **Phase 0 Includes a Concurrent Connection Lifecycle Stress Test** | `go-ladybug` issue #7 documented a Cgo/GC race (finalizer destroys `QueryResult` while another goroutine reads a derived `FlatTuple` → SIGSEGV), fixed via explicit `Close()`. Alfred's real access pattern is multi-goroutine (webhook handler, background jobs, query flow) against one `.lbug` file - Phase 0 must simulate this concurrency with disciplined `Close()` calls before Alfred depends on the binding. | Jun 15, 2026 |
| 58 | **Agentic Pipeline Specification Finalized** | Documented the three true autonomous pipelines: Ingestion (linear, ETL-style with RAG intermediary), Chat Flow (non-linear, full CRUD capability), and Nightwatch (cron-triggered graph maintenance). Eliminates manual traversal loops in Phase 2 extraction. | Jun 20, 2026 |
| 59 | **Edges Merged into Node Mutations (Option B)** | To prevent graph rot and history-bypassing, raw `create_edge` / `delete_edge` tools are banned. The agent modifies edges via `update_node`'s `add_edges` / `remove_edges` fields, which forces it to rewrite the node's `content` narrative concurrently. | Jun 20, 2026 |
| 60 | **Phase 0 Findings Locked** | Pinned APIs to `gemini-embedding-2` and `llama-3.3-70b-versatile`. `go-ladybug` requires Go 1.26+ (handled via `GOTOOLCHAIN=auto`) and CGO builds require `lbug.h` host headers. Vector Indexes strictly require a `FLOAT[768]` column rather than a string-based text index. | Jun 20, 2026 |
| 61 | **No Example Bleed in Prompts** | Specific placeholder data in prompt examples causes LLM hallucination and overfitting ("Example Bleed"). Prompt examples must use highly abstract placeholders like `[ABSTRACT_REASON]` to force reliance on raw transcript text. | Jun 21, 2026 |
| 62 | **Semantic Readable Node IDs** | Instead of generating purely random UUIDs (`node_1a2b3c...`) for new nodes, the orchestrator strips the LLM's `temp_` prefix and appends a 6-character hash to the semantic intent (e.g., `event_pembayaran_1a2b3c`). This guarantees DB uniqueness while keeping graph visualizer tooltips human-readable. | Jun 21, 2026 |
| 63 | **ConversationBlock replaced by node-level verbatim** | Dropped ConversationBlock to avoid graph pollution. Provenance is now maintained directly on nodes via 'verbatim' and edges via 'evidence_refs'. | Jun 21, 2026 |
| 64 | **Dynamic ReAct via Stateful Interceptor** | The monolithic ingestion prompt was split into decoupled skills (`discovery`, `commit`). A Go-side interceptor prunes the conversation array mid-flight and injects rigid schema constraints only upon a `[REQUEST_SCHEMA]` token, maximizing context window efficiency and recency bias. | Jun 22, 2026 |
| 65 | **Mandatory Hostile Persona in Simulation** | The Courtroom simulation now explicitly mandates a "Hostile Attacker (Prosecutor)" persona to aggressively stress-test structural vulnerabilities, prompt injection loopholes, and edge cases, preventing echo-chamber approvals for architectural changes. | Jun 22, 2026 |
| 66 | **Strict 5W Clarity Defaults** | Removed the "operationally necessary" loophole from the Clarity Guard. Any Task or Event missing explicit Who/What/When/Where/Why must be flagged `needs_clarification: true`. Null hypothesis for Event linking requires two unique matching keywords, preventing RAG-induced semantic hallucination. | Jun 22, 2026 |
| 67 | **1:1 Parallel Array Target Resolution** | To solve batch cardinality limits in `query_rag` while preserving explicit intent and result-based verification, the agent uses a 1:1 `target_speakers` array matching its `queries` array. A mismatched length instantly rejects the tool call, eliminating hidden state failures. | Jun 23, 2026 |
| 68 | **In-Memory Mock Database Pivot** | Replaced the actual C++ LadybugDB with a pure in-memory Go mock (`internal/ladybug/mock.go`) to bypass severe CGO compilation blocks during Phase 1 pipeline testing. No on-disk persistence until CGO issues are resolved. | Jun 23, 2026 |
| 69 | **Obligations Interceptor Gate** | Added `query_speaker_obligations` as a mandatory Go-side gate before schema request. Forces the agent to query the graph for existing `needs_clarification` nodes to perform temporal updates rather than creating duplicate nodes. | Jun 23, 2026 |
| 70 | **Documentation Suite Redesign** | Core bibles (`Alfred.md`, `Phase_1_Plan.md`) are stripped of low-level mechanics, pushing detailed architecture and flow documentation into `docs/architecture/`. `.geminirules` updated to mandate reading these before coding. | Jun 23, 2026 |
| 71 | **Layer 1 Mention Capture (Circle Deferral)** | Inline `Circle` node creation was identified as highly brittle and hallucination-prone. Circle creation is now deferred to a Layer 2 batch job. The ingestion agent captures structural mentions purely as data using the `group_mentions` property on `Task`/`Event`. | Jun 24, 2026 |
| 72 | **Mechanical Schema Guardrails** | The Go Orchestrator strictly enforces invariant structures: e.g., hard-rejecting `CREATE_NODE` for `Circle`, ensuring properties maps are evaluated accurately. This acts as an absolute backstop against LLM reasoning decay. | Jun 24, 2026 |
---

## 16. Documentation Suite Index

For low-level mechanics, refer to the modular documentation in `docs/`:

- `docs/architecture/database_mock_layer.md` — Explains the in-memory `mock.go` pivot and how to query it.
- `docs/architecture/agent_orchestrator.md` — Details the Go-side interceptor pattern and tool gates.
- `docs/ai_skills/courtroom.md` — Constraints and persona rules for running AI architectural debates.

