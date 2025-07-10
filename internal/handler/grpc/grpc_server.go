package grpc

import (
	"Gober/configs"
	"Gober/internal/generated/grpc/gober"
	"Gober/internal/repo/mysql"
	"Gober/internal/service"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"
)

type GRPCServer interface {
	StartGRPCServer(wg *sync.WaitGroup)
}

type grpcServer struct {
	handler gober.GoberServiceServer
}

func (g grpcServer) StartGRPCServer(wg *sync.WaitGroup) {
	defer wg.Done()

	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	// Load configuration
	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Initialize database connection
	db, err := mysql.InitDB(&config)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Initialize repository
	accountRepo := mysql.NewAccountDatabase(db)

	// Initialize service
	hashService := service.NewHash()
	tokenService := service.NewTokenService()
	accountService := service.NewAccountService(accountRepo, tokenService, hashService)

	// Initialize handler
	handler, err := NewAccountHandler(accountService)
	if err != nil {
		log.Fatalf("failed to create account handler: %v", err)
	}

	gober.RegisterGoberServiceServer(s, handler)

	log.Println("gRPC server listening at :8080")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func NewGRPCServer(handler gober.GoberServiceServer) GRPCServer {
	return &grpcServer{
		handler: handler,
	}
}
