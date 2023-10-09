package fiberserver

import (
	"github.com/Bit-Bridge-Source/BitBridge-UserService-Go/internal/service"
	publicModel "github.com/Bit-Bridge-Source/BitBridge-UserService-Go/public/model"
	"github.com/gofiber/fiber/v2"
)

type UserFiberHandler struct {
	UserService service.IUserService
}

func NewUserFiberHandler(userService service.IUserService) *UserFiberHandler {
	return &UserFiberHandler{
		UserService: userService,
	}
}

func (handler *UserFiberHandler) Create(c *fiber.Ctx) error {
	var createUserModel publicModel.CreateUserModel
	if err := c.BodyParser(&createUserModel); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	user, err := handler.UserService.Create(c.Context(), &createUserModel)

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Something went wrong")
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

func (handler *UserFiberHandler) FindByIdentifierPublic(c *fiber.Ctx) error {
	userIdentifier := c.Params("user_identifier")

	user, err := handler.UserService.FindByIdentifier(c.Context(), userIdentifier)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Something went wrong")
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

func (handler *UserFiberHandler) FindByIdentifierPrivate(c *fiber.Ctx) error {
	userIdentifier := c.Params("user_identifier")

	// Check if the user is authorized to access this resource
	user_id, ok := c.Locals("user_id").(string)
	if !ok || user_id == "" || user_id != "-1" {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorizedd")
	}

	user, err := handler.UserService.FindByIdentifier(c.Context(), userIdentifier)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Something went wrong")
	}

	return c.Status(fiber.StatusOK).JSON(user)
}
