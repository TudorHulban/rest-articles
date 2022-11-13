package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/TudorHulban/rest-articles/app/apperrors"
	"github.com/TudorHulban/rest-articles/app/service"
	"github.com/gofiber/fiber/v2"
)

func (rest *Rest) HandlerAlive() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).JSON(&fiber.Map{
			"alive": true,
		})
	}
}

func (rest *Rest) HandlerNewArticle() fiber.Handler {
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

		idInsert, errIns := rest.serv.CreateArticle(c.Context(), &service.ParamsCreateArticle{
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

func (rest *Rest) HandlerGetArticle() fiber.Handler {
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

		reconstructedItem, errFetch := rest.serv.GetArticle(c.Context(), int64(idItem))
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
func (rest *Rest) HandlerGetArticles() fiber.Handler {
	return func(c *fiber.Ctx) error {
		items, errReq := rest.serv.GetArticles(c.Context())
		if errReq != nil {
			return c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
				"success": false,
				"error":   errReq.Error(),
			})
		}

		return c.Status(http.StatusOK).JSON(&fiber.Map{
			"success":  true,
			"articles": items,
		})
	}
}

// handleGetArticlesWithPagination
func (rest *Rest) HandlerGetArticlesWithPagination() fiber.Handler {
	return func(c *fiber.Ctx) error {
		limit := c.Query("limit")
		page := c.Query("page")

		nLimit, _ := strconv.Atoi(limit)
		nPage, _ := strconv.Atoi(page)

		paginatedItems, errReq := rest.serv.GetArticlesPaginated(c.Context(), nLimit, nPage)
		if errReq != nil {
			return c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
				"success": false,
				"error":   errReq.Error(),
			})
		}

		return c.Status(http.StatusOK).JSON(&fiber.Map{
			"success":  true,
			"articles": paginatedItems,
		})
	}
}

func (rest *Rest) HandlerUpdateArticle() fiber.Handler {
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

		if errUpd := rest.serv.UpdateArticle(c.Context(), &params); errUpd != nil {
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

func (rest *Rest) HandlerDeleteArticle() fiber.Handler {
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

		if errDel := rest.serv.DeleteArticle(c.Context(), int64(idItem)); errDel != nil {
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
