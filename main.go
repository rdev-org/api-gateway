package main

import (
	"log"
	"net/http"
	"os"

	"github.com/rdev-org/api-gateway/internal/container"
)

func main() {
	appContainer, routes := container.NewContainer()
	defer appContainer.UserConn.Close()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("API Gateway is running on port %s...", port)
	if err := http.ListenAndServe("localhost:"+port, routes); err != nil {
		log.Fatalf("failed to start HTTP server: %v", err)
	}
}
