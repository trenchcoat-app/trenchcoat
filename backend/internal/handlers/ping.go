package handlers

import (
	"divvy/internal/api"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, api.PingResponse{
		Message: "pong",
	})
}
