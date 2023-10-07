package service

import (
	"context"
	"net/mail"

	"github.com/Bit-Bridge-Source/BitBridge-UserService-Go/internal/model"
	"github.com/Bit-Bridge-Source/BitBridge-UserService-Go/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IUserService interface {
	Create(ctx context.Context, user *model.PrivateUserModel) (*model.PrivateUserModel, error)
	FindById(ctx context.Context, id string) (*model.PrivateUserModel, error)
	FindByEmail(ctx context.Context, email string) (*model.PrivateUserModel, error)
	FindByUsername(ctx context.Context, username string) (*model.PrivateUserModel, error)
	FindByIdentifier(ctx context.Context, identifier string) (*model.PrivateUserModel, error)
}

type UserService struct {
	Repository repository.IUserRepository
}

// NewUserService creates a new instance of UserService.
func NewUserService(repository repository.IUserRepository) *UserService {
	return &UserService{
		Repository: repository,
	}
}

// Create implements IUserService.
func (s *UserService) Create(ctx context.Context, user *model.PrivateUserModel) (*model.PrivateUserModel, error) {
	err := s.Repository.Create(ctx, user)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// FindByEmail implements IUserService.
func (s *UserService) FindByEmail(ctx context.Context, email string) (*model.PrivateUserModel, error) {
	privateUser, err := s.Repository.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return privateUser, nil
}

// FindById implements IUserService.
func (s *UserService) FindById(ctx context.Context, id string) (*model.PrivateUserModel, error) {
	privateUser, err := s.Repository.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	return privateUser, nil
}

// FindByIdentifier implements IUserService.
func (s *UserService) FindByIdentifier(ctx context.Context, identifier string) (*model.PrivateUserModel, error) {
	_, err := mail.ParseAddress(identifier)
	if err == nil {
		return s.Repository.FindByEmail(ctx, identifier)
	}

	_, err = primitive.ObjectIDFromHex(identifier)
	if err == nil {
		return s.Repository.FindById(ctx, identifier)
	}

	return s.Repository.FindByUsername(ctx, identifier)
}

// FindByUsername implements IUserService.
func (s *UserService) FindByUsername(ctx context.Context, username string) (*model.PrivateUserModel, error) {
	privateUser, err := s.Repository.FindByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	return privateUser, nil
}

// Ensure UserService implements IUserService
var _ IUserService = &UserService{}
