package fiberserver

import "github.com/gofiber/fiber/v2"

type UserFiberServer struct {
	App *fiber.App
}

func (server *UserFiberServer) Run(port string) error {
	return server.App.Listen(port)
}
