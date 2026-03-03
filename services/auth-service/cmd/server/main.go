package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/Saad7890-web/FluxGuard/internal/config"
	"github.com/Saad7890-web/FluxGuard/internal/handler"
	"github.com/Saad7890-web/FluxGuard/internal/repository"
	"github.com/Saad7890-web/FluxGuard/internal/service"
	"github.com/Saad7890-web/FluxGuard/internal/token"
	authv1 "github.com/Saad7890-web/FluxGuard/proto/auth/v1"
)

func main() {
	cfg := config.Load()

	db, err := config.NewDB(cfg.DBUrl)
	if err != nil {
		log.Fatal(err)
	}

	if err := config.RunMigrations(cfg.DBUrl, cfg.MigrationsDir); err != nil {
		log.Fatal(err)
	}

	repo := repository.NewUserRepository(db)
	tokens := token.NewManager(cfg.JWTSecret)
	svc := service.NewAuthService(repo, tokens)
	handler := handler.NewGRPCHandler(svc)

	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatal(err)
	}

	server := grpc.NewServer()
	authv1.RegisterAuthServiceServer(server, handler)
	reflection.Register(server)

	go func() {
		log.Println("Auth gRPC running on port", cfg.GRPCPort)
		if err := server.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down...")
	server.GracefulStop()
}