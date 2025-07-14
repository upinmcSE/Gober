package grpc

import (
	"Gober/configs"
	"Gober/internal/generated/grpc/gober"
	"Gober/internal/repo/mysql"
	"Gober/internal/service"
	"context"
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
	mu      sync.RWMutex
}

func (g *grpcServer) StartGRPCServer(ctx context.Context) error {

	lis, err := net.Listen("tcp", ":"+g.config.Server.PortGrpc)
	if err != nil {
		return err
	}

	g.mu.Lock()
	g.server = grpc.NewServer()
	g.mu.Unlock()

	// Initialize dependencies
	if err := g.initializeDependencies(); err != nil {
		return err
	}

	// Start server in goroutine
	go func() {
		log.Println("gRPC server listening at :8080")
		if err := g.server.Serve(lis); err != nil {
			log.Printf("gRPC server error: %v", err)
		}
	}()

	// Wait for context cancellation
	<-ctx.Done()

	// Graceful shutdown
	g.mu.RLock()
	if g.server != nil {
		g.server.GracefulStop()
	}
	g.mu.RUnlock()

	return nil
}

func (g *grpcServer) initializeDependencies() error {

	// Initialize database connection
	db, err := mysql.InitDB(g.config)
	if err != nil {
		return err
	}

	// Initialize repository
	accountRepo := mysql.NewAccountDatabase(db)

	// Initialize services
	hashService := service.NewHash()
	tokenService := service.NewTokenService()
	accountService := service.NewAccountService(accountRepo, tokenService, hashService)

	// Initialize handler
	accountHandler, err := NewAccountHandler(accountService)
	if err != nil {
		return err
	}

	g.mu.RLock()
	gober.RegisterGoberServiceServer(g.server, accountHandler)
	g.mu.RUnlock()

	return nil
}

func NewGRPCServer(config *configs.Config) GRPCServer {
	return &grpcServer{
		config: config,
	}
}
