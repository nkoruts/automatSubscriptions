package server

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

type HTTPServer struct {
	httpHandlers *HTTPHandlers
}

func NewHTTPServer(httpHandler *HTTPHandlers) *HTTPServer {
	return &HTTPServer{
		httpHandlers: httpHandler,
	}
}

func (s *HTTPServer) StartServer() error {
	router := mux.NewRouter()

	router.Path("/subscriptions").Methods(http.MethodGet).HandlerFunc(s.httpHandlers.HandleGetAllSubscriptions)
	router.Path("/subscriptions").Methods(http.MethodPost).HandlerFunc(s.httpHandlers.HandleCreateSubscription)
	router.Path("/subscriptions/{key}").Methods(http.MethodDelete).HandlerFunc(s.httpHandlers.HandleDeleteSubscription)
	router.Path("/subscriptions/{key}").Methods(http.MethodPatch).HandlerFunc(s.httpHandlers.HandleUpdateSubscription)
	router.Path("/subscriptions/check").Methods(http.MethodPost).HandlerFunc(s.httpHandlers.HandleCheckSubscription)

	if err := http.ListenAndServe(":9091", router); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
	return nil
}
