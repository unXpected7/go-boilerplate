package service

import (
	"github.com/sriniously/go-boilerplate/apps/backend/internal/server"

	"github.com/clerk/clerk-sdk-go/v2"
)

type AuthService struct {
	server *server.Server
}

func NewAuthService(s *server.Server) *AuthService {
	clerk.SetKey(s.Config.AuthSecretKey)
	return &AuthService{
		server: s,
	}
}
