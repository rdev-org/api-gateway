package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/rdev-org/api-gateway/internal/notification"
	"github.com/rdev-org/api-gateway/internal/user_address"
	"github.com/rdev-org/api-gateway/internal/user_credential"
	"github.com/rdev-org/api-gateway/internal/user_preference"
	"github.com/rdev-org/api-gateway/internal/user_profile"

	pbAddress "github.com/rdev-org/proto-contracts/pb/user_address"
	pbCredential "github.com/rdev-org/proto-contracts/pb/user_credential"
	pbPreference "github.com/rdev-org/proto-contracts/pb/user_preference"
	pbProfile "github.com/rdev-org/proto-contracts/pb/user_profile"
)

type Router struct {
	Credential          pbCredential.UserCredentialServiceClient
	Profile             pbProfile.UserProfileServiceClient
	Address             pbAddress.UserAddressServiceClient
	Preference          pbPreference.UserPreferenceServiceClient
	NotificationHandler *notification.NotificationHandler
}

func NewRouter(r Router) *Router {
	return &Router{
		Credential:          r.Credential,
		Profile:             r.Profile,
		Address:             r.Address,
		Preference:          r.Preference,
		NotificationHandler: r.NotificationHandler,
	}
}

func (rt *Router) SetAllRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/health"))

	r.Route("/api/v1", func(subRouter chi.Router) {
		rt.SetUserCredentialRoutes(subRouter)
		rt.SetUserProfileRoutes(subRouter)
		rt.SetUserAddressRoutes(subRouter)
		rt.SetUserPreferenceRoutes(subRouter)
		rt.SetNotificationRoute(subRouter)
	})

	return r
}

func (rt *Router) SetNotificationRoute(r chi.Router) {
	if rt.NotificationHandler != nil {
		r.Get("/ws/notifications", rt.NotificationHandler.InAppNotification)
	}
}

func (rt *Router) SetUserCredentialRoutes(r chi.Router) {
	h := user_credential.NewUserCredentialHandler(rt.Credential)

	r.Post("/login", h.Login)
	r.Post("/register", h.Register)
	r.Post("/forget_password", h.ForgotPassword)

	r.Group(func(subRouter chi.Router) {
		subRouter.Delete("/{uuid}/delete", h.DeleteCredential)
		subRouter.Post("/update_password", h.UpdatePassword)
		subRouter.Post("/reset_password", h.ResetPassword)
	})
}

func (rt *Router) SetUserProfileRoutes(r chi.Router) {
	h := user_profile.NewUserProfileHandler(rt.Profile)

	r.Get("/profile", h.GetProfile)

	r.Route("/user_profile", func(subRouter chi.Router) {
		subRouter.Get("/list", h.ListProfile)
		subRouter.Post("/create", h.CreateProfile)
		subRouter.Put("/{uuid}/update", h.UpdateProfile)
		subRouter.Delete("/{uuid}/delete", h.DeleteProfile)
	})
}

func (rt *Router) SetUserAddressRoutes(r chi.Router) {
	h := user_address.NewUserAddressHandler(rt.Address)

	r.Route("/user_address", func(subRouter chi.Router) {
		subRouter.Get("/{uuid}", h.GetAddress)
		subRouter.Get("/list", h.ListAddress)
		subRouter.Post("/create", h.CreateAddress)
		subRouter.Put("/{uuid}/update", h.UpdateAddress)
		subRouter.Delete("/{uuid}/delete", h.DeleteAddress)
		subRouter.Get("/{uuid}/default", h.DefaultAddress)
	})
}

func (rt *Router) SetUserPreferenceRoutes(r chi.Router) {
	h := user_preference.NewUserPreferenceHandler(rt.Preference)

	r.Route("/user_preference", func(subRouter chi.Router) {
		subRouter.Get("/{uuid}", h.GetPreference)
		subRouter.Get("/list", h.ListPreference)
		subRouter.Post("/create", h.CreatePreference)
		subRouter.Put("/{uuid}/update", h.UpdatePreference)
		subRouter.Delete("/{uuid}/delete", h.DeletePreference)
	})
}
