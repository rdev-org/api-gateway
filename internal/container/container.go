package container

import (
	"log"

	"github.com/go-chi/chi/v5"
	"github.com/rdev-org/api-gateway/internal/config"
	"github.com/rdev-org/api-gateway/internal/router"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pbUser "github.com/rdev-org/proto-contracts/pb/user"
	pbUserAddress "github.com/rdev-org/proto-contracts/pb/user_address"
)

type AppContainer struct {
	Routes   *chi.Mux
	UserConn *grpc.ClientConn
}

func NewContainer(servicesURL config.ServicesURL) *AppContainer {
	userServiceAddr := servicesURL.User
	userConn, err := grpc.NewClient(userServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to user service: %v", err)
	}

	userClient := pbUser.NewUserServiceClient(userConn)
	userAddressClient := pbUserAddress.NewUserAddressServiceClient(userConn)

	r := router.NewRouter(
		userClient,
		userAddressClient,
	)

	return &AppContainer{
		Routes:   r.SetAllRoutes(),
		UserConn: userConn,
	}
}
