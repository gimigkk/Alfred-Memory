package db

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	ladybug "github.com/gimigkk/Alfred-Memory/internal/ladybug"
	_ "github.com/mattn/go-sqlite3"
)

// InitLadybugSchema runs all the DDL commands to create nodes and edges
func InitLadybugSchema(conn *ladybug.Connection) error {
	queries := []string{
		// Nodes
		"CREATE NODE TABLE Person (id STRING, name STRING, phone_number STRING, aliases STRING[], created_at TIMESTAMP, needs_clarification BOOLEAN, PRIMARY KEY(id))",
		"CREATE NODE TABLE Circle (id STRING, name STRING, aliases STRING[], title STRING, content STRING, verbatim STRING, history STRING[], created_at TIMESTAMP, needs_clarification BOOLEAN, embedding FLOAT[768], PRIMARY KEY(id))",
		"CREATE NODE TABLE Task (id STRING, title STRING, content STRING, aliases STRING[], verbatim STRING, status STRING, due_date TIMESTAMP, history STRING[], created_at TIMESTAMP, needs_clarification BOOLEAN, group_mentions STRING, embedding FLOAT[768], PRIMARY KEY(id))",
		"CREATE NODE TABLE Event (id STRING, title STRING, content STRING, aliases STRING[], verbatim STRING, status STRING, start_date TIMESTAMP, history STRING[], created_at TIMESTAMP, needs_clarification BOOLEAN, group_mentions STRING, embedding FLOAT[768], PRIMARY KEY(id))",
		"CREATE NODE TABLE Insight (id STRING, title STRING, content STRING, aliases STRING[], verbatim STRING, history STRING[], created_at TIMESTAMP, needs_clarification BOOLEAN, embedding FLOAT[768], PRIMARY KEY(id))",

		// Edges
		"CREATE REL TABLE ASSIGNED_TO (FROM Person TO Task, evidence_refs STRING)",
		"CREATE REL TABLE MENTIONED_IN (FROM Person TO Task, FROM Person TO Event, evidence_refs STRING)",
		"CREATE REL TABLE HAS_ROLE (FROM Person TO Event, evidence_refs STRING)",
		"CREATE REL TABLE MEMBER_OF (FROM Person TO Circle, role STRING, since TIMESTAMP, evidence_refs STRING)",
		"CREATE REL TABLE PART_OF (FROM Task TO Event, FROM Task TO Circle, evidence_refs STRING)",
		"CREATE REL TABLE DIR_TOWARDS (FROM Insight TO Person, FROM Insight TO Circle, evidence_refs STRING)",
		"CREATE REL TABLE LINKS_TO (FROM Task TO Event, FROM Task TO Insight, FROM Task TO Task, FROM Event TO Insight, FROM Event TO Event, FROM Insight TO Insight, context STRING, evidence_refs STRING)",
		"CREATE REL TABLE CONTRADICTS (FROM Insight TO Insight, detected_at TIMESTAMP, resolved BOOLEAN)",
		"CREATE REL TABLE KNOWS (FROM Person TO Person, descriptor STRING, context STRING, since TIMESTAMP, evidence_refs STRING)",

		// Vector Indexes
		"CREATE VECTOR INDEX circle_vec_idx ON Circle(embedding)",
		"CREATE VECTOR INDEX task_vec_idx ON Task(embedding)",
		"CREATE VECTOR INDEX event_vec_idx ON Event(embedding)",
		"CREATE VECTOR INDEX insight_vec_idx ON Insight(embedding)",
	}

	for _, q := range queries {
		res, err := conn.Query(q)
		if err != nil && !strings.Contains(err.Error(), "already exists") {
			// Ignore already exists, fail on actual errors
			log.Printf("DDL Warning: %v\n", err)
		}
		if res != nil {
			res.Close()
		}
	}

	log.Println("LadybugDB Schema Initialization complete.")
	return nil
}

// InitSQLite initializes the reminders.db SQLite database
func InitSQLite(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite database: %w", err)
	}

	query := `
	CREATE TABLE IF NOT EXISTS Reminders (
		id TEXT PRIMARY KEY,
		node_id TEXT,
		deadline DATETIME,
		is_sent BOOLEAN,
		message TEXT
	);
	`
	_, err = db.Exec(query)
	if err != nil {
		return nil, fmt.Errorf("failed to create Reminders table: %w", err)
	}

	log.Println("SQLite Schema Initialization complete.")
	return db, nil
}
