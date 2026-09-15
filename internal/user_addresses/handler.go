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
	sf := shared.Filter{Query: r.URL.Query()}
	filter := sf.SanitizeFilter([]string{"city"})

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
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var reqDTO UserAddressRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&reqDTO); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
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
