package main

import (
	"log"

	"divvy/config"
	"divvy/internal/db"
	"divvy/internal/server"
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

	server.Run()
}
