package server

import "github.com/nkoruts/automatSubscriptions/internal/subscription"

type HTTPHandlers struct {
	Subscriptions map[string]subscription.Subscription
}

func NewHTTPHandlers() *HTTPHandlers {
	return &HTTPHandlers{
		Subscriptions: make(map[string]subscription.Subscription),
	}
}
