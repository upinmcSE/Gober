package main

import (
	"Gober/internal/handler/grpc"
	"Gober/internal/handler/http"
	"log"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	// Start gRPC server
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered in gRPC server: %v", r)
			}
		}()
		grpcServer := grpc.NewGRPCServer(nil)
		grpcServer.StartGRPCServer(&wg)
	}()

	// Start HTTP server
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered in HTTP server: %v", r)
			}
		}()
		httpServer := http.NewHTTPServer()
		httpServer.StartHTTPServer(&wg)
	}()

	wg.Wait()
}
