package main

import (
	"log"
	"net/http"
	"os"

	"github.com/whu-ham/ham-gateway/gen/gatewayv1connect"
	"github.com/whu-ham/ham-gateway/internal/ham"
	"github.com/whu-ham/ham-gateway/internal/services"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	baseURL := os.Getenv("HAM_API_BASE_URL")
	appID := os.Getenv("HAM_OPEN_APP_ID")
	appSecret := os.Getenv("HAM_OPEN_APP_SECRET")
	certPath := os.Getenv("HAM_GRPC_MTLS_CLIENT_CRT")
	keyPath := os.Getenv("HAM_GRPC_MTLS_CLIENT_KEY")

	if baseURL == "" {
		baseURL = "https://ham.nowcent.cn"
	}

	hamClient, err := ham.NewClient(baseURL, appID, appSecret, certPath, keyPath)
	if err != nil {
		log.Fatalf("Failed to create HAM client: %v", err)
	}

	authToken := os.Getenv("GATEWAY_AUTH_TOKEN")
	gatewayServer := services.NewGatewayServer(hamClient, authToken)
	mux := http.NewServeMux()
	path, handler := gatewayv1connect.NewGatewayServiceHandler(gatewayServer)
	mux.Handle(path, handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting HAM Gateway on :%s", port)
	err = http.ListenAndServe(
		":"+port,
		h2c.NewHandler(mux, &http2.Server{}),
	)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
