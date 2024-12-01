package v1

import (
	"github.com/Nimartemoff/vk-api/internal/vk-api/usecase"
	"github.com/go-chi/chi"
)

type userRoutes struct {
	*usecase.UserUsecase
}

func newUserRoutes(r chi.Router, uc *usecase.UserUsecase) {
	ur := &userRoutes{uc}

	r.Get("/nodes", ur.getAllNodes)
	r.Get("/nodes/{id}", ur.getNode)

	r.With(userHasAnyRoleMiddleware("editor")).Group(func(r chi.Router) {
		r.Post("/nodes", ur.createNode)
		r.Delete("/nodes/{id}", ur.deleteNode)
	})
}
