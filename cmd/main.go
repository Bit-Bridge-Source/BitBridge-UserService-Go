package main

import (
	common_crypto "github.com/Bit-Bridge-Source/BitBridge-CommonService-Go/public/crypto"
	common_fiber "github.com/Bit-Bridge-Source/BitBridge-CommonService-Go/public/fiber"
	common_grpc "github.com/Bit-Bridge-Source/BitBridge-CommonService-Go/public/grpc"
	common_vault "github.com/Bit-Bridge-Source/BitBridge-CommonService-Go/public/vault"
	"github.com/Bit-Bridge-Source/BitBridge-UserService-Go/internal/app"
	"github.com/Bit-Bridge-Source/BitBridge-UserService-Go/internal/database"
	fiberserver "github.com/Bit-Bridge-Source/BitBridge-UserService-Go/internal/fiber"
	grpcserver "github.com/Bit-Bridge-Source/BitBridge-UserService-Go/internal/grpc"
	"github.com/Bit-Bridge-Source/BitBridge-UserService-Go/internal/repository"
	"github.com/Bit-Bridge-Source/BitBridge-UserService-Go/internal/service"
	"github.com/Bit-Bridge-Source/BitBridge-UserService-Go/proto/pb"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
)

func main() {
	// Initialize MongoDB database
	db, err := database.NewDatabase("mongodb://root:password@localhost:27017")
	if err != nil {
		panic(err)
	}
	defer db.Disconnect()
	if err := db.CreateIndexes(); err != nil {
		panic(err)
	}

	// Initialize Vault client and read secret for authentication
	vaultClient, err := common_vault.NewVault("http://127.0.0.1:8200", "XZ5!Ojk88#Ox8PoM!yZhiJfHs")
	if err != nil {
		panic(err)
	}
	vaultSecret, err := vaultClient.ReadSecret("secret/data/jwt_secret")
	if err != nil {
		panic(err)
	}

	// Initialize MongoDB adapter and repository
	mongoDBAdapter := repository.NewMongoAdapter(db.Collection)
	userRepository := repository.NewUserRepository(mongoDBAdapter)

	// Initialize user service and handler
	cryptoService := common_crypto.NewCrypto()
	userService := service.NewUserService(userRepository, cryptoService)
	userHandler := fiberserver.NewUserFiberHandler(userService)

	// Initialize Fiber server
	fiberServer := fiberserver.NewUserFiberServer(fiber.Config{
		ErrorHandler: common_fiber.FiberErrorHandler,
	})
	fiberServer.SetupRoutes(userHandler, common_fiber.FiberJWTAuthenticator(vaultSecret))

	// Initialize gRPC server
	var publicMethods = map[string]struct{}{pb.UserService_GetPublicUserByIdentifier_FullMethodName: {}}
	authInterceptor := common_grpc.AuthUnaryInterceptor(vaultSecret, publicMethods)
	errorInterceptor := common_grpc.GRPCErrorHandler
	grpcServer := grpcserver.NewUserGrpcServer(userService, []grpc.UnaryServerInterceptor{authInterceptor, errorInterceptor})

	// Create app and add servers
	app := app.NewApp(fiberServer, grpcServer)
	app.Run(":3000", ":3001")
}
