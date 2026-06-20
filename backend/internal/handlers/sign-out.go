package handlers

import (
	"divvy/internal/api"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) SignOut(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, api.ErrorResponse{
		Code:    "NOT_YET_IMPLEMENTED",
		Message: "This endpoint is not yet implemented.",
	})
}
