package server

import (
	"divvy/config"
	"divvy/internal/api"
	"divvy/internal/handlers"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func GetRouter() *gin.Engine {
	router := gin.Default()

	corsConfig := cors.Config{
		AllowOrigins:     config.AppConfig.CORS_ALLOWED_ORIGINS,
		AllowCredentials: true,
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
	}
	router.Use(cors.New(corsConfig))

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
