package users

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/rdev-org/api-gateway/internal/utils"
	pb "github.com/rdev-org/proto-contracts/pb/user"
)

type UserHandler struct {
	client pb.UserServiceClient
}

func NewUserHandler(client pb.UserServiceClient) *UserHandler {
	return &UserHandler{
		client: client,
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

	grpcResp, err := h.client.Login(r.Context(), &pb.UserLoginRequest{
		Username: reqDTO.Username,
		Password: reqDTO.Password,
	})

	utils.GetStatusFromError(err, w)
	if err != nil {
		log.Printf("error encounter with status %v", err)
		return
	}

	utils.ResponseOk(w, grpcResp)
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var reqDTO RegisterRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&reqDTO); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	grpcResp, err := h.client.Register(r.Context(), &pb.UserRegisterRequest{
		Username:     reqDTO.Username,
		Email:        reqDTO.Email,
		PasswordHash: reqDTO.PasswordHash,
		FirstName:    reqDTO.FirstName,
		MiddleName:   reqDTO.MiddleName,
		LastName:     reqDTO.LastName,
	})

	utils.GetStatusFromError(err, w)
	if err != nil {
		return
	}

	utils.ResponseCreated(w, map[string]string{"message": "account created for " + grpcResp.GetUsername()})
}
