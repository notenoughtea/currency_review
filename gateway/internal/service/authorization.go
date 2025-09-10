package service

import (
	"context"
	"errors"

	authClient "github.com/notenoughtea/currency_review/gateway/internal/clients/auth"
	"github.com/notenoughtea/currency_review/gateway/internal/dto"
)

var ErrInvalidToken = errors.New("invalid token")

type AuthServiceInterface interface {
	GenerateToken(ctx context.Context, login string) (string, error)
	ValidateToken(ctx context.Context, token string) error
	Register(req dto.RegisterRequest) error
	Login(ctx context.Context, username, password string) (string, error)
	Logout(token string) error
}

type AuthService struct {
	authClient *authClient.Client
}

func NewAuthService(authClient *authClient.Client) *AuthService {
	return &AuthService{authClient: authClient}
}

func (s *AuthService) GenerateToken(ctx context.Context, login string) (string, error) {
	return s.authClient.GenerateToken(ctx, login)
}

func (s *AuthService) ValidateToken(ctx context.Context, token string) error {
	return s.authClient.ValidateToken(ctx, token)
}

func (s *AuthService) Register(req dto.RegisterRequest) error {
	return nil
}

func (s *AuthService) Login(ctx context.Context, username, password string) (string, error) {
	return s.authClient.GenerateToken(ctx, username)
}

func (s *AuthService) Logout(token string) error {
	return nil
}
