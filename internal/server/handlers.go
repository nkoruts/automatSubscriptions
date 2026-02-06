package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/nkoruts/automatSubscriptions/internal/subscription"
)

type SubscriptionStorage interface {
	AddSubscription(owner string, days int) error
	DeleteSubscription(key string) error
	UpdateSubscription(key, deviceId string) error
	CheckSubscription(key, deviceId string) (bool, error)
	GetList() map[string]subscription.Subscription
}

type HTTPHandlers struct {
	storage SubscriptionStorage
}

func NewHTTPHandlers(storage SubscriptionStorage) *HTTPHandlers {
	return &HTTPHandlers{
		storage: storage,
	}
}

// GET /subscriptions
func (h *HTTPHandlers) HandleGetAllSubscriptions(w http.ResponseWriter, r *http.Request) {
	subscriptions := h.storage.GetList()
	writeJSON(w, http.StatusOK, subscriptions)
}

// POST /subscriptions
func (h *HTTPHandlers) HandleCreateSubscription(w http.ResponseWriter, r *http.Request) {
	var subsDTO SubscriptionDTO
	if err := json.NewDecoder(r.Body).Decode(&subsDTO); err != nil {
		httpError(w, err, http.StatusBadRequest)
		return
	}

	if err := subsDTO.ValidateRequest(); err != nil {
		httpError(w, err, http.StatusBadRequest)
		return
	}

	if err := h.storage.AddSubscription(subsDTO.Owner, subsDTO.Days); err != nil {
		httpErrorIs(err, subscription.ErrSubscriptionAlreadyExists, w)
		return
	}

	writeJSON(w, http.StatusCreated, SuccessDTO{Success: true})
}

// DELETE /subscriptions/{key}
func (h *HTTPHandlers) HandleDeleteSubscription(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]

	if err := h.storage.DeleteSubscription(key); err != nil {
		httpErrorIs(err, subscription.ErrSubscriptionNotFound, w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// PATCH /subscriptions/{key}
func (h *HTTPHandlers) HandleUpdateSubscription(w http.ResponseWriter, r *http.Request) {
	var updateDTO UpdateDTO
	if err := json.NewDecoder(r.Body).Decode(&updateDTO); err != nil {
		httpError(w, err, http.StatusBadRequest)
		return
	}

	if err := updateDTO.ValidateRequest(); err != nil {
		httpError(w, err, http.StatusBadRequest)
		return
	}

	key := mux.Vars(r)["key"]

	if err := h.storage.UpdateSubscription(key, updateDTO.DeviceId); err != nil {
		httpErrorIs(err, subscription.ErrSubscriptionNotFound, w)
		return
	}

	writeJSON(w, http.StatusOK, SuccessDTO{Success: true})
}

// POST /subscriptions/check
func (h *HTTPHandlers) HandleCheckSubscription(w http.ResponseWriter, r *http.Request) {
	var checkDTO CheckDTO
	if err := json.NewDecoder(r.Body).Decode(&checkDTO); err != nil {
		httpError(w, err, http.StatusBadRequest)
		return
	}

	if err := checkDTO.ValidateRequest(); err != nil {
		httpError(w, err, http.StatusBadRequest)
		return
	}

	success, err := h.storage.CheckSubscription(checkDTO.Key, checkDTO.DeviceID)
	if err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}
		switch {
		case errors.Is(err, subscription.ErrSubscriptionNotFound):
			http.Error(w, errDTO.ToString(), http.StatusNotFound)
		case errors.Is(err, subscription.ErrUnregisteredUserDevice):
			http.Error(w, errDTO.ToString(), http.StatusForbidden)
		default:
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		}
		return
	}

	if !success {
		if err := h.storage.UpdateSubscription(checkDTO.Key, checkDTO.DeviceID); err != nil {
			httpErrorIs(err, subscription.ErrSubscriptionNotFound, w)
			return
		}
		success = true
	}

	writeJSON(w, http.StatusOK, SuccessDTO{Success: success})
}
