package main

import (
	"fmt"
	"strings"
)

func findStringEnd(s string) int {
	for i := 0; i < len(s); i++ {
		if s[i] == '\'' {
			escaped := false
			j := i - 1
			for j >= 0 && s[j] == '\\' {
				escaped = !escaped
				j--
			}
			if !escaped {
				return i
			}
		}
	}
	return -1
}

func main() {
	query := "MATCH (n) WHERE n.id = 'task_alprog_ad46be' RETURN n.id"
	
	if strings.HasPrefix(query, "MATCH (n) WHERE n.id =") && strings.Contains(query, "RETURN n.id") {
		fmt.Println("Matched existence check block!")
		id := ""
		if idStart := strings.Index(query, "n.id = '"); idStart != -1 {
			idEnd := findStringEnd(query[idStart+8:])
			if idEnd != -1 {
				id = query[idStart+8 : idStart+8+idEnd]
			}
		}
		fmt.Printf("Extracted ID: '%s'\n", id)
		
		mockNodes := [][]any{
			{"task_alprog_ad46be", "Task", "some content", nil},
		}
		
		rows := [][]any{}
		for _, n := range mockNodes {
			if n[0].(string) == id {
				rows = append(rows, []any{id})
				break
			}
		}
		fmt.Printf("Rows length: %d\n", len(rows))
	} else if strings.HasPrefix(query, "MATCH (n) WHERE n.id =") && strings.Contains(query, "SET ") {
		fmt.Println("Matched SET block instead!")
	} else {
		fmt.Println("Matched nothing!")
	}
}
