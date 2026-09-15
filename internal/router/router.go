package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/rdev-org/api-gateway/internal/user_addresses"
	"github.com/rdev-org/api-gateway/internal/users"
	pbUser "github.com/rdev-org/proto-contracts/pb/user"
	pbUserAddress "github.com/rdev-org/proto-contracts/pb/user_address"
)

type Router struct {
	UserClient        pbUser.UserServiceClient
	UserAddressClient pbUserAddress.UserAddressServiceClient
}

func NewRouter(userClient pbUser.UserServiceClient, userAddressClient pbUserAddress.UserAddressServiceClient) *Router {
	return &Router{
		UserClient:        userClient,
		UserAddressClient: userAddressClient,
	}
}

func (rt *Router) SetAllRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/health"))

	r.Route("/api/v1", func(subRouter chi.Router) {
		rt.SetUserRoutes(subRouter)
		rt.SetUserAddressRoutes(subRouter)
	})

	return r
}

func (rt *Router) SetUserRoutes(r chi.Router) {
	h := users.NewUserHandler(rt.UserClient)

	r.Group(func(r chi.Router) {
		r.Post("/login", h.Login)
		r.Post("/register", h.Register)
	})
}

func (rt *Router) SetUserAddressRoutes(r chi.Router) {
	h := user_addresses.NewUserAddressHandler(rt.UserAddressClient)

	r.Route("/user_address", func(grpRoute chi.Router) {
		// Authentication
		grpRoute.Get("/{uuid}", h.GetAddress)
		grpRoute.Get("/", h.GetAllAddress)
		grpRoute.Post("/", h.CreateAddress)
		grpRoute.Put("/{uuid}", h.UpdateAddress)
		grpRoute.Delete("/{uuid}", h.DeleteAddress)
		grpRoute.Get("/{uuid}/default", h.SetDefaultAddress)
	})
}
