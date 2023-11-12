package router

import (
	"wb-l0/internal/controller/handler"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type Server struct {
	app     *fiber.App
	handler *handler.Handler
}

func NewServer(app *fiber.App, handler *handler.Handler) *Server {
	return &Server{app: app, handler: handler}
}

func (s *Server) Router() {
	api := s.app.Group("/api", logger.New())
	api.Get("/", func(c *fiber.Ctx) error {
		return c.Render("static/index.html", fiber.Map{})
	})

	order := api.Group("/order")

	//Order
	order.Get("/all", s.handler.FindAll)
	order.Get("/", s.handler.FindByUid)
}
