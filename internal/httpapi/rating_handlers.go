package httpapi

import (
	"net/http"
	"strconv"

	"example.com/bc-mvp/internal/store"
)

type RatingHandler struct {
	Stats *store.StatsStore
}

func (h *RatingHandler) Leaderboard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "use GET")
		return
	}
	limit := 50
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			limit = n
		}
	}

	items, err := h.Stats.Leaderboard(r.Context(), limit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal", "failed to load leaderboard")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

func (h *RatingHandler) Me(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "use GET")
		return
	}

	userID, ok := UserIDFromContext(r.Context())
	if !ok || userID == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized", "missing auth context")
		return
	}

	me, err := h.Stats.MyRating(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal", "failed to load rating")
		return
	}
	writeJSON(w, http.StatusOK, me)
}
