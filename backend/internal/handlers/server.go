package handlers

import (
	"trenchcoat/internal/interfaces"
)

// This struct's main purpose is to implement api.ServerInterface from
// oapi-codegen's auto-generated code based on the OpenAPI schema.
type Server struct {
	AuthService interfaces.AuthServiceInterface
}

func NewServer(authService interfaces.AuthServiceInterface) *Server {
	return &Server{AuthService: authService}
}
