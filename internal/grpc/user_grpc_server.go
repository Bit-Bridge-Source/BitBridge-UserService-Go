package grpcserver

import (
	"context"

	"github.com/Bit-Bridge-Source/BitBridge-UserService-Go/internal/service"
	"github.com/Bit-Bridge-Source/BitBridge-UserService-Go/proto/pb"
	publicModel "github.com/Bit-Bridge-Source/BitBridge-UserService-Go/public/model"
)

type IUserGrpcServer interface {
	CreateUser(ctx context.Context, createUserModel *pb.CreateUserRequest) (*pb.PublicUserResponse, error)
	GetPrivateUserByIdentifier(ctx context.Context, getUserByIdentifierModel *pb.IdentifierRequest) (*pb.UserResponse, error)
	GetPublicUserByIdentifier(ctx context.Context, getPublicUserByIdentifierModel *pb.IdentifierRequest) (*pb.PublicUserResponse, error)
}

type UserGrpcServer struct {
	UserService service.IUserService
	pb.UnimplementedUserServiceServer
}

func NewUserGrpcServer(userService service.IUserService) *UserGrpcServer {
	return &UserGrpcServer{
		UserService: userService,
	}
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
