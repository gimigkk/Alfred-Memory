# Database Mock Layer (`mock.go`)

> [!WARNING]
> The real C++ LadybugDB has been temporarily disabled. The system currently uses an in-memory mock. **Do not attempt to use `os.RemoveAll` or file-based cleanup scripts, as there is no persistent storage.**

## The Pivot Context
During Phase 1 pipeline testing, severe CGO compilation issues blocked progress with `go-ladybug`. To decouple pipeline development from C/C++ build errors, the actual LadybugDB client was swapped out for a pure **in-memory Go mock** located at `internal/ladybug/mock.go`.

## How It Works
The `internal/ladybug` package exports a mock `Database` and `Connection` struct that perfectly mirrors the API expected by the Go Orchestrator and extraction pipeline. 

### The Source of Truth
All data is stored inside two package-level variables:
```go
var mockNodes = [][]any{ ... }
var mockEdges = [][]any{ ... }
```
When the `main.go` server boots up, these slices are instantiated in memory. This provides an **automatically clean, perfectly seeded test environment on every single run.**

### Query Interception
The mock `Connection.Query()` method uses hardcoded string matching (e.g. `strings.Contains(query, "MATCH (a)-[r]->(b)")`) to return the expected slices to the pipeline.

Mutation queries (`CREATE NODE`, `UPDATE_NODE`, etc.) are similarly intercepted and appended to the in-memory slices.

## Restoring the Real Database
Once the CGO headers are fixed on the host environment:
1. Delete `internal/ladybug/mock.go`.
2. Update all imports referencing `github.com/gimigkk/Alfred-Memory/internal/ladybug` back to `github.com/LadybugDB/go-ladybug`.
3. The `.lbug` directory will resume functioning as the persistent on-disk data store.
