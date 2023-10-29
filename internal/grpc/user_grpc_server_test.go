package grpcserver_test

import (
	"context"
	"testing"
	"time"

	grpcserver "github.com/Bit-Bridge-Source/BitBridge-UserService-Go/internal/grpc"
	privateModel "github.com/Bit-Bridge-Source/BitBridge-UserService-Go/internal/model"
	"github.com/Bit-Bridge-Source/BitBridge-UserService-Go/internal/service"
	"github.com/Bit-Bridge-Source/BitBridge-UserService-Go/proto/pb"
	publicModel "github.com/Bit-Bridge-Source/BitBridge-UserService-Go/public/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
)

type MockIUserService struct {
	mock.Mock
}

// Create implements service.IUserService.
func (m *MockIUserService) Create(ctx context.Context, user *publicModel.CreateUserModel) (*privateModel.PrivateUserModel, error) {
	args := m.Called(ctx, user)

	userModelArgs, ok := args.Get(0).(*privateModel.PrivateUserModel)
	if !ok && args.Get(0) != nil {
		return nil, args.Error(1)
	}

	return userModelArgs, args.Error(1)
}

// FindByEmail implements service.IUserService.
func (m *MockIUserService) FindByEmail(ctx context.Context, email string) (*privateModel.PrivateUserModel, error) {
	args := m.Called(ctx, email)

	userModelArgs, ok := args.Get(0).(*privateModel.PrivateUserModel)
	if !ok && args.Get(0) != nil {
		return nil, args.Error(1)
	}

	return userModelArgs, args.Error(1)
}

// FindById implements service.IUserService.
func (m *MockIUserService) FindById(ctx context.Context, id string) (*privateModel.PrivateUserModel, error) {
	args := m.Called(ctx, id)

	userModelArgs, ok := args.Get(0).(*privateModel.PrivateUserModel)
	if !ok && args.Get(0) != nil {
		return nil, args.Error(1)
	}

	return userModelArgs, args.Error(1)
}

// FindByIdentifier implements service.IUserService.
func (m *MockIUserService) FindByIdentifier(ctx context.Context, identifier string) (*privateModel.PrivateUserModel, error) {
	args := m.Called(ctx, identifier)

	userModelArgs, ok := args.Get(0).(*privateModel.PrivateUserModel)
	if !ok && args.Get(0) != nil {
		return nil, args.Error(1)
	}

	return userModelArgs, args.Error(1)
}

// FindByUsername implements service.IUserService.
func (m *MockIUserService) FindByUsername(ctx context.Context, username string) (*privateModel.PrivateUserModel, error) {
	args := m.Called(ctx, username)

	userModelArgs, ok := args.Get(0).(*privateModel.PrivateUserModel)
	if !ok && args.Get(0) != nil {
		return nil, args.Error(1)
	}

	return userModelArgs, args.Error(1)
}

// Ensure that the mock implements the interface
var _ service.IUserService = &MockIUserService{}

func TestCreateUser_Success(t *testing.T) {
	userResponse := &privateModel.PrivateUserModel{
		ID:        primitive.NewObjectID(),
		Email:     "test@mail.com",
		Username:  "test",
		Hash:      "test",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Setup
	mockUserService := new(MockIUserService)
	mockUserService.On("Create", mock.Anything, mock.Anything).Return(userResponse, nil)
	grpcserver := grpcserver.NewUserGrpcServer(mockUserService, []grpc.UnaryServerInterceptor{})

	// Test
	resp, err := grpcserver.CreateUser(context.Background(), &pb.CreateUserRequest{
		Email:    userResponse.Email,
		Username: userResponse.Username,
		Password: userResponse.Hash,
	})

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, userResponse.ID.Hex(), resp.Id)
	mockUserService.AssertExpectations(t)
}

func TestCreateUser_Fail(t *testing.T) {
	mockUserService := new(MockIUserService)
	mockUserService.On("Create", mock.Anything, mock.Anything).Return(nil, assert.AnError)
	grpcserver := grpcserver.NewUserGrpcServer(mockUserService, []grpc.UnaryServerInterceptor{})

	// Test
	resp, err := grpcserver.CreateUser(context.Background(), &pb.CreateUserRequest{
		Email:    "test@mail.com",
		Username: "test",
		Password: "test",
	})

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resp)
	mockUserService.AssertExpectations(t)
}

func TestGetPrivateUserByIdentifier_Success(t *testing.T) {
	userResponse := &privateModel.PrivateUserModel{
		ID:        primitive.NewObjectID(),
		Email:     "test@mail.com",
		Username:  "test",
		Hash:      "test",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Setup
	mockUserService := new(MockIUserService)
	mockUserService.On("FindByIdentifier", mock.Anything, mock.Anything).Return(userResponse, nil)
	grpcserver := grpcserver.NewUserGrpcServer(mockUserService, []grpc.UnaryServerInterceptor{})

	// Test
	resp, err := grpcserver.GetPrivateUserByIdentifier(context.Background(), &pb.IdentifierRequest{
		UserIdentifier: userResponse.ID.Hex(),
	})

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, userResponse.ID.Hex(), resp.Id)
	mockUserService.AssertExpectations(t)
}

func TestGetPrivateUserByIdentifier_Fail(t *testing.T) {
	mockUserService := new(MockIUserService)
	mockUserService.On("FindByIdentifier", mock.Anything, mock.Anything).Return(nil, assert.AnError)
	grpcserver := grpcserver.NewUserGrpcServer(mockUserService, []grpc.UnaryServerInterceptor{})

	// Test
	resp, err := grpcserver.GetPrivateUserByIdentifier(context.Background(), &pb.IdentifierRequest{
		UserIdentifier: "test",
	})

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resp)
	mockUserService.AssertExpectations(t)
}

func TestGetPublicUserByIdentifier_Success(t *testing.T) {
	userResponse := &privateModel.PrivateUserModel{
		ID:        primitive.NewObjectID(),
		Email:     "test@mail.com",
		Username:  "test",
		Hash:      "test",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Setup
	mockUserService := new(MockIUserService)
	mockUserService.On("FindByIdentifier", mock.Anything, mock.Anything).Return(userResponse, nil)
	grpcserver := grpcserver.NewUserGrpcServer(mockUserService, []grpc.UnaryServerInterceptor{})

	// Test
	resp, err := grpcserver.GetPublicUserByIdentifier(context.Background(), &pb.IdentifierRequest{
		UserIdentifier: userResponse.ID.Hex(),
	})

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, userResponse.ID.Hex(), resp.Id)
	mockUserService.AssertExpectations(t)
}

func TestGetPublicUserByIdentifier_Fail(t *testing.T) {
	mockUserService := new(MockIUserService)
	mockUserService.On("FindByIdentifier", mock.Anything, mock.Anything).Return(nil, assert.AnError)
	grpcserver := grpcserver.NewUserGrpcServer(mockUserService, []grpc.UnaryServerInterceptor{})

	// Test
	resp, err := grpcserver.GetPublicUserByIdentifier(context.Background(), &pb.IdentifierRequest{
		UserIdentifier: "test",
	})

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resp)
	mockUserService.AssertExpectations(t)
}
