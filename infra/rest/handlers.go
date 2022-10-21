package rest

import (
	"context"
	"net/http"

	"github.com/TudorHulban/rest-articles/app/service"
	"github.com/gofiber/fiber/v2"
)

func (s *WebServer) handleNewArticle() fiber.Handler {
	return func(c *fiber.Ctx) error {
		type request struct {
			Title string `json:"title"`
			URL   string `json:"url"`
		}

		var req request

		if errBody := c.BodyParser(&req); errBody != nil {
			return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
				"success": false,
				"error":   errBody.Error(),
			})
		}

		idInsert, errIns := s.serv.CreateArticle(context.Background(), &service.ParamsCreateArticle{
			Title: req.Title,
			URL:   req.URL,
		})
		if errIns != nil {
			return c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
				"success": false,
				"error":   errIns.Error(),
			})
		}

		return c.Status(http.StatusOK).JSON(&fiber.Map{
			"success": true,
			"id":      idInsert,
		})
	}
}

func (s *WebServer) handleGetArticle() fiber.Handler {
	return func(c *fiber.Ctx) error { return nil }
}

func (s *WebServer) handleGetArticles() fiber.Handler {
	return func(c *fiber.Ctx) error { return nil }
}
