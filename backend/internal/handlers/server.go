package handlers

import (
	"trenchcoat/internal/services/auth"
)

// This struct's main purpose is to implement api.ServerInterface from
// oapi-codegen's auto-generated code based on the OpenAPI schema.
type Server struct {
	AuthService *auth.AuthService
}

func NewServer(authService *auth.AuthService) *Server {
	return &Server{AuthService: authService}
}
