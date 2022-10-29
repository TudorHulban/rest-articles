package rest

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/TudorHulban/rest-articles/app/service"
	"github.com/gofiber/fiber/v2"
)

type WebServer struct {
	app  *fiber.App
	serv *service.Service

	errShutdown error
	port        uint
}

const (
	_routeItem  = "/api/v1/article"
	_routeAlive = "/"
)

func NewWebServer(port uint, service *service.Service) (*WebServer, error) {
	if service == nil {
		return nil, errors.New("passed service is nil")
	}

	return &WebServer{
		app:  fiber.New(),
		port: port,
		serv: service,
	}, nil
}

func (s *WebServer) addRoutes() {
	s.app.Post(_routeItem, s.handleNewArticle())
	s.app.Get(_routeItem+"/:id", s.handleGetArticle())
	s.app.Put(_routeItem+"/:id", s.handleUpdateArticle())
	s.app.Delete(_routeItem+"/:id", s.handleDeleteArticle())

	s.app.Get(_routeItem, s.handleGetArticles())
	s.app.Get(_routeAlive, s.handleAlive())
}

func (s *WebServer) Start() {
	s.addRoutes()

	fmt.Println("web server started")

	s.errShutdown = s.app.Listen(":" + strconv.Itoa(int(s.port)))

	if s.errShutdown != nil {
		fmt.Printf("start - stopped now: %s\n", s.errShutdown)

		return
	}

	fmt.Println("start - stopped now: no error")
}

func (s *WebServer) Stop() {
	fmt.Println("stopping Fiber")

	if errShut := s.app.Shutdown(); errShut != nil {
		fmt.Printf("error Fiber: %s\n", errShut.Error())
	}
}
