package main

import (
	"fmt"
	"github.com/gimigkk/Alfred-Memory/internal/ladybug"
)

func main() {
	db, err := ladybug.NewDatabase(".lbug")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer db.Close()

	conn, _ := ladybug.NewConnection(db)
	defer conn.Close()

	res, err := conn.Query("MATCH (n:Event) RETURN n.id, n.needs_clarification")
	if err != nil {
		fmt.Println("Query Error:", err)
		return
	}
	for res.HasNext() {
		fmt.Println(res.GetNext())
	}
}
