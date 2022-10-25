package rest

import (
	"context"
	"net/http"
	"strconv"

	"github.com/TudorHulban/rest-articles/app/service"
	"github.com/gofiber/fiber/v2"
)

func (s *WebServer) handleAlive() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).JSON(&fiber.Map{
			"alive": true,
		})
	}
}

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
	return func(c *fiber.Ctx) error {
		idRequest := c.Params("id")
		idItem, errReq := strconv.Atoi(idRequest)
		if errReq != nil {
			return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
				"success": false,
				"error":   errReq.Error(),
			})
		}

		if idItem < 1 {
			return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
				"success": false,
				"error":   "ID should be at least 1",
			})
		}

		reconstructedItem, errFetch := s.serv.GetArticle(context.Background(), idItem)
		if errFetch != nil {
			return c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
				"success": false,
				"error":   errFetch.Error(),
			})
		}

		return c.Status(http.StatusOK).JSON(&fiber.Map{
			"success": true,
			"company": reconstructedItem,
		})
	}
}

// handleGetArticles - no pagination
func (s *WebServer) handleGetArticles() fiber.Handler {
	return func(c *fiber.Ctx) error { return nil }
}
