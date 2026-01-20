package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"example.com/bc-mvp/internal/auth"
	"example.com/bc-mvp/internal/store"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	Users    *store.UserStore
	Stats    *store.StatsStore
	Auth     *auth.Service
	TokenTTL time.Duration
}

type RegisterRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	DisplayName string `json:"displayName"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"accessToken"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "use POST")
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", "invalid json")
		return
	}
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.DisplayName = strings.TrimSpace(req.DisplayName)

	if req.Email == "" || req.Password == "" || req.DisplayName == "" {
		writeError(w, http.StatusBadRequest, "bad_request", "email, password and displayName are required")
		return
	}
	if len(req.Password) < 1 {
		writeError(w, http.StatusBadRequest, "bad_request", "password must be at least 1 chars")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal", "failed to hash password")
		return
	}

	userID := uuid.NewString()
	u := store.User{
		ID:           userID,
		Email:        req.Email,
		PasswordHash: string(hash),
		DisplayName:  req.DisplayName,
	}

	if err := h.Users.Create(r.Context(), u); err != nil {
		if err == store.ErrEmailTaken {
			writeError(w, http.StatusConflict, "email_taken", "email already exists")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal", "failed to create user")
		return
	}

	// создаём пустую статистику
	_ = h.Stats.InitForUser(r.Context(), userID)

	// MVP: register возвращает 201 без токена (можно сделать сразу login, если хочешь)
	w.WriteHeader(http.StatusCreated)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "use POST")
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", "invalid json")
		return
	}
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	if req.Email == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "bad_request", "email and password are required")
		return
	}

	u, err := h.Users.GetByEmail(r.Context(), req.Email)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "invalid_credentials", "invalid email or password")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password)); err != nil {
		writeError(w, http.StatusUnauthorized, "invalid_credentials", "invalid email or password")
		return
	}

	token, err := h.Auth.SignWithName(u.ID, u.DisplayName, h.TokenTTL)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal", "failed to sign token")
		return
	}

	writeJSON(w, http.StatusOK, LoginResponse{AccessToken: token})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserIDFromContext(r.Context())
	if !ok || userID == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized", "missing auth context")
		return
	}

	u, err := h.Users.GetByID(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "unauthorized", "user not found")
		return
	}

	st, err := h.Stats.Get(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal", "failed to load stats")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"id":          u.ID,
		"email":       u.Email,
		"displayName": u.DisplayName,
		"createdAt":   u.CreatedAt,
		"stats": map[string]any{
			"wins":   st.Wins,
			"losses": st.Losses,
			"draws":  st.Draws,
		},
	})
}
