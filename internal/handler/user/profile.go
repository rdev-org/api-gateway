package user

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rdev-org/api-gateway/internal/dto"
	"github.com/rdev-org/api-gateway/internal/shared"
	"github.com/rdev-org/api-gateway/internal/utils"
	pbCommon "github.com/rdev-org/proto-contracts/pb/common"
	pb "github.com/rdev-org/proto-contracts/pb/user"
)

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.ResponseError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	uuid := r.Header.Get("X-User-Id")
	grpcResp, err := h.c.GetProfile(r.Context(), &pb.GetProfileRequest{
		Uuid: uuid,
	})

	utils.GetStatusFromError(err, w)
	if err != nil {
		log.Printf("error encounter with status %v", err)
		return
	}

	utils.ResponseOk(w, grpcResp)
}

func (h *UserHandler) ListProfile(w http.ResponseWriter, r *http.Request) {
	sf := shared.Filter{Query: r.URL.Query()}
	filter := sf.SanitizeFilter([]string{
		"first_name",
		"middle_name",
		"last_name",
		"mobile_number",
		"date_of_birth",
		"gender",
	})

	grpcResp, err := h.c.ListProfile(r.Context(), &pb.ListProfileRequest{
		Filter: &pbCommon.Filter{
			Sort:   filter.Sort,
			Search: filter.Search,
		},
		Pagination: &pbCommon.PaginationRequest{
			PageSize: int64(filter.PageSize),
			Page:     int64(filter.Page),
		},
	})

	utils.GetStatusFromError(err, w)
	if err != nil {
		log.Printf("error encounter with status %v", err)
		return
	}

	utils.ResponseOk(w, grpcResp)
}

func (h *UserHandler) CreateProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.ResponseError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var reqDTO dto.ProfileRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&reqDTO); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user_id := r.Header.Get("X-User-Id")
	email := r.URL.Query().Get("email")
	grpcResp, err := h.c.CreateProfile(r.Context(), &pb.CreateProfileRequest{
		Payload: &pb.CreateProfilePayload{
			UserId:       user_id,
			Email:        email,
			FirstName:    reqDTO.FirstName,
			MiddleName:   reqDTO.MiddleName,
			LastName:     reqDTO.LastName,
			MobileNumber: reqDTO.MobileNumber,
			DateOfBirth:  reqDTO.DateOfBirth,
			Gender:       reqDTO.Gender,
			Avatar:       reqDTO.Avatar,
			Address: &pb.Address{
				AddressLine1:  reqDTO.Address.AddressLine1,
				AddressLine2:  reqDTO.Address.AddressLine2,
				City:          reqDTO.Address.City,
				StateProvince: reqDTO.Address.StateProvince,
				PostalCode:    reqDTO.Address.PostalCode,
				Country:       reqDTO.Address.Country,
				IsDefault:     true,
				AddressType:   reqDTO.Address.AddressType,
			},
		},
	})

	utils.GetStatusFromError(err, w)
	if err != nil {
		log.Printf("error encounter with status %v", err)
		return
	}

	utils.ResponseCreated(w, grpcResp)
}

func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.ResponseError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var reqDTO dto.ProfileRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&reqDTO); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	uuid := chi.URLParam(r, "uuid")
	user_id := r.Header.Get("X-User-Id")
	grpcResp, err := h.c.UpdateProfile(r.Context(), &pb.UpdateProfileRequest{
		Uuid: uuid,
		Profile: &pb.Profile{
			UserId:       user_id,
			FirstName:    reqDTO.FirstName,
			MiddleName:   reqDTO.MiddleName,
			LastName:     reqDTO.LastName,
			MobileNumber: reqDTO.MobileNumber,
			DateOfBirth:  reqDTO.DateOfBirth,
			Gender:       reqDTO.Gender,
			Avatar:       reqDTO.Avatar,
		},
	})

	utils.GetStatusFromError(err, w)
	if err != nil {
		log.Printf("error encounter with status %v", err)
		return
	}

	utils.ResponseOk(w, grpcResp)
}

func (h *UserHandler) DeleteProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.ResponseError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	uuid := chi.URLParam(r, "uuid")
	grpcResp, err := h.c.DeleteProfile(r.Context(), &pb.DeleteProfileRequest{
		Uuid: uuid,
	})

	utils.GetStatusFromError(err, w)
	if err != nil {
		log.Printf("error encounter with status %v", err)
		return
	}

	utils.ResponseOk(w, grpcResp)
}
