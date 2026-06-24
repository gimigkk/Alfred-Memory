# Alfred's Memory

> Chat logs are not memory. Memory is an evolving world model.

Humans do not store life as transcripts. We compress experiences into relationships, expectations, and stories that help us act in the future.

[***Alfred***](./Alfred.md) is a biomimetic agentic memory system that converts fragmented conversations into an evolving knowledge graph. Experiences are extracted, linked, revised, and retrieved as context rather than replayed as history. The graph is not the memory itself. It is the shape memory leaves behind.

---

## High-Level Architecture

Alfred is composed of three primary layers that work in tandem to construct and maintain the world model:

1. **The Go Orchestrator (`cmd/alfred`)**: The strict, mechanical backbone of the system. It handles webhook debouncing, manages database connections, intercepts LLM hallucinations, and executes Cypher queries. 
2. **The Dynamic ReAct Loop (`internal/agent`)**: A highly specialized Gemini-based agent that iteratively queries the existing graph (via RAG), identifies speakers, infers context, and proposes structured JSON graph mutations.
3. **The Knowledge Graph (`LadybugDB` / `mock`)**: An advanced graph topology that stores entities (`Person`, `Task`, `Event`, `Insight`, `Circle`) and their semantic relationships (`ASSIGNED_TO`, `MENTIONED_IN`, `HAS_ROLE`, `PART_OF`, `LINKS_TO`, etc.) rather than flat conversational logs.

---

## Architectural Novelty

Alfred solves several fundamental problems with LLM-based memory agents through mechanical constraints rather than relying purely on prompt engineering:

### 1. Additive-with-Pruning Context Caching
LLMs suffer from "Lost in the Middle" instruction decay. If you put 2,000 tokens of strict JSON schema rules at the top of a prompt, the LLM will fail at nuanced RAG discovery. 
**The Solution:** Alfred starts the ReAct loop with a purely exploratory persona. Once the agent finishes RAG and signals `[REQUEST_SCHEMA]`, a Go Interceptor pauses the loop, sweeps the chat history to delete old schema constraints, and *appends* the massive rulebook as a fresh `User` message at the absolute bottom of the context window. This mathematically forces peak Recency Bias on schema compliance at the exact moment the agent writes its final JSON payload.

### 2. Mechanical Deflection (Guardrails > Prompting)
We don't just ask the LLM nicely to output valid graph structures. The Go Orchestrator acts as an iron-clad bouncer. If the LLM generates a bad edge direction (e.g., pointing `MENTIONED_IN` backward), forgets to link a participant, or hallucinates nested properties, Go's strict JSON `DisallowUnknownFields` and Gate Logic throw a hard `fmt.Errorf`. That literal error string is fed straight back into the ReAct loop as a User message, forcing the LLM to dynamically self-correct its own hallucinations without breaking the pipeline.

### 3. Layer 1 vs Layer 2 Graph Segregation
Inline creation of structural/organizational nodes (like `Circle`) in a messy group chat is highly brittle and leads to the LLM hallucinating five duplicate nodes for "divisi acara". 
**The Solution:** The Orchestrator hard-rejects any attempt to create a `Circle` during Layer 1 Ingestion. Instead, the agent is forced to capture structural mentions purely as passive data (`group_mentions` JSON arrays on transient Tasks). A deterministic Layer 2 cron job later sweeps these mentions, clusters them safely, and promotes them into permanent `Circle` nodes.

### 4. Atomic Temporal History
LLMs are terrible at updating complex chronological arrays. Instead of forcing the LLM to format and prepend old node content to a history log, the Go backend natively intercepts `UPDATE_NODE` operations. It dynamically generates Atomic Cypher that shifts the current node state into an immutable `history STRING[]` array before executing the `SET` for the new state, guaranteeing zero data loss.

### 5. Node-Level Verbatim (No Chat Logs)
Storing raw transcripts as graph nodes pollutes the world model. Alfred completely excised the `ConversationBlock` node concept. Instead, exact quotes and their origin timestamps are attached directly to the semantic edges (`evidence_refs`) and the nodes themselves (`verbatim`). The graph contains only the extracted knowledge, anchored by exact citations, leaving the noise behind.

---

## Current Status

We are currently deep into **Phase 1: The Core Ingestion Engine**.
- [x] Modular Prompt Refactoring & Dynamic ReAct logic.
- [x] Strict JSON guardrails, DisallowUnknownFields, and hallucination deflection.
- [x] Atomic Cypher temporal history tracking.
- [ ] SQLite Reminders integration (pending).
- [ ] Layer 2 Mention Promotion cron job (pending).
- [ ] The Chat Agent endpoint and PWA frontend (pending).

For the full, extensive decision log, edge definitions, and node topology, see [Alfred.md](./Alfred.md) and the sprint tracker in [Phase_1_Plan.md](./Phase_1_Plan.md).
