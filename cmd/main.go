package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/nkoruts/automatSubscriptions/internal/server"
	"github.com/nkoruts/automatSubscriptions/internal/storage"
)

func main() {
	fmt.Println("Start HTTP server.")

	dsn := os.Getenv("DATABASE_URL")
	ctx := context.Background()

	storage, err := storage.NewPostgresStorage(ctx, dsn)
	if err != nil {
		log.Fatal(err)
	}

	err = storage.CreateTable()
	if err != nil {
		log.Fatal(err)
	}

	httpHandlers := server.NewHTTPHandlers(storage)
	httpServer := server.NewHTTPServer(httpHandlers)
	if err := httpServer.StartServer(); err != nil {
		fmt.Println("HTTP server error:", err)
	}
}
