package user_addresses

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rdev-org/api-gateway/internal/shared"
	"github.com/rdev-org/api-gateway/internal/utils"
	pb "github.com/rdev-org/proto-contracts/pb/user_address"
)

type UserAddressHandler struct {
	client pb.UserAddressServiceClient
}

func NewUserAddressHandler(client pb.UserAddressServiceClient) *UserAddressHandler {
	return &UserAddressHandler{
		client: client,
	}
}

func (h *UserAddressHandler) GetAddress(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.ResponseError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	uuid := chi.URLParam(r, "uuid")

	grpcResp, err := h.client.GetAddress(r.Context(), &pb.GetAddressRequest{
		Uuid: uuid,
	})

	utils.GetStatusFromError(err, w)
	if err != nil {
		log.Printf("error encounter with status %v", err)
		return
	}

	utils.ResponseOk(w, grpcResp)
}

func (h *UserAddressHandler) GetAllAddress(w http.ResponseWriter, r *http.Request) {
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

	grpcResp, err := h.client.GetAllAddress(r.Context(), &pb.GetAllAddressRequest{
		Filter: &pb.AddressFilter{
			Sort:   filter.Sort,
			Search: filter.Search,
		},
		Pagination: &pb.PaginationRequest{
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

	grpcResp, err := h.client.CreateAddress(r.Context(), &pb.CreateAddressRequest{
		Address: &pb.Address{
			UserId:        reqDTO.UserID,
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

	utils.ResponseOk(w, map[string]string{"uuid": grpcResp.GetAddress().GetUuid()})
}

func (h *UserAddressHandler) UpdateAddress(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.ResponseError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	uuid := chi.URLParam(r, "uuid")
	var reqDTO UserAddressRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&reqDTO); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	grpcResp, err := h.client.UpdateAddress(r.Context(), &pb.UpdateAddressRequest{
		Address: &pb.Address{
			UserId:        reqDTO.UserID,
			AddressLine1:  reqDTO.AddressLine1,
			AddressLine2:  reqDTO.AddressLine2,
			City:          reqDTO.City,
			StateProvince: reqDTO.StateProvince,
			PostalCode:    reqDTO.PostalCode,
			Country:       reqDTO.Country,
			IsDefault:     reqDTO.IsDefault,
			AddressType:   reqDTO.AddressType,
			Uuid:          uuid,
		},
	})

	utils.GetStatusFromError(err, w)
	if err != nil {
		log.Printf("error encounter with status %v", err)
		return
	}

	utils.ResponseOk(w, map[string]string{"uuid": grpcResp.GetAddress().GetUuid()})
}

func (h *UserAddressHandler) DeleteAddress(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.ResponseError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	uuid := chi.URLParam(r, "uuid")
	grpcResp, err := h.client.DeleteAddress(r.Context(), &pb.DeleteAddressRequest{
		Uuid: uuid,
	})

	utils.GetStatusFromError(err, w)
	if err != nil {
		log.Printf("error encounter with status %v", err)
		return
	}

	utils.ResponseOk(w, grpcResp)
}

func (h *UserAddressHandler) SetDefaultAddress(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.ResponseError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	user_id := r.Header.Get("x-user-id")
	uuid := chi.URLParam(r, "uuid")
	grpcResp, err := h.client.SetDefaultAddress(r.Context(), &pb.SetDefaultAddressRequest{
		UserId: user_id,
		Uuid:   uuid,
	})

	utils.GetStatusFromError(err, w)
	if err != nil {
		log.Printf("error encounter with status %v", err)
		return
	}

	utils.ResponseOk(w, map[string]string{"uuid": grpcResp.GetAddress().GetUuid()})
}
