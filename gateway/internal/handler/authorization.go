package handler

import (
	"encoding/json"
	"net/http"

	"github.com/notenoughtea/currency_review/gateway/internal/dto"
	"github.com/notenoughtea/currency_review/gateway/internal/logger"
)

type registerRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (s *controller) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	if err := s.authService.Register(dto.RegisterRequest(req)); err != nil {
		logger.Log.Errorf("register error: %v", err)
		s.handleError(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (s *controller) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	token, err := s.authService.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		logger.Log.Errorf("login error: %v", err)
		s.handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"token": token})
}

func (s *controller) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	token := r.Header.Get("Authorization")
	if token == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Authorization token is required"})
		return
	}
	if err := s.authService.Logout(token); err != nil {
		logger.Log.Errorf("logout error: %v", err)
		s.handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "logout successful"})
}

func (s *controller) handleError(w http.ResponseWriter, err error) {
	writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
}
