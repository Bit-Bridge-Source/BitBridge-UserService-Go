package fiberserver

import "github.com/gofiber/fiber/v2"

type UserFiberServer struct {
	App *fiber.App
}

func NewUserFiberServer(config fiber.Config) *UserFiberServer {
	return &UserFiberServer{
		App: fiber.New(config),
	}
}

func (server *UserFiberServer) Run(port string) error {
	return server.App.Listen(port)
}

func (server *UserFiberServer) Stop() error {
	return server.App.Shutdown()
}

func (server *UserFiberServer) SetupRoutes(userHandler *UserFiberHandler) {
	private := server.App.Group("/private")
	private.Get("/user/:user_identifier", userHandler.FindByIdentifierPrivate)
	private.Post("/user", userHandler.Create)

	server.App.Get("/user/:user_identifier", userHandler.FindByIdentifierPublic)
}
