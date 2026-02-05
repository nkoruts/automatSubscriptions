package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/nkoruts/automatSubscriptions/internal/subscription"
)

type HTTPHandlers struct {
	subscriptionsList *subscription.List
}

func NewHTTPHandlers(list *subscription.List) *HTTPHandlers {
	return &HTTPHandlers{
		subscriptionsList: list,
	}
}

/*
pattern: /subscriptions
method:  GET
info:    -
*/
func (h *HTTPHandlers) HandleGetAllSubscriptions(w http.ResponseWriter, r *http.Request) {
	subscriptions := h.subscriptionsList.GetList()
	bytes, err := json.MarshalIndent(subscriptions, "", "    ")
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(bytes); err != nil {
		fmt.Println("failed to write http response:", err)
		return
	}
}

/*
pattern: /subscriptions
method: POST
info: JSON in HTTP request body
*/
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

	if err := h.subscriptionsList.AddSubscription(subsDTO.Owner, subsDTO.Days); err != nil {
		httpErrorIs(subscription.ErrSubscriptionAlreadyExists, err, w)
		return
	}

	successResp := SuccessDTO{Success: true}
	b, err := json.MarshalIndent(successResp.ToString(), "", " ")
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write http response:", err)
		return
	}
}

/*
pattern: /subscriptions/{key}
method: DELETE
info: pattern
*/
func (h *HTTPHandlers) HandleDeleteSubscription(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]

	if err := h.subscriptionsList.DeleteSubscription(key); err != nil {
		httpErrorIs(subscription.ErrSubscriptionNotFound, err, w)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

/*
pattern: /subscriptions/{key}
method: PATCH
info: pattern + JSON in body request
*/
func (h *HTTPHandlers) HandleUpdateSubscription(w http.ResponseWriter, r *http.Request) {
	var updateDTO UpdateDTO
	if err := json.NewDecoder(r.Body).Decode(&updateDTO); err != nil {
		httpError(w, err, http.StatusBadRequest)
		return
	}

	key := mux.Vars(r)["key"]

	if err := h.subscriptionsList.UpdateSubscription(key, updateDTO.DeviceId); err != nil {
		httpErrorIs(subscription.ErrSubscriptionNotFound, err, w)
		return
	}

	successResp := SuccessDTO{Success: true}
	b, err := json.MarshalIndent(successResp, "", " ")
	if err != nil {
		panic(err)
	}

	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write http response:", err)
		return
	}
}

/*
pattern: /subscriptions/check
method: POST
info: JSON in body request
*/
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

	if err := h.subscriptionsList.CheckSubscription(checkDTO.Key, checkDTO.DeviceID); err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}

		if errors.Is(err, subscription.ErrSubscriptionNotFound) {
			http.Error(w, errDTO.ToString(), http.StatusNotFound)
		} else if errors.Is(err, subscription.ErrUnregisteredDevice) {
			http.Error(w, errDTO.ToString(), http.StatusForbidden)
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		}
		return
	}

	successResp := SuccessDTO{Success: true}
	b, err := json.MarshalIndent(successResp, "", " ")
	if err != nil {
		panic(err)
	}

	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write http response:", err)
		return
	}
}
