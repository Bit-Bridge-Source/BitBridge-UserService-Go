package fiberserver

import (
	"github.com/Bit-Bridge-Source/BitBridge-UserService-Go/internal/service"
)

type UserFiberHandler struct {
	UserService service.IUserService
}

func NewUserFiberHandler(userService service.IUserService) *UserFiberHandler {
	return &UserFiberHandler{
		UserService: userService,
	}
}

// func (handler *UserFiberHandler) Create(c *fiber.Ctx) {
// 	var createUserModel publicModel.CreateUserModel
// 	if err := c.BodyParser(&createUserModel); err != nil {
// 		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
// 	}

// 	// TODO: Add validation

// 	user, err := handler.UserService.Create(createUserModel)
// }
