package main

import (
	"log"

	"planning-trip-be/internal/config"
	"planning-trip-be/internal/database"
)

func main() {
	env, err := config.Load()
	if err != nil {
		log.Fatalf("config load failed: %v", err)
	}

	db, err := database.OpenPostgres(env.DatabaseURL)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}

	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	var currentDB string
	var currentSchema string
	if err := db.Raw("select current_database(), current_schema()").Row().Scan(&currentDB, &currentSchema); err != nil {
		log.Fatalf("cannot read database context: %v", err)
	}

	var tableCount int64
	if err := db.Raw("select count(*) from information_schema.tables where table_schema = 'public'").Scan(&tableCount).Error; err != nil {
		log.Fatalf("cannot count tables: %v", err)
	}

	log.Printf("database=%s schema=%s public_tables=%d", currentDB, currentSchema, tableCount)
	log.Println("migration completed successfully")
}
