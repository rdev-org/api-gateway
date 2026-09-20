package container

import (
	"context"
	"log"

	"github.com/go-chi/chi/v5"
	"github.com/rdev-org/api-gateway/internal/config"
	"github.com/rdev-org/api-gateway/internal/notification"
	"github.com/rdev-org/api-gateway/internal/router"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pbAuth "github.com/rdev-org/proto-contracts/pb/auth"
	pb "github.com/rdev-org/proto-contracts/pb/user"
)

type UserServiceContainer struct {
	conn *grpc.ClientConn
	user pb.UserServiceClient
	auth pbAuth.AuthenticationServiceClient
}

type Container struct {
	Routes              *chi.Mux
	UserConn            *grpc.ClientConn
	NotificationHandler *notification.NotificationHandler
}

func NewContainer(ctx context.Context, servicesURL config.ServicesURL, ck config.KafkaConfig) *Container {
	userClient, err := setUserService(servicesURL.User)
	if err != nil {
		log.Fatalf("failed to set user client %v", err)
	}

	notifHandler := notification.NewNotificationHandler(ctx, &ck)

	r := router.NewRouter(router.Router{
		User:                userClient.user,
		Auth:                userClient.auth,
		NotificationHandler: notifHandler,
	})

	return &Container{
		Routes:              r.SetAllRoutes(),
		UserConn:            userClient.conn,
		NotificationHandler: notifHandler,
	}
}

func (c *Container) Close() {
	if c.UserConn != nil {
		c.UserConn.Close()
	}
	if c.NotificationHandler != nil {
		c.NotificationHandler.Close()
	}
}

func setUserService(url string) (UserServiceContainer, error) {
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to user service: %v", err)
	}

	auth := pbAuth.NewAuthenticationServiceClient(conn)
	user := pb.NewUserServiceClient(conn)

	return UserServiceContainer{
		conn: conn,
		auth: auth,
		user: user,
	}, err
}
