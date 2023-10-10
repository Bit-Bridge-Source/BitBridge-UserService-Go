package app

import (
	"log"

	fiberserver "github.com/Bit-Bridge-Source/BitBridge-UserService-Go/internal/fiber"
	grpcserver "github.com/Bit-Bridge-Source/BitBridge-UserService-Go/internal/grpc"
)

type App struct {
	FiberServer *fiberserver.UserFiberServer
	GRPCServer  *grpcserver.UserGrpcServer
}

func NewApp(fiberServer *fiberserver.UserFiberServer, gRPCServer *grpcserver.UserGrpcServer) *App {
	return &App{
		FiberServer: fiberServer,
		GRPCServer:  gRPCServer,
	}
}

func (app *App) Run(httpPort string, gRPCPort string) {
	// Channels to collect errors from the servers
	httpErrChan := make(chan error)
	grpcErrChan := make(chan error)

	// Run Fiber server
	go func() {
		log.Println("Starting Fiber server on port", httpPort)
		httpErrChan <- app.FiberServer.Run(httpPort)
	}()

	// Run gRPC server
	go func() {
		log.Println("Starting gRPC server on port", gRPCPort)
		grpcErrChan <- app.GRPCServer.Run(gRPCPort)
	}()

	// Wait for errors from the servers
	select {
	case err := <-httpErrChan:
		log.Fatal("Fiber server error: ", err)
	case err := <-grpcErrChan:
		log.Fatal("gRPC server error: ", err)
	}
}
