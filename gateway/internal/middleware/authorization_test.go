package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/notenoughtea/currency_review/gateway/internal/dto"
	"github.com/notenoughtea/currency_review/gateway/internal/logger"
	"github.com/notenoughtea/currency_review/gateway/internal/service"
)

type mockAuthService struct {
	validateTokenFunc func(ctx context.Context, token string) error
}

func (m *mockAuthService) GenerateToken(ctx context.Context, login string) (string, error) {
	return "mock-token", nil
}

func (m *mockAuthService) ValidateToken(ctx context.Context, token string) error {
	if m.validateTokenFunc != nil {
		return m.validateTokenFunc(ctx, token)
	}
	return nil
}

func (m *mockAuthService) Register(req dto.RegisterRequest) error {
	return nil
}

func (m *mockAuthService) Login(ctx context.Context, username, password string) (string, error) {
	return "", nil
}

func (m *mockAuthService) Logout(token string) error {
	return nil
}

func TestAuthMiddleware_Authorize(t *testing.T) {
	logger.Init()

	tests := []struct {
		name           string
		path           string
		authHeader     string
		validateFunc   func(ctx context.Context, token string) error
		expectedStatus int
	}{
		{
			name:           "Valid token",
			path:           "/api/v1/rate",
			authHeader:     "Bearer valid-token",
			validateFunc:   func(ctx context.Context, token string) error { return nil },
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid token",
			path:           "/api/v1/rate",
			authHeader:     "Bearer invalid-token",
			validateFunc:   func(ctx context.Context, token string) error { return service.ErrInvalidToken },
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Missing Authorization header",
			path:           "/api/v1/rate",
			authHeader:     "",
			validateFunc:   nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid Authorization header format",
			path:           "/api/v1/rate",
			authHeader:     "InvalidFormat",
			validateFunc:   nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Skip auth for register",
			path:           "/api/v1/register",
			authHeader:     "",
			validateFunc:   nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Skip auth for login",
			path:           "/api/v1/login",
			authHeader:     "",
			validateFunc:   nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Skip auth for ping",
			path:           "/ping",
			authHeader:     "",
			validateFunc:   nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Skip auth for metrics",
			path:           "/api/v1/metrics",
			authHeader:     "",
			validateFunc:   nil,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuth := &mockAuthService{
				validateTokenFunc: tt.validateFunc,
			}

			middleware := NewAuthMiddleware(mockAuth)

			req := httptest.NewRequest("GET", tt.path, nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			rr := httptest.NewRecorder()

			handler := middleware.Authorize(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))

			handler.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}
