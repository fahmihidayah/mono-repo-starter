package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/fahmihidayah/go-api-orchestrator/internal/config"
	"github.com/fahmihidayah/go-api-orchestrator/internal/db"
	"github.com/fahmihidayah/go-api-orchestrator/internal/db/migrations"
)

func main() {
	// Define command line flags
	action := flag.String("action", "up", "Migration action: up, down, fresh, seed")
	flag.Parse()

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	database, err := db.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB, err := database.DB()
	if err != nil {
		log.Fatalf("Failed to get database connection: %v", err)
	}
	defer sqlDB.Close()

	// Execute migration based on action
	switch *action {
	case "up":
		fmt.Println("Running migrations...")
		if err := migrations.AutoMigrate(database); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
		fmt.Println("✓ Migration completed successfully")

	case "down":
		fmt.Println("Rolling back migrations...")
		if err := migrations.DropAllTables(database); err != nil {
			log.Fatalf("Rollback failed: %v", err)
		}
		fmt.Println("✓ Rollback completed successfully")

	case "fresh":
		fmt.Println("Running fresh migration (drop + migrate)...")
		if err := migrations.DropAllTables(database); err != nil {
			log.Fatalf("Drop tables failed: %v", err)
		}
		if err := migrations.AutoMigrate(database); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
		if err := migrations.CreateIndexes(database); err != nil {
			log.Fatalf("Create indexes failed: %v", err)
		}
		fmt.Println("✓ Fresh migration completed successfully")

	case "seed":
		fmt.Println("Seeding database...")
		if err := migrations.SeedData(database); err != nil {
			log.Fatalf("Seeding failed: %v", err)
		}
		fmt.Println("✓ Seeding completed successfully")

	case "status":
		fmt.Println("Checking migration status...")
		// Check if tables exist
		tables := []string{"users", "accounts", "user_sessions"}
		for _, table := range tables {
			if database.Migrator().HasTable(table) {
				fmt.Printf("✓ Table '%s' exists\n", table)
			} else {
				fmt.Printf("✗ Table '%s' does not exist\n", table)
			}
		}

	default:
		fmt.Printf("Unknown action: %s\n", *action)
		fmt.Println("Available actions: up, down, fresh, seed, status")
		os.Exit(1)
	}
}
