package middleware

import (
	"net/http"
	"strings"

	"github.com/notenoughtea/currency_review/gateway/internal/logger"
	"github.com/notenoughtea/currency_review/gateway/internal/service"
)

type AuthMiddleware struct {
	authService service.AuthServiceInterface
}

func NewAuthMiddleware(authService service.AuthServiceInterface) *AuthMiddleware {
	return &AuthMiddleware{authService: authService}
}

func (m *AuthMiddleware) Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Пропускаем авторизацию для определенных маршрутов
		if m.shouldSkipAuth(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusBadRequest)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid Authorization header format", http.StatusBadRequest)
			return
		}

		token := parts[1]
		err := m.authService.ValidateToken(r.Context(), token)
		if err != nil {
			logger.Log.Errorf("invalid token: %v", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (m *AuthMiddleware) shouldSkipAuth(path string) bool {
	// Маршруты, которые не требуют авторизации
	skipPaths := []string{
		"/ping",
		"/api/v1/register",
		"/api/v1/login",
		"/api/v1/metrics",
	}

	for _, skipPath := range skipPaths {
		if path == skipPath {
			return true
		}
	}
	return false
}
