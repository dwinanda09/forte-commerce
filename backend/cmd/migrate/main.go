package main

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"

	"github.com/dwinanda09/forte-commerce/internal/config"
	"github.com/dwinanda09/forte-commerce/internal/resource"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to database
	db, err := resource.NewDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Connected to database")

	// Get migration files
	migrationDir := "migrations"
	files, err := ioutil.ReadDir(migrationDir)
	if err != nil {
		log.Fatalf("Failed to read migrations directory: %v", err)
	}

	// Filter and sort SQL files
	var sqlFiles []string
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".sql" {
			sqlFiles = append(sqlFiles, file.Name())
		}
	}
	sort.Strings(sqlFiles)

	if len(sqlFiles) == 0 {
		log.Println("No migration files found")
		return
	}

	// Execute migrations
	for _, filename := range sqlFiles {
		filepath := filepath.Join(migrationDir, filename)

		// Read SQL file
		content, err := ioutil.ReadFile(filepath)
		if err != nil {
			log.Fatalf("Failed to read migration file %s: %v", filename, err)
		}

		// Execute SQL
		_, err = db.Exec(string(content))
		if err != nil {
			log.Fatalf("Failed to execute migration %s: %v", filename, err)
		}

		log.Printf("Executed migration: %s", filename)
	}

	log.Println("All migrations executed successfully")
}
