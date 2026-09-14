package container

import (
	"log"

	"github.com/go-chi/chi/v5"
	"github.com/rdev-org/api-gateway/internal/config"
	"github.com/rdev-org/api-gateway/internal/router"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pbUser "github.com/rdev-org/proto-contracts/pb/user"
)

type AppContainer struct {
	UserConn *grpc.ClientConn
}

func NewContainer() (*AppContainer, *chi.Mux) {
	appConfig, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("failed to load config %v", err)
	}

	userServiceAddr := appConfig.Services.User
	userConn, err := grpc.NewClient(userServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to user service: %v", err)
	}

	userClient := pbUser.NewUserServiceClient(userConn)
	r := router.NewRouter(userClient)

	appContainer := &AppContainer{
		UserConn: userConn,
	}

	return appContainer, r.SetAllRoutes()
}
