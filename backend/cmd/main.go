package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"divvy/config"
	"divvy/internal/db"
)

func main() {
	if err := config.Init(); err != nil {
		log.Fatalf("config init failed: %v", err)
	}

	if err := db.RunMigrations(); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	database := db.OpenDB()
	defer database.Close()

	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	if err := router.Run(); err != nil {
		log.Fatal(err)
	}
}
