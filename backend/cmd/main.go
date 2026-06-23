package main

import (
	"log"

	"trenchcoat/config"
	"trenchcoat/internal/db"
	"trenchcoat/internal/server"
)

func main() {
	parseFlags()

	if err := config.Init(); err != nil {
		log.Fatalf("config init failed: %v", err)
	}

	if !skipMigrations {
		if err := db.RunMigrations(); err != nil {
			log.Fatalf("migration failed: %v", err)
		}
	} else {
		log.Println("skipping database migrations")
	}

	database := db.OpenDB()
	defer database.Close()

	server.Run(database)
}
