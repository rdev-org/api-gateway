package user_preference

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rdev-org/api-gateway/internal/shared"
	"github.com/rdev-org/api-gateway/internal/utils"
	pbCommon "github.com/rdev-org/proto-contracts/pb/common"
	pb "github.com/rdev-org/proto-contracts/pb/user_preference"
	"google.golang.org/protobuf/types/known/structpb"
)

type UserPreferenceHandler struct {
	c pb.UserPreferenceServiceClient
}

func NewUserPreferenceHandler(c pb.UserPreferenceServiceClient) *UserPreferenceHandler {
	return &UserPreferenceHandler{
		c: c,
	}
}

func (h *UserPreferenceHandler) GetPreference(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.ResponseError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	uuid := r.Header.Get("X-User-Id")
	grpcResp, err := h.c.GetPreference(r.Context(), &pb.GetPreferenceRequest{
		Uuid: uuid,
	})

	utils.GetStatusFromError(err, w)
	if err != nil {
		log.Printf("error encounter with status %v", err)
		return
	}

	utils.ResponseOk(w, grpcResp)
}

func (h *UserPreferenceHandler) ListPreference(w http.ResponseWriter, r *http.Request) {
	sf := shared.Filter{Query: r.URL.Query()}
	filter := sf.SanitizeFilter([]string{
		"language",
		"currency",
	})

	grpcResp, err := h.c.ListPreference(r.Context(), &pb.ListPreferenceRequest{
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

func (h *UserPreferenceHandler) CreatePreference(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.ResponseError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var reqDTO PreferenceRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&reqDTO); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	settingsAny := make(map[string]any, len(reqDTO.Settings))
	for k, v := range reqDTO.Settings {
		settingsAny[k] = v
	}

	settingsStruct, err := structpb.NewStruct(settingsAny)
	if err != nil {
		log.Printf("failed to convert settings to structpb: %v", err)
	}

	user_id := r.Header.Get("X-User-Id")
	grpcResp, err := h.c.CreatePreference(r.Context(), &pb.CreatePreferenceRequest{
		Preference: &pb.Preference{
			UserId:     user_id,
			Language:   reqDTO.Language,
			Currency:   reqDTO.Currency,
			EmailOptIn: reqDTO.EmailOptIn,
			SmsOptIn:   reqDTO.SmsOptIn,
			Settings:   settingsStruct,
		},
	})

	utils.GetStatusFromError(err, w)
	if err != nil {
		log.Printf("error encounter with status %v", err)
		return
	}

	utils.ResponseCreated(w, grpcResp)
}

func (h *UserPreferenceHandler) UpdatePreference(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.ResponseError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var reqDTO PreferenceRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&reqDTO); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	uuid := chi.URLParam(r, "uuid")
	user_id := r.Header.Get("X-User-Id")

	settingsAny := make(map[string]any, len(reqDTO.Settings))
	for k, v := range reqDTO.Settings {
		settingsAny[k] = v
	}

	settingsStruct, err := structpb.NewStruct(settingsAny)
	if err != nil {
		log.Printf("failed to convert settings to structpb: %v", err)
	}

	grpcResp, err := h.c.UpdatePreference(r.Context(), &pb.UpdatePreferenceRequest{
		Uuid: uuid,
		Preference: &pb.Preference{
			UserId:     user_id,
			Language:   reqDTO.Language,
			Currency:   reqDTO.Currency,
			EmailOptIn: reqDTO.EmailOptIn,
			SmsOptIn:   reqDTO.SmsOptIn,
			Settings:   settingsStruct,
		},
	})

	utils.GetStatusFromError(err, w)
	if err != nil {
		log.Printf("error encounter with status %v", err)
		return
	}

	utils.ResponseOk(w, grpcResp)
}

func (h *UserPreferenceHandler) DeletePreference(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.ResponseError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	uuid := chi.URLParam(r, "uuid")
	grpcResp, err := h.c.DeletePreference(r.Context(), &pb.DeletePreferenceRequest{
		Uuid: uuid,
	})

	utils.GetStatusFromError(err, w)
	if err != nil {
		log.Printf("error encounter with status %v", err)
		return
	}

	utils.ResponseOk(w, grpcResp)
}
