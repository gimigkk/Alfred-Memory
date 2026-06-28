package agent

import (
	"regexp"
	"strings"
)

func normalizeStr(s string) string {
	reg := regexp.MustCompile("[^a-zA-Z0-9]+")
	return strings.ToLower(reg.ReplaceAllString(s, ""))
}

func tokenize(s string) []string {
	reg := regexp.MustCompile("[^a-zA-Z0-9]+")
	parts := reg.Split(s, -1)
	var tokens []string
	for _, p := range parts {
		p = strings.ToLower(p)
		if p != "" {
			tokens = append(tokens, p)
		}
	}
	return tokens
}

func isTokenSubset(sub, main []string) bool {
	if len(sub) == 0 {
		return false
	}
	mainSet := make(map[string]bool, len(main))
	for _, t := range main {
		mainSet[t] = true
	}
	for _, t := range sub {
		if !mainSet[t] {
			return false
		}
	}
	return true
}

func hasTokenOverlap(sub, main []string) bool {
	if len(sub) == 0 {
		return false
	}
	mainSet := make(map[string]bool, len(main))
	for _, t := range main {
		mainSet[t] = true
	}
	for _, t := range sub {
		if mainSet[t] {
			return true
		}
	}
	return false
}

// longestCommonRun returns the length of the longest contiguous run of tokens
// shared between a and b, preserving order (classic longest-common-substring,
// applied to token sequences rather than characters). Used where a single shared
// token would be too weak a signal but exact whole-sequence containment (isSubslice)
// would be too strict — e.g. matching a transcript line against a paraphrased
// node content where only a meaningful phrase, not the whole line, recurs verbatim.
func longestCommonRun(a, b []string) int {
	if len(a) == 0 || len(b) == 0 {
		return 0
	}
	prev := make([]int, len(b)+1)
	best := 0
	for i := 1; i <= len(a); i++ {
		curr := make([]int, len(b)+1)
		for j := 1; j <= len(b); j++ {
			if a[i-1] == b[j-1] {
				curr[j] = prev[j-1] + 1
				if curr[j] > best {
					best = curr[j]
				}
			}
		}
		prev = curr
	}
	return best
}

func filterNoiseTokens(tokens []string) []string {
	var filtered []string
	numericRe := regexp.MustCompile(`^[0-9]+$`)
	for _, t := range tokens {
		if numericRe.MatchString(t) {
			continue // drop purely numeric cohort years / IDs
		}
		// TODO: "ieee" is dataset-specific; consider a configurable noise-word list if Alfred is used beyond this org
		if t == "ieee" {
			continue // drop cohort/org suffix
		}
		if len(t) < 3 {
			continue // drop short tokens like "m", "w"
		}
		filtered = append(filtered, t)
	}
	if len(filtered) == 0 {
		return tokens // filtering would remove everything — fall back to unfiltered
	}
	return filtered
}
