package handlers

import (
	"trenchcoat/internal/services"
)

// This struct's main purpose is to implement api.ServerInterface from
// oapi-codegen's auto-generated code based on the OpenAPI schema.
type Server struct {
	AuthService *services.AuthService
}

func NewServer(authService *services.AuthService) *Server {
	return &Server{AuthService: authService}
}
