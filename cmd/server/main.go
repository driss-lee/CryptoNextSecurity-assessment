package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/cryptonextsecurity/network-sniffer/docs" // Swagger docs
	"github.com/cryptonextsecurity/network-sniffer/internal/api"
	"github.com/cryptonextsecurity/network-sniffer/internal/config"
	"github.com/cryptonextsecurity/network-sniffer/internal/services"
	"github.com/cryptonextsecurity/network-sniffer/internal/storage"
	"github.com/cryptonextsecurity/network-sniffer/pkg/sniffing"
)

func main() {
	log.Println("Starting Network Sniffing Service...")

	// Load configuration from environment variables
	cfg := config.Load()

	log.Printf("Configuration: Storage Max Size=%d, Sniffing Interval=%v, Server Port=%s, Shutdown Timeout=%v",
		cfg.StorageMaxSize, cfg.SniffingInterval, cfg.ServerPort, cfg.ShutdownTimeout)

	// Create storage
	storage := storage.NewInMemoryStorage(cfg.StorageMaxSize)

	// Create sniffer
	sniffer := sniffing.NewPacketSniffer(storage, cfg.SniffingInterval)

	// Create service
	packetService := services.NewPacketService(storage, sniffer, nil)

	// Create handler and router
	handler := api.NewHandler(packetService, nil)
	router := api.NewRouter(handler, nil)
	ginRouter := router.Setup()

	// Start sniffing
	ctx := context.Background()
	log.Println("Starting packet sniffing...")
	packetService.StartSniffing(ctx)

	// Setup server
	server := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: ginRouter,
	}

	// Start server
	go func() {
		log.Printf("Server starting on port %s", cfg.ServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Stop sniffing
	packetService.StopSniffing(ctx)
	log.Println("Packet sniffing stopped")

	// Shutdown server
	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()
	server.Shutdown(shutdownCtx)
	log.Println("Server stopped")
}
