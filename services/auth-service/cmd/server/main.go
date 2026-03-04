package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/Saad7890-web/FluxGuard/internal/config"
	"github.com/Saad7890-web/FluxGuard/internal/handler"
	"github.com/Saad7890-web/FluxGuard/internal/repository"
	"github.com/Saad7890-web/FluxGuard/internal/service"
	"github.com/Saad7890-web/FluxGuard/internal/token"
	"github.com/Saad7890-web/FluxGuard/pkg/logger"
	authv1 "github.com/Saad7890-web/FluxGuard/proto/auth/v1"
	"go.uber.org/zap"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg := config.Load()

	// Run migrations FIRST
	if err := config.RunMigrations(cfg.DBUrl, cfg.MigrationsDir); err != nil {
		log.Fatal(err)
	}

	db, err := config.NewDB(cfg.DBUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	zapLogger, _ := logger.New()
	defer zapLogger.Sync()

	repo := repository.NewUserRepository(db)
	tokens := token.NewManager(cfg.JWTSecret)
	svc := service.NewAuthService(repo, tokens)
	grpcHandler := handler.NewGRPCHandler(svc)

	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatal(err)
	}

	server := grpc.NewServer(
		grpc.UnaryInterceptor(handler.UnaryLoggingInterceptor(zapLogger)),
	)

	authv1.RegisterAuthServiceServer(server, grpcHandler)
	reflection.Register(server)

	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(server, healthServer)
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	go func() {
		zapLogger.Info("auth service started", zap.String("port", cfg.GRPCPort))
		if err := server.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	zapLogger.Info("shutting down auth service")
	server.GracefulStop()
}