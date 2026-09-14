package users

import (
	"encoding/json"
	"net/http"

	"github.com/rdev-org/api-gateway/internal/utils"
	pb "github.com/rdev-org/proto-contracts/pb/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserHandler struct {
	userClient pb.UserServiceClient
}

func NewUserHandler(userClient pb.UserServiceClient) *UserHandler {
	return &UserHandler{
		userClient: userClient,
	}
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var reqDTO LoginRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&reqDTO); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	grpcResp, err := h.userClient.Login(r.Context(), &pb.UserLoginRequest{
		Username: reqDTO.Username,
		Password: reqDTO.Password,
	})

	if err := getError(err, w); err != nil {
		return
	}

	payload := LoginResponseDTO{
		Token:    grpcResp.GetToken(),
		Username: grpcResp.GetUsername(),
	}

	utils.ResponseOk(w, payload)
}

func getError(err error, w http.ResponseWriter) error {
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.InvalidArgument:
				http.Error(w, st.Message(), http.StatusBadRequest)
				return err
			case codes.Unauthenticated:
				http.Error(w, st.Message(), http.StatusUnauthorized)
				return err
			}
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return err
	}

	return err
}
