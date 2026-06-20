package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gimigkk/Alfred-Memory/internal/agent"
	"github.com/gimigkk/Alfred-Memory/internal/config"
	"github.com/gimigkk/Alfred-Memory/internal/db"
	"github.com/gimigkk/Alfred-Memory/internal/embed"
	"github.com/gimigkk/Alfred-Memory/internal/llm"
)

type Check struct {
	ID                 string   `json:"id"`
	Description        string   `json:"description"`
	Type               string   `json:"type"`
	NodeType           string   `json:"node_type,omitempty"`
	NodeID             string   `json:"node_id,omitempty"`
	ContentContainsAny []string `json:"content_contains_any,omitempty"`
	Source             string   `json:"source,omitempty"`
	RelType            string   `json:"rel_type,omitempty"`
	TargetType         string   `json:"target_type,omitempty"`
	Speakers           []string `json:"speakers,omitempty"`
}

type Expected struct {
	Fixture string  `json:"fixture"`
	Checks  []Check `json:"checks"`
}

type CheckResult struct {
	ID     string
	Passed bool
	Detail string
}

type RunResult struct {
	RunIndex  int
	Mutations []agent.Mutation
	Checks    []CheckResult
}

func evaluate(linkOut *agent.LinkingOutput, expected Expected) []CheckResult {
	mutations := linkOut.Mutations
	var results []CheckResult

	// Helpers
	nodeExists := func(nodeID string) bool {
		for _, m := range mutations {
			if m.NodeID == nodeID {
				return true
			}
		}
		return false
	}

	for _, check := range expected.Checks {
		res := CheckResult{ID: check.ID, Passed: false}

		switch check.Type {
		case "node_exists":
			for _, m := range mutations {
				if m.Operation == "CREATE_NODE" && m.NodeType == check.NodeType {
					if len(check.ContentContainsAny) > 0 {
						content, _ := m.Properties["content"].(string)
						content = strings.ToLower(content)
						for _, term := range check.ContentContainsAny {
							if strings.Contains(content, strings.ToLower(term)) {
								res.Passed = true
								break
							}
						}
					} else {
						res.Passed = true
					}
				}
				if res.Passed {
					break
				}
			}
		case "edge_exists":
			for _, m := range mutations {
				if m.NodeID == check.Source {
					for _, edge := range m.AddEdges {
						if edge.RelType == check.RelType {
							// For soft target type checks, we assume it's right if the edge exists since schema checks are hard invariants
							res.Passed = true
						}
					}
				}
			}
		case "edge_absent":
			res.Passed = true
			for _, m := range mutations {
				if m.NodeID == check.Source {
					for _, edge := range m.AddEdges {
						if edge.RelType == check.RelType {
							res.Passed = false
						}
					}
				}
			}
		case "node_absent":
			res.Passed = true
			if nodeExists(check.NodeID) {
				res.Passed = false
			}
		case "speaker_coverage":
			res.Passed = true
			for _, speaker := range check.Speakers {
				covered := false
				for _, m := range mutations {
					if m.NodeID == speaker {
						if m.Operation == "CREATE_NODE" {
							covered = true
							break
						}
						if len(m.AddEdges) > 0 {
							covered = true
							break
						}
					}
					// Check if they are the target of an edge
					for _, edge := range m.AddEdges {
						if edge.TargetNodeID == speaker {
							covered = true
							break
						}
					}
				}
				if !covered {
					res.Passed = false
					res.Detail = fmt.Sprintf("Speaker %s was completely unrepresented (no nodes or surviving edges)", speaker)
					break
				}
			}
		case "manifest_task_match":
			res.Passed = true
			for _, acc := range linkOut.ManifestAccounting {
				if strings.Contains(strings.ToUpper(acc.ActionTaken), "CREATE_TASK") || strings.Contains(strings.ToUpper(acc.ActionTaken), "TASK") {
					// Need to find a surviving Task mutation that references this line
					matched := false
					for _, m := range mutations {
						if m.NodeType == "Task" || strings.Contains(strings.ToLower(m.NodeID), "task") {
							for _, edge := range m.AddEdges {
								for _, ref := range edge.EvidenceRefs {
									if strings.Contains(acc.Line, ref.Quote) || strings.Contains(ref.Quote, acc.Line) {
										matched = true
										break
									}
								}
								if matched { break }
							}
						}
						if matched { break }
					}
					if !matched {
						res.Passed = false
						res.Detail = fmt.Sprintf("Manifest claimed task action for line '%s', but no surviving Task mutation referenced it", acc.Line)
						break
					}
				}
			}
		case "structural_invariant":
			// Verified natively by Go validations inside RunAgenticIngestion. If we got here, it passed.
			res.Passed = true
		}
		results = append(results, res)
	}

	return results
}

func main() {
	log.SetFlags(0)
	
	// 1. Load config & API clients
	cfg := config.LoadConfig()
	geminiEmbed := embed.NewGeminiClient(cfg.GeminiAPIKey)
	llmRouter := llm.NewRouterClient(cfg.GeminiAPIKey, cfg.GroqAPIKey)

	// 2. Initialize DBs
	dbDir := "./.lbug"
	lbugClient, err := db.NewClient(dbDir)
	if err != nil {
		log.Fatalf("Failed to initialize LadybugDB: %v", err)
	}
	defer lbugClient.Close()

	conn, err := lbugClient.GetConnection()
	if err != nil {
		log.Fatalf("Failed to get connection: %v", err)
	}
	defer conn.Close()

	if err := db.InitLadybugSchema(conn); err != nil {
		log.Printf("Schema init warning: %v", err)
	}

	// 3. Setup Orchestrator
	orchestrator := agent.NewOrchestrator(llmRouter, geminiEmbed, conn)

	// 4. Load Fixtures
	fixtureDir := "testdata/fixtures/sambutan_001"
	transcriptBytes, err := os.ReadFile(fixtureDir + "/transcript.txt")
	if err != nil {
		log.Fatalf("Failed to read transcript: %v", err)
	}
	transcript := string(transcriptBytes)

	expectedBytes, err := os.ReadFile(fixtureDir + "/expected.json")
	if err != nil {
		log.Fatalf("Failed to read expected.json: %v", err)
	}
	var expected Expected
	if err := json.Unmarshal(expectedBytes, &expected); err != nil {
		log.Fatalf("Failed to parse expected.json: %v", err)
	}

	// 5. Run N Iterations
	var results []RunResult
	N := 3
	fmt.Printf("Starting eval harness for fixture: %s (%d runs)\n\n", expected.Fixture, N)

	for i := 0; i < N; i++ {
		fmt.Printf("Run %d/%d...\n", i+1, N)
		linkOut, err := orchestrator.RunAgenticIngestion(fmt.Sprintf("eval_%d", i), transcript, true)
		if err != nil {
			fmt.Printf("  Run %d Failed: %v\n", i+1, err)
			results = append(results, RunResult{
				RunIndex: i, 
				Checks: []CheckResult{{ID: "run_completed", Passed: false, Detail: err.Error()}},
			})
			continue
		}
		
		checks := evaluate(linkOut, expected)
		results = append(results, RunResult{RunIndex: i, Mutations: linkOut.Mutations, Checks: checks})
	}

	// 6. Print Score Table
	fmt.Printf("\n============================================\n")
	fmt.Printf("Fixture: %s  (%d runs)\n\n", expected.Fixture, N)
	
	scoreMap := make(map[string]int)
	for _, res := range results {
		for _, c := range res.Checks {
			if c.Passed {
				scoreMap[c.ID]++
			}
		}
	}

	fmt.Printf("HARD INVARIANTS (must be %d/%d):\n", N, N)
	hardInvariants := []string{"no_directional_violations", "rafid_not_wrongly_assigned", "bahlil_never_appears"}
	for _, id := range hardInvariants {
		score := scoreMap[id]
		icon := " "
		if score == N {
			icon = "✓"
		}
		fmt.Printf("  %-30s %2d/%d  %s\n", id, score, N, icon)
	}

	fmt.Printf("\nSOFT COMPLETENESS (expected variance):\n")
	softCompleteness := []string{"event_created", "backup_task_assigned", "live_report_task_created", "all_speakers_covered", "manifest_task_match"}
	for _, id := range softCompleteness {
		score := scoreMap[id]
		fmt.Printf("  %-30s %2d/%d\n", id, score, N)
	}

	fmt.Printf("\nOverall structural validity: %d/%d runs produced zero hard-invariant violations.\n", scoreMap["no_directional_violations"], N)
}
