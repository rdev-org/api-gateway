package user_address

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rdev-org/api-gateway/internal/shared"
	"github.com/rdev-org/api-gateway/internal/utils"
	pbCommon "github.com/rdev-org/proto-contracts/pb/common"
	pb "github.com/rdev-org/proto-contracts/pb/user_address"
)

type UserAddressHandler struct {
	c pb.UserAddressServiceClient
}

func NewUserAddressHandler(c pb.UserAddressServiceClient) *UserAddressHandler {
	return &UserAddressHandler{
		c: c,
	}
}

func (h *UserAddressHandler) GetAddress(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.ResponseError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	uuid := chi.URLParam(r, "uuid")

	grpcResp, err := h.c.GetAddress(r.Context(), &pb.GetAddressRequest{
		Uuid: uuid,
	})

	utils.GetStatusFromError(err, w)
	if err != nil {
		log.Printf("error encounter with status %v", err)
		return
	}

	utils.ResponseOk(w, grpcResp)
}

func (h *UserAddressHandler) ListAddress(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.ResponseError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	sf := shared.Filter{Query: r.URL.Query()}
	filter := sf.SanitizeFilter([]string{
		"address_line1",
		"address_line2",
		"city",
		"state_province",
		"postal_code",
		"country",
		"address_type",
	})

	grpcResp, err := h.c.ListAddress(r.Context(), &pb.ListAddressRequest{
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

func (h *UserAddressHandler) CreateAddress(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.ResponseError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var reqDTO UserAddressRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&reqDTO); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user_id := r.Header.Get("X-User-Id")
	grpcResp, err := h.c.CreateAddress(r.Context(), &pb.CreateAddressRequest{
		Address: &pb.Address{
			UserId:        user_id,
			AddressLine1:  reqDTO.AddressLine1,
			AddressLine2:  reqDTO.AddressLine2,
			City:          reqDTO.City,
			StateProvince: reqDTO.StateProvince,
			PostalCode:    reqDTO.PostalCode,
			Country:       reqDTO.Country,
			IsDefault:     reqDTO.IsDefault,
			AddressType:   reqDTO.AddressType,
		},
	})

	utils.GetStatusFromError(err, w)
	if err != nil {
		log.Printf("error encounter with status %v", err)
		return
	}

	utils.ResponseOk(w, map[string]string{"uuid": grpcResp.GetAddress().GetMeta().GetUuid()})
}

func (h *UserAddressHandler) UpdateAddress(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.ResponseError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var reqDTO UserAddressRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&reqDTO); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	uuid := chi.URLParam(r, "uuid")
	user_id := r.Header.Get("X-User-Id")
	grpcResp, err := h.c.UpdateAddress(r.Context(), &pb.UpdateAddressRequest{
		Uuid: uuid,
		Address: &pb.Address{
			UserId:        user_id,
			AddressLine1:  reqDTO.AddressLine1,
			AddressLine2:  reqDTO.AddressLine2,
			City:          reqDTO.City,
			StateProvince: reqDTO.StateProvince,
			PostalCode:    reqDTO.PostalCode,
			Country:       reqDTO.Country,
			IsDefault:     reqDTO.IsDefault,
			AddressType:   reqDTO.AddressType,
		},
	})

	utils.GetStatusFromError(err, w)
	if err != nil {
		log.Printf("error encounter with status %v", err)
		return
	}

	utils.ResponseOk(w, map[string]string{"uuid": grpcResp.GetAddress().GetMeta().GetUuid()})
}

func (h *UserAddressHandler) DeleteAddress(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.ResponseError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	uuid := chi.URLParam(r, "uuid")
	grpcResp, err := h.c.DeleteAddress(r.Context(), &pb.DeleteAddressRequest{
		Uuid: uuid,
	})

	utils.GetStatusFromError(err, w)
	if err != nil {
		log.Printf("error encounter with status %v", err)
		return
	}

	utils.ResponseOk(w, grpcResp)
}

func (h *UserAddressHandler) DefaultAddress(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.ResponseError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	user_id := r.Header.Get("X-User-Id")
	uuid := chi.URLParam(r, "uuid")
	grpcResp, err := h.c.DefaultAddress(r.Context(), &pb.DefaultAddressRequest{
		UserId: user_id,
		Uuid:   uuid,
	})

	utils.GetStatusFromError(err, w)
	if err != nil {
		log.Printf("error encounter with status %v", err)
		return
	}

	utils.ResponseOk(w, grpcResp)
}
