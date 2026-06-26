package handlers

// This struct's main purpose is to implement api.ServerInterface from
// oapi-codegen's auto-generated code based on the OpenAPI schema.
type Server struct {
	AuthService AuthServiceInterface
}

func NewServer(authService AuthServiceInterface) *Server {
	return &Server{AuthService: authService}
}
