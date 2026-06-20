# Phase 0: Go Hybrid GraphRAG Port - Implementation Plan

This document serves as the absolute "bible" for **Phase 0**. Any AI agent working on this phase must refer to this document for structural, architectural, and procedural constraints.

---

## 1. Phase Context & Goals
**Objective:** Build a standalone, robust Hybrid GraphRAG package in Go on top of LadybugDB. 
**Why:** The entirety of Alfred's memory system relies on `go-ladybug`. Because it is an early, low-activity binding, we must validate the foundation (specifically memory safety and concurrent connection handling) in an isolated ~200-line testbench before building the massive Alfred infrastructure on top of it.
**Success Criteria:** 
1. Ingest dummy documents, embed them using Gemini, and store them in LadybugDB.
2. Successfully execute the 4-stage Hybrid RAG retrieval (Vector → Graph Expand → RRF → PageRank).
3. Survive a brutal concurrent stress test without crashing (validating the CGO GC/`Close()` fix).

---

## 2. Architecture & Components

### A. Core Engine
- **Language:** Go (1.22+)
- **Database:** LadybugDB running in-process via CGO bindings (`github.com/LadybugDB/go-ladybug`).
- **Database File:** Local `.lbug` directory.

### B. The Embedding Layer
- **Provider:** Google Gemini API
- **Model:** `text-embedding-004`
- **Why:** Massive free tier (1,500 requests/day), excellent multilingual support (Indonesian), and removes the need to load vector models into limited RAM.
- **Output:** 768-dimensional float32 vector.

### C. Data Schema (Phase 0 Test Bench)
We use a stripped-down schema to test the mechanics without the complexity of Alfred's persona.
- **Node Table:** `Document (id STRING, content STRING)`
- **Edge Table:** `LINKS_TO (FROM Document TO Document, context STRING)`
- **Vector Index:** `CREATE VECTOR INDEX doc_idx ON Document(content)`

### D. The Retrieval Pipeline (The "rag" package)
1. **Embed:** Call Gemini API to get the vector for the user's query string.
2. **Vector Search:** Execute Cypher `CALL QUERY_VECTOR_INDEX(...)` to get the Top-K semantic hits.
3. **Graph Expand:** Execute Cypher to find nodes 1-hop away from the vector hits.
4. **Rank:** Execute Cypher `CALL pagerank(...)`. Merge the vector similarity score with the PageRank score using **Reciprocal Rank Fusion (RRF)** in native Go logic.

### E. The LLM Synthesis Layer
After the subgraph is retrieved, it must be formatted and passed to an LLM to prove that an agent can actually read and reason over the raw graph data.
- **Provider:** Groq API (Primary inference provider).
- **Model:** `llama3-70b-8192` (or similar fast Groq model).

---

## 3. Directory Structure

The Go module will be strictly isolated in a new folder: `/home/gimigkk/Desktop/Projects/Alfred/phase0-rag/`

```text
phase0-rag/
├── main.go               # Entry point, CLI test runner
├── go.mod                
├── .env                  # Stores GEMINI_API_KEY
├── config/               
│   └── env.go            # Loads and validates environment variables
├── db/                   
│   └── client.go         # LadybugDB connection lifecycle (STRICT explicit Close() enforcement)
├── schema/               
│   ├── init.go           # DDL creation and Vector Index setup
│   └── seed.go           # Injects 5-10 connected dummy nodes for testing
├── embed/                
│   └── gemini.go         # HTTP client for text-embedding-004
├── llm/
│   └── groq.go           # Groq API client for final answer synthesis
├── rag/                  
│   ├── types.go          # Structs for Subgraph return types
│   ├── hybrid.go         # The 4-stage retrieval orchestrator
│   └── rrf.go            # Math logic for Reciprocal Rank Fusion
└── stress/               
    └── concurrent.go     # Multi-goroutine lifecycle stress test
```

---

## 4. Step-by-Step Execution Plan

### Step 1: Initialization & External APIs
*Prerequisite: Ensure `gcc` and `g++` are installed on the Linux host, as CGO requires a C++ compiler to build `go-ladybug`.*
1. Initialize the Go module (`go mod init`).
2. Implement the `config` package to load the `.env` file (requires `GEMINI_API_KEY` and `GROQ_API_KEY`).
3. Implement `embed/gemini.go`: A simple REST client to fetch embeddings.
4. Implement `llm/groq.go`: A simple REST client pointing to `https://api.groq.com/openai/v1/chat/completions` to pass the context and get an answer.

### Step 2: Database Layer & Schema Setup
1. Import `github.com/LadybugDB/go-ladybug`.
2. Implement `db/client.go`. **Critical Rule:** Every `Connection`, `QueryResult`, and `FlatTuple` allocated via CGO *must* have an explicit `.Close()` call using `defer`. We cannot rely on the Go Garbage Collector.
3. Implement `schema/init.go`: Write the Cypher queries to create the `Document` node, `LINKS_TO` edge, and the Vector Index.
4. Implement `schema/seed.go`: Create a hardcoded set of ~10 interconnected nodes, fetch their Gemini embeddings, and insert them.

### Step 3: The Hybrid Retrieval Engine
1. Implement `rag/hybrid.go`. The main function `QueryRAG(query string) Subgraph` will:
   - Call `embed.GetVector(query)`.
   - Run vector search Cypher to get top hits.
   - Run a second Cypher query to expand edges from those hits.
   - Fetch PageRank scores.
2. Implement `rag/rrf.go`: Calculate `1 / (k + rank)` for both the vector hit list and the PageRank score list, combine them, and sort the final subgraph.

### Step 4: LLM Synthesis (The "Generation" in RAG)
1. In `rag/hybrid.go` (or `main.go`), take the sorted subgraph output from Step 3.
2. Format the nodes and edges into a readable string/JSON context block.
3. Pass the user's original query + the context block to `llm/groq.go`.
4. Return the natural language answer. This proves the LLM can comprehend the graph structure.

### Step 5: The Concurrency Stress Test
This is the most important part of Phase 0.
1. Implement `stress/concurrent.go`.
2. Spin up a `sync.WaitGroup` with 50 goroutines.
3. Inside each goroutine, open a *new* LadybugDB connection to the shared `.lbug` file. Run a random read or write query. Ensure the `QueryResult` and `Connection` are explicitly closed.
4. If this triggers a `SIGSEGV` (Segmentation Fault) or a C++ lock error, the architecture fails here and must be debugged before moving to Phase 1.

### Step 6: Wrap & Review
1. Wire everything into `main.go` so running `go run .` executes the setup, seeds the database, runs a test query printing the retrieved subgraph **and the LLM's answer**, and then runs the stress test.
2. If successful, this entire codebase becomes the foundational dependency for Alfred in Phase 1.
