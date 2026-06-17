package handlers

// This struct's main purpose is to implement api.ServerInterface from
// oapi-codegen's auto-generated code based on the OpenAPI schema.
type Server struct{}

func NewServer() *Server {
	return &Server{}
}
