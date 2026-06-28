package main

import (
	"fmt"
	"github.com/gimigkk/Alfred-Memory/internal/ladybug"
)

func main() {
	conn, _ := ladybug.NewConnection(nil)

	// Inject a task node
	createQ := "CREATE (n:Task {id: 'task_alprog_123', content: 'test'})"
	conn.Query(createQ)

	// Run existence check
	checkQ := "MATCH (n) WHERE n.id = 'task_alprog_123' RETURN n.id"
	res, err := conn.Query(checkQ)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else if res == nil {
		fmt.Println("Res is nil")
	} else if !res.HasNext() {
		fmt.Println("Res has no next")
	} else {
		fmt.Printf("Success! Row: %v\n", res.GetNext())
	}
}
