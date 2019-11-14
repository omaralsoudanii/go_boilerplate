package routes

import (
	"go_boilerplate/item"
	_itemHttpDelivery "go_boilerplate/item/delivery/http"
	"go_boilerplate/middleware"

	"github.com/go-chi/chi"
)

func ItemHttpRouter(router *chi.Mux, UseCase item.UseCase) {
	handler := &_itemHttpDelivery.NewHttpItemHandler{
		ItemUseCase: UseCase,
	}
	r := chi.NewRouter()
	r.Get("/{id}", handler.GetByID)
	r.Group(func(r chi.Router) {
		r.Use(middleware.AccessTokenVerifier)
		r.Delete("/{id}", handler.Delete)
		r.Post("/", handler.Store)
	})
	r.Get("/", handler.GetAllItem)

	router.Mount("/items", r)
}
