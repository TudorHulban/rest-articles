package rest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/TudorHulban/rest-articles/app/apperrors"
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

		reconstructedItem, errFetch := s.serv.GetArticle(context.Background(), int64(idItem))
		if errFetch != nil {
			if errors.Is(errFetch, apperrors.ErrObjectNotFound{}) {
				return c.Status(http.StatusOK).JSON(&fiber.Map{
					"success": true,
					"error":   errFetch.Error(),
				})
			}

			return c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
				"success": false,
				"error":   errFetch.Error(),
			})
		}

		return c.Status(http.StatusOK).JSON(&fiber.Map{
			"success": true,
			"article": reconstructedItem,
		})
	}
}

// handleGetArticles - no pagination
func (s *WebServer) handleGetArticles() fiber.Handler {
	return func(c *fiber.Ctx) error {
		articles, errReq := s.serv.GetArticles(context.Background())
		if errReq != nil {
			return c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
				"success": false,
				"error":   errReq.Error(),
			})
		}

		return c.Status(http.StatusOK).JSON(&fiber.Map{
			"success":  true,
			"articles": articles,
		})
	}
}

func (s *WebServer) handleUpdateArticle() fiber.Handler {
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

		type request struct {
			Title string `json:"title"`
			URL   string `json:"url"`
		}

		var req request

		if errUn := json.Unmarshal(c.Body(), &req); errUn != nil {
			return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
				"success":     false,
				"error":       fmt.Sprintf("for ID: %d: %s", idItem, errUn.Error()),
				"requestbody": string(c.Body()),
			})
		}

		// TODO: investigate why unit test did not work with body parser
		// if errBody := c.BodyParser(&req); errBody != nil {
		// 	return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
		// 		"success":     false,
		// 		"error":       fmt.Sprintf("for ID: %d: %s", idItem, errBody.Error()),
		// 		"requestbody": string(c.Body()),
		// 	})
		// }

		params := service.ParamsUpdateArticle{
			ID:    int64(idItem),
			Title: &req.Title,
			URL:   &req.URL,
		}

		if errUpd := s.serv.UpdateArticle(context.Background(), &params); errUpd != nil {
			return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
				"success": false,
				"error":   errUpd.Error(),
			})
		}

		return c.Status(http.StatusOK).JSON(&fiber.Map{
			"success": true,
		})
	}
}

func (s *WebServer) handleDeleteArticle() fiber.Handler {
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

		if errDel := s.serv.DeleteArticle(context.Background(), int64(idItem)); errDel != nil {
			return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
				"success": false,
				"error":   errDel.Error(),
			})
		}

		return c.Status(http.StatusOK).JSON(&fiber.Map{
			"success": true,
		})
	}
}
