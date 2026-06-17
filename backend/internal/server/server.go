package server

import (
	"divvy/internal/api"
	"divvy/internal/handlers"
	"log"

	"github.com/gin-gonic/gin"
)

func GetRouter() *gin.Engine {
	router := gin.Default()

	srv := handlers.NewServer()
	api.RegisterHandlers(router, srv)

	return router
}

func Run() {
	router := GetRouter()

	if err := router.Run(); err != nil {
		log.Fatal(err)
	}
}
