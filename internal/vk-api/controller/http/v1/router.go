package v1

import (
	"encoding/json"
	"github.com/Nimartemoff/vk-api/cmd/vk-api/config"
	"github.com/Nimartemoff/vk-api/internal/vk-api/usecase"
	"github.com/go-chi/chi"
	"net/http"
)

func NewRouter(cfg *config.Config, r *chi.Mux, uc *usecase.UserUsecase) {
	r.Route("/api/v1", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			newUserRoutes(r, uc)
		})
	})
}

func renderJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		renderError(w, http.StatusInternalServerError, err)
	}
}

type jsonError struct {
	Msg string `json:"error"`
}

func renderError(w http.ResponseWriter, statusCode int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(jsonError{Msg: err.Error()})
}
