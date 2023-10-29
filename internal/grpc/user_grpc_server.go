package grpcserver

import (
	"context"
	"log"
	"net"

	common_grpc "github.com/Bit-Bridge-Source/BitBridge-CommonService-Go/public/grpc"
	"github.com/Bit-Bridge-Source/BitBridge-UserService-Go/internal/service"
	"github.com/Bit-Bridge-Source/BitBridge-UserService-Go/proto/pb"
	publicModel "github.com/Bit-Bridge-Source/BitBridge-UserService-Go/public/model"
	"google.golang.org/grpc"
)

type IUserGrpcServer interface {
	CreateUser(ctx context.Context, createUserModel *pb.CreateUserRequest) (*pb.PublicUserResponse, error)
	GetPrivateUserByIdentifier(ctx context.Context, getUserByIdentifierModel *pb.IdentifierRequest) (*pb.UserResponse, error)
	GetPublicUserByIdentifier(ctx context.Context, getPublicUserByIdentifierModel *pb.IdentifierRequest) (*pb.PublicUserResponse, error)
}

type UserGrpcServer struct {
	UserService  service.IUserService
	Interceptors []grpc.UnaryServerInterceptor
	pb.UnimplementedUserServiceServer
}

func NewUserGrpcServer(userService service.IUserService, interceptors []grpc.UnaryServerInterceptor) *UserGrpcServer {
	return &UserGrpcServer{
		UserService:  userService,
		Interceptors: interceptors,
	}
}

func (s *UserGrpcServer) Run(port string) error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen on port: %v", err)
		return err
	}
	gRPCServer := grpc.NewServer(
		grpc.UnaryInterceptor(common_grpc.ChainUnaryInterceptors(s.Interceptors...)),
	)
	pb.RegisterUserServiceServer(gRPCServer, s)
	return gRPCServer.Serve(lis)
}

func (s *UserGrpcServer) CreateUser(ctx context.Context, createUserModel *pb.CreateUserRequest) (*pb.PublicUserResponse, error) {
	createUserModelInternal := &publicModel.CreateUserModel{
		Email:    createUserModel.Email,
		Username: createUserModel.Username,
		Password: createUserModel.Password,
	}

	publicUserResponse, err := s.UserService.Create(ctx, createUserModelInternal)
	if err != nil {
		return nil, err
	}

	return &pb.PublicUserResponse{
		Id:        publicUserResponse.ID.Hex(),
		Username:  publicUserResponse.Username,
		CreatedAt: publicUserResponse.CreatedAt.String(),
		UpdatedAt: publicUserResponse.UpdatedAt.String(),
	}, nil
}

func (s *UserGrpcServer) GetPrivateUserByIdentifier(ctx context.Context, getUserByIdentifierModel *pb.IdentifierRequest) (*pb.UserResponse, error) {
	userResponse, err := s.UserService.FindByIdentifier(ctx, getUserByIdentifierModel.UserIdentifier)
	if err != nil {
		return nil, err
	}

	return &pb.UserResponse{
		Id:        userResponse.ID.Hex(),
		Email:     userResponse.Email,
		Username:  userResponse.Username,
		Hash:      userResponse.Hash,
		CreatedAt: userResponse.CreatedAt.String(),
		UpdatedAt: userResponse.UpdatedAt.String(),
	}, nil
}

func (s *UserGrpcServer) GetPublicUserByIdentifier(ctx context.Context, getPublicUserByIdentifierModel *pb.IdentifierRequest) (*pb.PublicUserResponse, error) {
	publicUserResponse, err := s.UserService.FindByIdentifier(ctx, getPublicUserByIdentifierModel.UserIdentifier)
	if err != nil {
		return nil, err
	}

	return &pb.PublicUserResponse{
		Id:        publicUserResponse.ID.Hex(),
		Username:  publicUserResponse.Username,
		CreatedAt: publicUserResponse.CreatedAt.String(),
		UpdatedAt: publicUserResponse.UpdatedAt.String(),
	}, nil
}
