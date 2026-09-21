package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/rdev-org/api-gateway/internal/handler/auth"
	"github.com/rdev-org/api-gateway/internal/handler/user"
	"github.com/rdev-org/api-gateway/internal/notification"

	pbAuth "github.com/rdev-org/proto-contracts/pb/auth"
	pb "github.com/rdev-org/proto-contracts/pb/user"
)

type Router struct {
	User                pb.UserServiceClient
	Auth                pbAuth.AuthenticationServiceClient
	NotificationHandler *notification.NotificationHandler
}

func NewRouter(r Router) *Router {
	return &Router{
		User:                r.User,
		Auth:                r.Auth,
		NotificationHandler: r.NotificationHandler,
	}
}

func (rt *Router) SetAllRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/health"))

	r.Route("/api/v1", func(subRouter chi.Router) {
		rt.RegisterAuthenticationRoute(subRouter)
		rt.SetNotificationRoute(subRouter)

		rt.RegisterUserRoutes(subRouter)
	})

	return r
}

func (rt *Router) SetNotificationRoute(r chi.Router) {
	if rt.NotificationHandler != nil {
		r.Get("/ws/notifications", rt.NotificationHandler.InAppNotification)
	}
}

func (rt *Router) RegisterAuthenticationRoute(r chi.Router) {
	h := auth.NewAuthenticationHandler(rt.Auth)

	r.Post("/login", h.Login)
	r.Post("/register", h.Register)
	r.Post("/forget_password", h.ForgotPassword)

	r.Group(func(subRouter chi.Router) {
		subRouter.Delete("/{uuid}/delete", h.DeleteCredential)
		subRouter.Post("/update_password", h.UpdatePassword)
		subRouter.Post("/reset_password", h.ResetPassword)
	})
}

func (rt *Router) RegisterUserRoutes(r chi.Router) {
	h := user.NewUserService(rt.User)

	rt.SetUserProfileRoutes(h, r)
	rt.SetUserAddressRoutes(h, r)
	rt.SetUserPreferenceRoutes(h, r)
}

func (rt *Router) SetUserProfileRoutes(h *user.UserHandler, r chi.Router) {
	r.Get("/profile", h.MyProfile)

	r.Route("/user_profile", func(subRouter chi.Router) {
		subRouter.Get("/{uuid}/profile", h.GetProfile)
		subRouter.Get("/list", h.ListProfile)
		subRouter.Post("/create", h.CreateProfile)
		subRouter.Put("/{uuid}/update", h.UpdateProfile)
		subRouter.Delete("/{uuid}/delete", h.DeleteProfile)
	})
}

func (rt *Router) SetUserAddressRoutes(h *user.UserHandler, r chi.Router) {
	r.Route("/user_address", func(subRouter chi.Router) {
		subRouter.Get("/{uuid}", h.GetAddress)
		subRouter.Get("/list", h.ListAddress)
		subRouter.Post("/create", h.CreateAddress)
		subRouter.Put("/{uuid}/update", h.UpdateAddress)
		subRouter.Delete("/{uuid}/delete", h.DeleteAddress)
		subRouter.Get("/{uuid}/default", h.DefaultAddress)
	})
}

func (rt *Router) SetUserPreferenceRoutes(h *user.UserHandler, r chi.Router) {
	r.Route("/user_preference", func(subRouter chi.Router) {
		subRouter.Get("/{uuid}", h.GetPreference)
		subRouter.Get("/list", h.ListPreference)
		subRouter.Post("/create", h.CreatePreference)
		subRouter.Put("/{uuid}/update", h.UpdatePreference)
		subRouter.Delete("/{uuid}/delete", h.DeletePreference)
	})
}
