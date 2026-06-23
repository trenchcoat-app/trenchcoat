package handlers

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

// This struct's main purpose is to implement api.ServerInterface from
// oapi-codegen's auto-generated code based on the OpenAPI schema.
type Server struct {
	DB *pgxpool.Pool
}

func NewServer(db *pgxpool.Pool) *Server {
	return &Server{DB: db}
}
