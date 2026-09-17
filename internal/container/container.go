package container

import (
	"log"

	"github.com/go-chi/chi/v5"
	"github.com/rdev-org/api-gateway/internal/config"
	"github.com/rdev-org/api-gateway/internal/router"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pbAddress "github.com/rdev-org/proto-contracts/pb/user_address"
	pbCredential "github.com/rdev-org/proto-contracts/pb/user_credential"
	pbPreference "github.com/rdev-org/proto-contracts/pb/user_preference"
	pbProfile "github.com/rdev-org/proto-contracts/pb/user_profile"
)

type UserServiceContainer struct {
	conn       *grpc.ClientConn
	credential pbCredential.UserCredentialServiceClient
	profile    pbProfile.UserProfileServiceClient
	address    pbAddress.UserAddressServiceClient
	preference pbPreference.UserPreferenceServiceClient
}

type NotificationService struct {
	conn *grpc.ClientConn
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

	r := router.NewRouter(router.Router{
		Credential: userClient.credential,
		Profile:    userClient.profile,
		Address:    userClient.address,
		Preference: userClient.preference,
	})

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

	credential := pbCredential.NewUserCredentialServiceClient(conn)
	profile := pbProfile.NewUserProfileServiceClient(conn)
	address := pbAddress.NewUserAddressServiceClient(conn)
	preference := pbPreference.NewUserPreferenceServiceClient(conn)

	return UserServiceContainer{
		conn:       conn,
		credential: credential,
		profile:    profile,
		address:    address,
		preference: preference,
	}, err
}
