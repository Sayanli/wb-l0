package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) FindByUid(c *fiber.Ctx) error {
	m := c.Queries()
	uid := m["uid"]
	fmt.Println(uid)
	order, err := h.services.Order.FindByUid(context.Background(), uid)
	if err != nil {
		return c.Render("static/error_order.html", fiber.Map{})
	}
	jsonBytes, err := json.Marshal(&order)
	return c.Render("static/order.html", fiber.Map{
		"uid":   "uid",
		"order": string(jsonBytes),
	})
}

func (h *Handler) FindAll(c *fiber.Ctx) error {
	orders, err := h.services.Order.FindAll(context.Background())
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": err.Error(), "data": nil})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Orders found", "data": orders})
}
