package main

import (
	"context"
	"log"
	"net"

	cfg "github.com/Denis/project_auth/internal/config"
	repoAuth "github.com/Denis/project_auth/internal/repository/auth"
	serviceAuth "github.com/Denis/project_auth/internal/service/auth"
	desc "github.com/Denis/project_auth/pkg/user_v1"

	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	ctx := context.Background()

	// Конфигурация
	err := cfg.Load(".env")
	if err != nil {
		log.Fatal("failed to load config:", err)
	}

	grpcConfig, err := cfg.NewGRPCConfig()
	if err != nil {
		log.Fatal("failed to get grpc config:", err)
	}

	pgConfig, err := cfg.NewPGConfig()
	if err != nil {
		log.Fatal("failed to get pg config:", err)
	}

	// Пул соединений с БД
	pool, err := pgxpool.Connect(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}
	defer pool.Close()

	// Создаем репозиторий
	authRepo := repoAuth.NewRepository(pool)

	// Создаем сервис (он уже реализует desc.UserV1Server)
	authService := serviceAuth.NewService(authRepo)

	// Создаем gRPC сервер
	lis, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatal("failed to listen:", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)

	// Регистрируем сервис напрямую (без дополнительного wrapper)
	desc.RegisterUserV1Server(s, authService)

	log.Println("server listening at", lis.Addr())
	if err = s.Serve(lis); err != nil {
		log.Fatal("failed to serve:", err)
	}
}
