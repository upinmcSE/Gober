package http

import (
	"Gober/configs"
	"Gober/internal/generated/grpc/gober"
	"Gober/internal/middleware"
	"Gober/pkg/cache"
	"Gober/pkg/jwt"
	"Gober/pkg/logger"
	"context"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"sync"
	"time"
)

type HTTPServer interface {
	StartHTTPServer(ctx context.Context) error
}

type httpServer struct {
	config *configs.Config
	server *http.Server
	conn   *grpc.ClientConn
	cache  cache.RedisCacheService
	token  jwt.TokenService
	mu     sync.RWMutex
}

func (h *httpServer) StartHTTPServer(ctx context.Context) error {
	// Initialize gRPC client connection
	if err := h.initGRPCClient(h.config.Server.PortGrpc); err != nil {
		return err
	}

	// Setup HTTP server
	router := h.setupRouter()

	h.mu.Lock()
	h.server = &http.Server{
		Addr:           ":" + h.config.Server.PortHttp,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    15 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	h.mu.Unlock()

	// Start server in goroutine
	go func() {
		log.Println("HTTP server listening at :8082")
		if err := h.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	// Wait for context cancellation
	<-ctx.Done()

	// Graceful shutdown
	return h.shutdown()
}

func (h *httpServer) initGRPCClient(portGrpc string) error {
	var err error
	h.conn, err = grpc.Dial("localhost:"+portGrpc,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithTimeout(5*time.Second),
	)
	return err
}

func (h *httpServer) setupRouter() *gin.Engine {
	r := gin.Default()

	// Setup API versioning
	v1 := r.Group("/v1/2025")

	// Initialize loggers
	httpLogger := logger.NewLoggerWithPath("logs/http.log", "info")
	rateLimiterLogger := logger.NewLoggerWithPath("logs/rate_limiter.log", "warning")

	// Initialize middleware
	middleware.InitAuthMiddleware(h.token, h.cache)

	// Apply middleware
	v1.Use(
		middleware.ApikeyMiddleware(),
		middleware.CORSMiddleware(),
		middleware.LoggerMiddleware(httpLogger),
		middleware.RateLimiterMiddleware(rateLimiterLogger),
		//gzip.Gzip(gzip.BestCompression),
	)

	// Setup routes
	h.setupRoutes(v1)

	return r
}

func (h *httpServer) setupRoutes(rg *gin.RouterGroup) {
	client := gober.NewGoberServiceClient(h.conn)
	accountHandler := NewAccountHandler(client)
	eventHandler := NewEventHandler(client)
	ticketHandler := NewTicketHandler(client)

	account := rg.Group("/accounts")
	account.POST("/create", accountHandler.CreateHandler)
	account.POST("/session", accountHandler.CreateSessionHandler)
	account.POST("/refresh", accountHandler.RefreshSessionHandler)
	account.GET("/:id", middleware.AuthMiddleware(), accountHandler.GetAccountHandler)
	account.DELETE("/delete-session", middleware.AuthMiddleware(), accountHandler.DeleteSessionHandler)

	event := rg.Group("/events")
	event.Use(middleware.AuthMiddleware())
	event.POST("/create", eventHandler.CreateEventHandler)
	event.GET("/:id", eventHandler.GetEventHandler)
	event.GET("/", eventHandler.ListEventsHandler)
	event.PATCH("/:id", eventHandler.UpdateEventHandler)
	event.DELETE("/:id", eventHandler.DeleteEventHandler)

	ticket := rg.Group("/tickets")
	ticket.Use(middleware.AuthMiddleware())
	ticket.POST("/create", ticketHandler.CreateTicketHandler)
	ticket.GET("/:id", ticketHandler.GetTicketHandler)
	ticket.GET("/", ticketHandler.ListTicketsHandler)
	ticket.PATCH("/:id", ticketHandler.UpdateTicketHandler)
}

func (h *httpServer) shutdown() error {
	// Shutdown HTTP server
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	h.mu.RLock()
	server := h.server
	h.mu.RUnlock()

	if server != nil {
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("HTTP server shutdown error: %v", err)
		}
	}

	// Close gRPC connection
	if h.conn != nil {
		if err := h.conn.Close(); err != nil {
			log.Printf("gRPC connection close error: %v", err)
		}
	}

	return nil
}

func NewHTTPServer(config *configs.Config, cache cache.RedisCacheService, token jwt.TokenService) HTTPServer {
	return &httpServer{
		config: config,
		cache:  cache,
		token:  token,
	}
}
