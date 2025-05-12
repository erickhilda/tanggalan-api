package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

func main() {
	dbPath := "data/tanggalan.db"
	schemaPath := "internal/database/schema.sql"

	schema, err := os.ReadFile(schemaPath)
	if err != nil {
		fmt.Println("❌ Failed to read schema.sql:", err)
		os.Exit(1)
	}

	err = os.MkdirAll("data", 0755)
	if err != nil {
		fmt.Println("❌ Failed to create data folder:", err)
		os.Exit(1)
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		fmt.Println("❌ Failed to open DB:", err)
		os.Exit(1)
	}
	defer db.Close()

	_, err = db.Exec(string(schema))
	if err != nil {
		fmt.Println("❌ Failed to execute schema:", err)
		os.Exit(1)
	}

	fmt.Println("✅ Database initialized at:", dbPath)
}
