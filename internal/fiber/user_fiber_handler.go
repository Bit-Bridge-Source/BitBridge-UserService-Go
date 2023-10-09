package fiberserver

import "github.com/Bit-Bridge-Source/BitBridge-UserService-Go/internal/service"

type UserFiberHandler struct {
	UserService service.IUserService
}

func NewUserFiberHandler(userService service.IUserService) *UserFiberHandler {
	return &UserFiberHandler{
		UserService: userService,
	}
}
