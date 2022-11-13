package web

import (
	"fmt"
	"strconv"

	"github.com/TudorHulban/rest-articles/infra/rest"
	"github.com/gofiber/fiber/v2"
)

type WebServer struct {
	App  *fiber.App
	rest *rest.Rest

	errShutdown error
	port        uint
}

func NewWebServer(port uint, rest *rest.Rest) (*WebServer, error) {
	return &WebServer{
		App:  fiber.New(),
		port: port,
		rest: rest,
	}, nil
}

func (s *WebServer) Start() {
	s.AddRoutes()

	fmt.Println("web server started")

	s.errShutdown = s.App.Listen(":" + strconv.Itoa(int(s.port)))

	if s.errShutdown != nil {
		fmt.Printf("start - stopped now: %s\n", s.errShutdown)

		return
	}

	fmt.Println("start - stopped now: no error")
}

// Stop relases web server and service.
// Returns closing errors from web and service.
func (s *WebServer) Stop() (error, error) {
	fmt.Println("stopping Fiber")

	var errWebClose, errServiceClose error

	if errShut := s.App.Shutdown(); errShut != nil {
		fmt.Printf("error Fiber: %s\n", errShut.Error()) // local handling can be improved

		errWebClose = errShut
	}

	return errWebClose, errServiceClose
}

func (s *WebServer) AddRoutes() {
	s.App.Post(rest.RouteItem, s.rest.HandlerNewArticle())
	s.App.Get(rest.RouteItem+"/:id", s.rest.HandlerGetArticle())
	s.App.Put(rest.RouteItem+"/:id", s.rest.HandlerUpdateArticle())
	s.App.Delete(rest.RouteItem+"/:id", s.rest.HandlerDeleteArticle())

	s.App.Get(rest.RouteItems+"/all", s.rest.HandlerGetArticles())
	s.App.Get(rest.RouteItems, s.rest.HandlerGetArticlesWithPagination())

	s.App.Get(rest.RouteAlive, s.rest.HandlerAlive())
}
