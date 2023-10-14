package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	common_crypto "github.com/Bit-Bridge-Source/BitBridge-CommonService-Go/public/crypto"
	"github.com/Bit-Bridge-Source/BitBridge-UserService-Go/internal/model"
	"github.com/Bit-Bridge-Source/BitBridge-UserService-Go/internal/repository"
	"github.com/Bit-Bridge-Source/BitBridge-UserService-Go/internal/service"
	publicModel "github.com/Bit-Bridge-Source/BitBridge-UserService-Go/public/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockIUserRepository struct {
	mock.Mock
}

// Create implements repository.IUserRepository.
func (m *MockIUserRepository) Create(ctx context.Context, user *model.PrivateUserModel) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

// Delete implements repository.IUserRepository.
func (m *MockIUserRepository) Delete(ctx context.Context, user *model.PrivateUserModel) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

// FindByEmail implements repository.IUserRepository.
func (m *MockIUserRepository) FindByEmail(ctx context.Context, email string) (*model.PrivateUserModel, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	if user, ok := args.Get(0).(*model.PrivateUserModel); ok {
		return user, args.Error(1)
	}
	return nil, errors.New("type assertion to *model.PrivateUserModel failed")
}

// FindById implements repository.IUserRepository.
func (m *MockIUserRepository) FindById(ctx context.Context, id string) (*model.PrivateUserModel, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	if user, ok := args.Get(0).(*model.PrivateUserModel); ok {
		return user, args.Error(1)
	}
	return nil, errors.New("type assertion to *model.PrivateUserModel failed")
}

// FindByUsername implements repository.IUserRepository.
func (m *MockIUserRepository) FindByUsername(ctx context.Context, username string) (*model.PrivateUserModel, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	if user, ok := args.Get(0).(*model.PrivateUserModel); ok {
		return user, args.Error(1)
	}
	return nil, errors.New("type assertion to *model.PrivateUserModel failed")
}

func (m *MockIUserRepository) FindByIdentifier(ctx context.Context, identifier string) (*model.PrivateUserModel, error) {
	args := m.Called(ctx, identifier)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	if user, ok := args.Get(0).(*model.PrivateUserModel); ok {
		return user, args.Error(1)
	}
	return nil, errors.New("type assertion to *model.PrivateUserModel failed")
}

// Update implements repository.IUserRepository.
func (m *MockIUserRepository) Update(ctx context.Context, user *model.PrivateUserModel) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

// Ensure that MockIUserRepository implements IUserRepository.
var _ repository.IUserRepository = &MockIUserRepository{}

type MockCryptoService struct {
	mock.Mock
}

func (m *MockCryptoService) GenerateFromPassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockCryptoService) CompareHashAndPassword(hashedPassword, password string) error {
	args := m.Called(hashedPassword, password)
	return args.Error(0)
}

// Ensure that MockCryptoService implements ICryptoService.
var _ common_crypto.ICrypto = &MockCryptoService{}

func TestCreate_Success(t *testing.T) {
	// Arrange
	mockRepository := new(MockIUserRepository)
	mockCrypto := new(MockCryptoService)
	service := service.NewUserService(mockRepository, mockCrypto)
	ctx := context.Background()
	privateUser := &model.PrivateUserModel{
		ID:        primitive.NewObjectID(),
		Email:     "test@mail.com",
		Username:  "test",
		Hash:      "test",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	userToBeCreated := &publicModel.CreateUserModel{
		Email:    privateUser.Email,
		Username: privateUser.Username,
		Password: "test",
	}

	mockRepository.On("Create", ctx, mock.Anything).Return(nil)
	mockCrypto.On("GenerateFromPassword", mock.Anything).Return("hashedPassword", nil)

	// Act
	result, err := service.Create(ctx, userToBeCreated)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, privateUser.Email, result.Email)
	mockRepository.AssertExpectations(t)
}

func TestCreate_Failure(t *testing.T) {
	// Arrange
	mockRepository := new(MockIUserRepository)
	mockCrypto := new(MockCryptoService)
	service := service.NewUserService(mockRepository, mockCrypto)
	ctx := context.Background()
	privateUser := &model.PrivateUserModel{
		ID:        primitive.NewObjectID(),
		Email:     "test@mail.com",
		Username:  "test",
		Hash:      "test",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	userToBeCreated := &publicModel.CreateUserModel{
		Email:    privateUser.Email,
		Username: privateUser.Username,
		Password: "test",
	}

	mockRepository.On("Create", ctx, mock.Anything).Return(assert.AnError)
	mockCrypto.On("GenerateFromPassword", mock.Anything).Return("hashedPassword", nil)
	// Act
	result, err := service.Create(ctx, userToBeCreated)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepository.AssertExpectations(t)
}

func TestCreate_Failure_InvalidEncryption(t *testing.T) {
	// Arrange
	mockRepository := new(MockIUserRepository)
	mockCrypto := new(MockCryptoService)
	service := service.NewUserService(mockRepository, mockCrypto)
	ctx := context.Background()
	privateUser := &model.PrivateUserModel{
		ID:        primitive.NewObjectID(),
		Email:     "test@mail.com",
		Username:  "test",
		Hash:      "test",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	userToBeCreated := &publicModel.CreateUserModel{
		Email:    privateUser.Email,
		Username: privateUser.Username,
		Password: "test",
	}

	mockCrypto.On("GenerateFromPassword", mock.Anything).Return("", assert.AnError)

	// Act
	result, err := service.Create(ctx, userToBeCreated)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepository.AssertExpectations(t)
}

func TestFindByEmail_Success(t *testing.T) {
	// Arrange
	mockRepository := new(MockIUserRepository)
	mockCrypto := new(MockCryptoService)
	service := service.NewUserService(mockRepository, mockCrypto)
	ctx := context.Background()
	privateUser := &model.PrivateUserModel{
		ID:        primitive.NewObjectID(),
		Email:     "test@mail.com",
		Username:  "test",
		Hash:      "test",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepository.On("FindByEmail", ctx, privateUser.Email).Return(privateUser, nil)

	// Act
	result, err := service.FindByEmail(ctx, privateUser.Email)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, privateUser.Email, result.Email)
	mockRepository.AssertExpectations(t)
	mockCrypto.AssertExpectations(t)
}

func TestFindByEmail_Failure(t *testing.T) {
	// Arrange
	mockRepository := new(MockIUserRepository)
	mockCrypto := new(MockCryptoService)
	service := service.NewUserService(mockRepository, mockCrypto)
	ctx := context.Background()
	privateUser := &model.PrivateUserModel{
		ID:        primitive.NewObjectID(),
		Email:     "test@mail.com",
		Username:  "test",
		Hash:      "test",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepository.On("FindByEmail", ctx, privateUser.Email).Return(nil, assert.AnError)

	// Act
	result, err := service.FindByEmail(ctx, privateUser.Email)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepository.AssertExpectations(t)
}

func TestFindById_Success(t *testing.T) {
	// Arrange
	mockRepository := new(MockIUserRepository)
	mockCrypto := new(MockCryptoService)
	service := service.NewUserService(mockRepository, mockCrypto)
	ctx := context.Background()
	privateUser := &model.PrivateUserModel{
		ID:        primitive.NewObjectID(),
		Email:     "test@mail.com",
		Username:  "test",
		Hash:      "test",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepository.On("FindById", ctx, privateUser.ID.Hex()).Return(privateUser, nil)

	// Act
	result, err := service.FindById(ctx, privateUser.ID.Hex())

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, privateUser.ID, result.ID)
	mockRepository.AssertExpectations(t)
}

func TestFindById_Failure(t *testing.T) {
	// Arrange
	mockRepository := new(MockIUserRepository)
	mockCrypto := new(MockCryptoService)
	service := service.NewUserService(mockRepository, mockCrypto)
	ctx := context.Background()
	privateUser := &model.PrivateUserModel{
		ID:        primitive.NewObjectID(),
		Email:     "test@mail.com",
		Username:  "test",
		Hash:      "test",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepository.On("FindById", ctx, privateUser.ID.Hex()).Return(nil, assert.AnError)

	// Act
	result, err := service.FindById(ctx, privateUser.ID.Hex())

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepository.AssertExpectations(t)

}

func TestFindByUsername_Success(t *testing.T) {
	mockRepo := new(MockIUserRepository)
	mockCrypto := new(MockCryptoService)
	service := service.NewUserService(mockRepo, mockCrypto)

	ctx := context.Background()
	testUser := &model.PrivateUserModel{
		ID:        primitive.NewObjectID(),
		Email:     "test@mail.com",
		Username:  "test",
		Hash:      "test",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo.On("FindByUsername", ctx, testUser.Username).Return(testUser, nil)

	user, err := service.FindByUsername(ctx, testUser.Username)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, testUser.Username, user.Username)
	mockRepo.AssertExpectations(t)
}

func TestFindByUsername_Failure(t *testing.T) {
	mockRepo := new(MockIUserRepository)
	mockCrypto := new(MockCryptoService)
	service := service.NewUserService(mockRepo, mockCrypto)

	ctx := context.Background()
	testUser := &model.PrivateUserModel{
		ID:        primitive.NewObjectID(),
		Email:     "test@mail.com",
		Username:  "test",
		Hash:      "test",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo.On("FindByUsername", ctx, testUser.Username).Return(nil, assert.AnError)

	user, err := service.FindByUsername(ctx, testUser.Username)

	assert.Error(t, err)
	assert.Nil(t, user)
	mockRepo.AssertExpectations(t)
}

func TestFindByIdentifier(t *testing.T) {
	mockRepo := new(MockIUserRepository)
	mockCrypto := new(MockCryptoService)
	service := service.NewUserService(mockRepo, mockCrypto)

	ctx := context.Background()
	testUser := &model.PrivateUserModel{
		ID:        primitive.NewObjectID(),
		Email:     "test@mail.com",
		Username:  "test",
		Hash:      "test",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Test with ObjectID
	mockRepo.On("FindById", ctx, testUser.ID.Hex()).Return(testUser, nil)
	user, err := service.FindByIdentifier(ctx, testUser.ID.Hex())
	assert.NoError(t, err)
	assert.Equal(t, testUser.ID, user.ID)
	mockRepo.AssertExpectations(t)

	// Test with Email
	mockRepo.On("FindByEmail", ctx, "test@email.com").Return(testUser, nil)
	user, err = service.FindByIdentifier(ctx, "test@email.com")
	assert.NoError(t, err)
	assert.Equal(t, testUser.ID, user.ID)
	mockRepo.AssertExpectations(t)

	// Test with Username
	mockRepo.On("FindByUsername", ctx, "username").Return(testUser, nil)
	user, err = service.FindByIdentifier(ctx, "username")
	assert.NoError(t, err)
	assert.Equal(t, testUser.ID, user.ID)
	mockRepo.AssertExpectations(t)
}
