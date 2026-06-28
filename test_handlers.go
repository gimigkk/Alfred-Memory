package main

import (
	"fmt"
	"github.com/gimigkk/Alfred-Memory/internal/agent"
	"github.com/gimigkk/Alfred-Memory/internal/ladybug"
)

func main() {
	db := &ladybug.Database{}
	conn, _ := ladybug.NewConnection(db)
	
	orch := &agent.Orchestrator{
		DBConn: conn,
	}

	createPayload := `{"mutations": [{"operation": "CREATE_NODE", "node_type": "Task", "node_id": "temp_task_alprog", "properties": {"title": "Kuis Alprog", "content": "Kuis"}}], "thought": "creating"}`
	
	fmt.Println("=== CREATING NODE ===")
	res, err := orch.ExecuteChatTool("commit_chat_mutations", createPayload)
	fmt.Printf("Result: %s, Error: %v\n", res, err)

	updatePayload := `{"mutations": [{"operation": "UPDATE_NODE", "node_type": "Task", "node_id": "temp_task_alprog", "properties": {"title": "Kuis Alprog 2", "content": "Kuis 2"}}], "thought": "updating"}`
	
	fmt.Println("\n=== UPDATING NODE (USING TEMP) ===")
	res, err = orch.ExecuteChatTool("commit_chat_mutations", updatePayload)
	fmt.Printf("Result: %s, Error: %v\n", res, err)
}
