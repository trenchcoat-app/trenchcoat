package server

import (
	"log"
	"trenchcoat/config"
	"trenchcoat/internal/api"
	"trenchcoat/internal/handlers"
	"trenchcoat/internal/services/auth"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetRouter(db *pgxpool.Pool) *gin.Engine {
	router := gin.Default()

	corsConfig := cors.Config{
		AllowOrigins:     config.AppConfig.CORS_ALLOWED_ORIGINS,
		AllowCredentials: true,
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
	}
	router.Use(cors.New(corsConfig))

	authService := auth.NewAuthService(db)

	srv := handlers.NewServer(authService)
	api.RegisterHandlers(router, srv)

	return router
}

func Run(db *pgxpool.Pool) {
	router := GetRouter(db)

	if err := router.Run(); err != nil {
		log.Fatal(err)
	}
}
