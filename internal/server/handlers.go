package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/nkoruts/automatSubscriptions/internal/subscription"
)

type HTTPHandlers struct {
	subscriptionsList subscription.List
}

func NewHTTPHandlers() *HTTPHandlers {
	return &HTTPHandlers{
		subscriptionsList: *subscription.NewList(),
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
	var subscriptionDTO SubscriptionDTO
	if err := json.NewDecoder(r.Body).Decode(&subscriptionDTO); err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}

		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	if err := subscriptionDTO.ValidateForCreate(); err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}

		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	if err := h.subscriptionsList.AddSubscription(subscriptionDTO.Owner, subscriptionDTO.Days); err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}

		if errors.Is(err, subscription.ErrSubscriptionAlreadyExists) {
			http.Error(w, errDTO.ToString(), http.StatusConflict)
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		}

		return
	}

	successDTO := SuccessDTO{Success: true}
	b, err := json.MarshalIndent(successDTO.ToString(), "", " ")
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write http response:", err)
		return
	}
}
