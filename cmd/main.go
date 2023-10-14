package main

import (
	common_crypto "github.com/Bit-Bridge-Source/BitBridge-CommonService-Go/public/crypto"
	"github.com/Bit-Bridge-Source/BitBridge-UserService-Go/internal/app"
	"github.com/Bit-Bridge-Source/BitBridge-UserService-Go/internal/database"
	fiberserver "github.com/Bit-Bridge-Source/BitBridge-UserService-Go/internal/fiber"
	grpcserver "github.com/Bit-Bridge-Source/BitBridge-UserService-Go/internal/grpc"
	"github.com/Bit-Bridge-Source/BitBridge-UserService-Go/internal/repository"
	"github.com/Bit-Bridge-Source/BitBridge-UserService-Go/internal/service"
	"github.com/gofiber/fiber/v2"
)

func main() {
	db, err := database.NewDatabase("mongodb://root:password@localhost:27017")
	if err != nil {
		panic(err)
	}

	if err := db.CreateIndexes(); err != nil {
		panic(err)
	}

	cryptoService := common_crypto.NewCrypto()

	defer db.Disconnect()

	// Initialize MongoDB adapter and repository
	mongoDBAdapter := repository.NewMongoAdapter(db.Collection)
	userRepository := repository.NewUserRepository(mongoDBAdapter)

	// Initialize user service and handler
	userService := service.NewUserService(userRepository, cryptoService)
	userHandler := fiberserver.NewUserFiberHandler(userService)

	// Initialize Fiber server
	fiberServer := fiberserver.NewUserFiberServer(fiber.Config{})
	fiberServer.SetupRoutes(userHandler)

	// Initialize gRPC server
	grpcServer := grpcserver.NewUserGrpcServer(userService)

	// Create app and add servers
	app := app.NewApp(fiberServer, grpcServer)
	app.Run(":3000", ":3001")
}
