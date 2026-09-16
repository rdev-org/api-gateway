package user_credential

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/rdev-org/api-gateway/internal/utils"
	pb "github.com/rdev-org/proto-contracts/pb/user_credential"
)

type UserCredentialHandler struct {
	c pb.UserCredentialServiceClient
}

func NewUserCredentialHandler(c pb.UserCredentialServiceClient) *UserCredentialHandler {
	return &UserCredentialHandler{
		c: c,
	}
}

func (h *UserCredentialHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.ResponseError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var reqDTO LoginRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&reqDTO); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	grpcResp, err := h.c.Login(r.Context(), &pb.LoginRequest{
		Username: reqDTO.Username,
		Password: reqDTO.Password,
	})

	utils.GetStatusFromError(err, w)
	if err != nil {
		log.Printf("error encounter with status %v", err)
		return
	}

	w.Header().Set("X-User-Id", grpcResp.GetUuid())

	utils.ResponseOk(w, grpcResp)
}

func (h *UserCredentialHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.ResponseError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var reqDTO RegisterRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&reqDTO); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	grpcResp, err := h.c.Register(r.Context(), &pb.RegisterRequest{
		Email:    reqDTO.Email,
		Password: reqDTO.Password,
	})

	utils.GetStatusFromError(err, w)
	if err != nil {
		return
	}

	utils.ResponseCreated(w, map[string]string{"message": grpcResp.GetMessage()})
}

func (h *UserCredentialHandler) DeleteCredential(w http.ResponseWriter, r *http.Request) {
	//
}

func (h *UserCredentialHandler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.ResponseError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var reqDTO PasswordRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&reqDTO); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	uuid := r.Header.Get("X-User-Id")
	grpcResp, err := h.c.UpdatePassword(r.Context(), &pb.UpdatePasswordRequest{
		Uuid:     uuid,
		Password: reqDTO.Password,
	})

	utils.GetStatusFromError(err, w)
	if err != nil {
		return
	}

	utils.ResponseCreated(w, map[string]string{"message": grpcResp.GetMessage()})
}

func (h *UserCredentialHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.ResponseError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var reqDTO PasswordRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&reqDTO); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	email := r.URL.Query().Get("email")
	uuid := r.URL.Query().Get("uuid")
	token := r.URL.Query().Get("token")
	grpcResp, err := h.c.ResetPassword(r.Context(), &pb.ResetPasswordRequest{
		Email:    email,
		Password: reqDTO.Password,
		Uuid:     uuid,
		Token:    token,
	})

	utils.GetStatusFromError(err, w)
	if err != nil {
		return
	}

	utils.ResponseCreated(w, map[string]string{"message": grpcResp.GetMessage()})
}
