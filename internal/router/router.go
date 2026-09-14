package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/rdev-org/api-gateway/internal/users"
	pbUser "github.com/rdev-org/proto-contracts/pb/user"
)

type Router struct {
	UserClient pbUser.UserServiceClient
}

func NewRouter(userClient pbUser.UserServiceClient) *Router {
	return &Router{
		UserClient: userClient,
	}
}

func (rt *Router) SetAllRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/health"))

	r.Route("/api/v1", func(subRouter chi.Router) {
		rt.SetUserRoutes(subRouter)
	})

	return r
}

func (rt *Router) SetUserRoutes(r chi.Router) {
	userHandler := users.NewUserHandler(rt.UserClient)

	r.Group(func(r chi.Router) {
		r.Post("/login", userHandler.Login)
		r.Post("/register", userHandler.Register)
	})
}
