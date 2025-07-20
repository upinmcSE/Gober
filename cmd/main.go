package main

import (
	configs2 "Gober/configs"
	"Gober/internal/handler/grpc"
	"Gober/internal/handler/http"
	redis "Gober/pkg/cache"
	"Gober/pkg/jwt"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Initialize configurations
	config, err := configs2.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	cache := redis.NewRedisClient(&config)
	cacheService := redis.NewRedisCacheService(cache)

	token := jwt.NewTokenService(cacheService, &config)

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Channel to receive OS signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Channel to receive errors from servers
	errChan := make(chan error, 2)

	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Start gRPC server
	grpcServer := grpc.NewGRPCServer(&config, cacheService, token)
	go func() {
		if err := grpcServer.StartGRPCServer(ctx); err != nil {
			errChan <- err
		}
	}()

	// Wait for gRPC server to be ready
	time.Sleep(2 * time.Second)

	// Start HTTP server
	httpServer := http.NewHTTPServer(&config, cacheService, token)
	go func() {
		if err := httpServer.StartHTTPServer(ctx); err != nil {
			errChan <- err
		}
	}()

	// Wait for termination signal or error
	select {
	case <-sigChan:
		log.Println("Received termination signal, shutting down...")
	case err := <-errChan:
		log.Printf("Server error: %v", err)
	}

	// Graceful shutdown
	cancel()
	time.Sleep(5 * time.Second) // Give servers time to shutdown
	log.Println("Servers stopped")
}
