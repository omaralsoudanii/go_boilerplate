package routes

import (
	"go_boilerplate/middleware"
	"go_boilerplate/user"
	_userHttpDelivery "go_boilerplate/user/delivery/http"

	"github.com/go-chi/chi"
)

func UserHttpRouter(router *chi.Mux, UseCase user.UseCase) {
	handler := &_userHttpDelivery.NewHttpUserHandler{
		UserUseCase: UseCase,
	}
	r := chi.NewRouter()
	r.Post("/register", handler.Register)
	r.Post("/login", handler.SignIn)
	r.Group(func(r chi.Router) {
		r.Use(middleware.RefreshTokenVerifier)
		r.Post("/refresh", handler.Refresh)
	})
	r.Group(func(r chi.Router) {
		r.Use(middleware.AccessTokenVerifier)
		r.Post("/logout", handler.SignOut)
	})
	router.Mount("/auth", r)
}
