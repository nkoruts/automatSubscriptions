package main

import (
	"fmt"

	"github.com/nkoruts/automatSubscriptions/internal/server"
	"github.com/nkoruts/automatSubscriptions/internal/subscription"
)

func main() {
	subsList := subscription.NewList()
	httpHandlers := server.NewHTTPHandlers(subsList)
	httpServer := server.NewHTTPServer(httpHandlers)

	if err := httpServer.StartServer(); err != nil {
		fmt.Println("HTTP server error:", err)
	}
}
