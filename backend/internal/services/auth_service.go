package services

import "github.com/jackc/pgx/v5/pgxpool"

type AuthService struct {
	DB *pgxpool.Pool
}

func NewAuthService(db *pgxpool.Pool) *AuthService {
	return &AuthService{DB: db}
}
