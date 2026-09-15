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

type UserServiceContainer struct {
	conn        *grpc.ClientConn
	user        pbUser.UserServiceClient
	userAddress pbUserAddress.UserAddressServiceClient
}

type AppContainer struct {
	Routes   *chi.Mux
	UserConn *grpc.ClientConn
}

func NewContainer(servicesURL config.ServicesURL) *AppContainer {
	userClient, err := setUserService(servicesURL.User)
	if err != nil {
		log.Fatalf("failed to set user client %v", err)
	}

	r := router.NewRouter(
		userClient.user,
		userClient.userAddress,
	)

	return &AppContainer{
		Routes:   r.SetAllRoutes(),
		UserConn: userClient.conn,
	}
}

func setUserService(url string) (UserServiceContainer, error) {
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to user service: %v", err)
	}

	userClient := pbUser.NewUserServiceClient(conn)
	userAddressClient := pbUserAddress.NewUserAddressServiceClient(conn)

	return UserServiceContainer{
		conn:        conn,
		user:        userClient,
		userAddress: userAddressClient,
	}, err
}
