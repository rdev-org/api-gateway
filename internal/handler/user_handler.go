package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

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

type LoginRequestDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponseDTO struct {
	Token string `json:"token"`
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

	fmt.Println("ERROR: ", err)

	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.Unauthenticated:
				http.Error(w, st.Message(), http.StatusUnauthorized)
				return
			case codes.InvalidArgument:
				http.Error(w, st.Message(), http.StatusBadRequest)
				return
			}
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(LoginResponseDTO{
		Token: grpcResp.GetToken(),
	})
}
