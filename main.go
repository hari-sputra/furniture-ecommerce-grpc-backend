package main

import (
	"context"
	"log"
	"net"
	"os"
	"time"

	"github.com/hari-sputra/furniture-ecommerce-grpc-be/internal/handler"
	"github.com/hari-sputra/furniture-ecommerce-grpc-be/internal/middleware"
	"github.com/hari-sputra/furniture-ecommerce-grpc-be/internal/repository"
	"github.com/hari-sputra/furniture-ecommerce-grpc-be/internal/service"
	"github.com/hari-sputra/furniture-ecommerce-grpc-be/pb/auth"
	"github.com/hari-sputra/furniture-ecommerce-grpc-be/pkg/database"
	"github.com/joho/godotenv"
	gocache "github.com/patrickmn/go-cache"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	ctx := context.Background()

	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db := database.InitDB(ctx)
	defer db.Close()

	lis, err := net.Listen("tcp", ":50051")

	if err != nil {
		log.Panicf("Error when listening %v", err)
	}

	// cache
	cache := gocache.New(time.Hour*24, time.Hour)

	// middleware
	authMiddleware := middleware.NewAuthMiddleware(cache)

	// auth
	authRepository := repository.NewAuthRepository(db)
	authService := service.NewAuthService(authRepository, cache)
	authHandler := handler.NewAuthHandler(authService)

	serv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.ErrorMiddleware,
			authMiddleware.Middleware,
		),
	)

	auth.RegisterAuthServiceServer(serv, authHandler)

	if os.Getenv("ENV") == "dev" {
		reflection.Register(serv)
		println("reflection registered")
	}

	log.Println("Server is running on port 50051")
	if err := serv.Serve(lis); err != nil {
		log.Panicf("Error when serving %v", err)
	}
}
