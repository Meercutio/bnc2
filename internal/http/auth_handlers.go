package httpapi

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"example.com/bc-mvp/internal/auth"
	"example.com/bc-mvp/internal/store"
)

type AuthHandler struct {
	Users     *store.UserStore
	JWTSecret []byte
	TokenTTL  time.Duration
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email       string `json:"email"`
		Password    string `json:"password"`
		DisplayName string `json:"displayName"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)

	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	u := store.User{
		ID:           uuid.NewString(),
		Email:        req.Email,
		PasswordHash: string(hash),
		DisplayName:  req.DisplayName,
	}

	if err := h.Users.Create(r.Context(), u); err != nil {
		http.Error(w, "email already exists", http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)

	u, err := h.Users.GetByEmail(r.Context(), req.Email)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password)) != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	token, _ := auth.Sign(h.JWTSecret, u.ID, h.TokenTTL)

	json.NewEncoder(w).Encode(map[string]string{
		"accessToken": token,
	})
}
