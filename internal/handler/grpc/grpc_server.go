package grpc

import (
	"Gober/configs"
	"Gober/internal/generated/grpc/gober"
	"Gober/internal/repo/mysql"
	"Gober/internal/service"
	"Gober/pkg/cache"
	"Gober/pkg/jwt"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"
)

type GRPCServer interface {
	StartGRPCServer(ctx context.Context) error
}

type grpcServer struct {
	config  *configs.Config
	handler gober.GoberServiceServer
	server  *grpc.Server
	cache   cache.RedisCacheService
	token   jwt.TokenService
	mu      sync.RWMutex
}

func (g *grpcServer) StartGRPCServer(ctx context.Context) error {
	lis, err := net.Listen("tcp", ":"+g.config.Server.PortGrpc)
	if err != nil {
		return fmt.Errorf("failed to listen on port %s: %w", g.config.Server.PortGrpc, err)
	}

	g.mu.Lock()
	g.server = grpc.NewServer()
	g.mu.Unlock()

	// Initialize dependencies
	if err := g.initializeDependencies(); err != nil {
		return fmt.Errorf("failed to initialize dependencies: %w", err)
	}

	// Channel to capture server errors
	serverErr := make(chan error, 1)

	// Start server in goroutine
	go func() {
		log.Printf("gRPC server listening at :%s", g.config.Server.PortGrpc)
		if err := g.server.Serve(lis); err != nil {
			serverErr <- fmt.Errorf("gRPC server serve error: %w", err)
		}
	}()

	// Wait for context cancellation or server error
	select {
	case <-ctx.Done():
		log.Println("Context cancelled, shutting down gRPC server...")
	case err := <-serverErr:
		log.Printf("Server error: %v", err)
		return err
	}

	// Graceful shutdown
	g.mu.RLock()
	if g.server != nil {
		log.Println("Gracefully stopping gRPC server...")
		g.server.GracefulStop()
	}
	g.mu.RUnlock()

	return nil
}

func (g *grpcServer) initializeDependencies() error {
	// Initialize database connection
	db, err := mysql.InitDB(g.config)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	// Initialize repository
	accountRepo := mysql.NewAccountDatabase(db)
	eventRepo := mysql.NewEventDatabase(db)
	ticketRepo := mysql.NewTicketDatabase(db)

	// Initialize services
	hashService := service.NewHash()
	accountService := service.NewAccountService(accountRepo, g.token, hashService, g.cache)
	eventService := service.NewEventService(eventRepo)
	ticketService := service.NewTicketService(ticketRepo)

	// Initialize handler
	handler, err := NewGoberHandler(accountService, eventService, ticketService)

	if err != nil {
		return err
	}

	g.mu.Lock()
	gober.RegisterGoberServiceServer(g.server, handler)
	g.mu.Unlock()

	return nil
}

func NewGRPCServer(config *configs.Config, cache cache.RedisCacheService, token jwt.TokenService) GRPCServer {
	return &grpcServer{
		config: config,
		cache:  cache,
		token:  token,
	}
}
