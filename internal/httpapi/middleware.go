package httpapi

import (
	"context"
	"net/http"
	"strings"

	"example.com/bc-mvp/internal/auth"
)

type ctxKey string

const userIDKey ctxKey = "userID"

func AuthMiddleware(jwtSecret []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h := r.Header.Get("Authorization")
			if !strings.HasPrefix(h, "Bearer ") {
				writeError(w, http.StatusUnauthorized, "unauthorized", "missing bearer token")
				return
			}
			token := strings.TrimPrefix(h, "Bearer ")

			claims, err := auth.Verify(jwtSecret, token)
			if err != nil {
				writeError(w, http.StatusUnauthorized, "unauthorized", "invalid token")
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserIDFromContext(ctx context.Context) (string, bool) {
	v := ctx.Value(userIDKey)
	s, ok := v.(string)
	return s, ok
}
