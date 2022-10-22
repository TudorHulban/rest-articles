package rest

import (
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
	_routeArticle = "/api/v1/article"
	_routeAlive   = "/"
)

func NewWebServer(port uint) *WebServer {
	return &WebServer{
		app:  fiber.New(),
		port: port,
	}
}

func (s *WebServer) addRoutes() {
	s.app.Post(_routeArticle, s.handleNewArticle())
	s.app.Get(_routeArticle+"/:id", s.handleGetArticle())
	s.app.Get(_routeArticle, s.handleGetArticles())
	s.app.Get(_routeAlive, s.handleAlive())
}

func (s *WebServer) Start() {
	s.addRoutes()

	fmt.Println("web server started")

	s.errShutdown = s.app.Listen(":" + strconv.Itoa(int(s.port)))
}

func (s *WebServer) Stop() {
	fmt.Println("stopping Fiber")

	if errShut := s.app.Shutdown(); errShut != nil {
		fmt.Printf("error Fiber: %s\n", errShut.Error())
	}
}
