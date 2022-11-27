package web

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/TudorHulban/rest-articles/app/service"
	"github.com/TudorHulban/rest-articles/infra/graphql/generated"
	"github.com/TudorHulban/rest-articles/infra/graphql/resolvers"
	"github.com/TudorHulban/rest-articles/infra/rest"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

type WebServer struct {
	App *fiber.App

	serv *service.Service
	rest *rest.Rest

	errShutdown error
	port        uint
}

func NewWebServerREST(port uint, rest *rest.Rest) (*WebServer, error) {
	return &WebServer{
		App:  fiber.New(),
		port: port,
		rest: rest,
	}, nil
}

func NewWebServerWService(port uint, serv *service.Service) (*WebServer, error) {
	crud, errREST := rest.NewRESTWService(serv)
	if errREST != nil {
		return nil, errREST
	}

	return &WebServer{
		App:  fiber.New(),
		port: port,
		rest: crud,
	}, nil
}

func (s *WebServer) Start() {
	s.AddRESTRoutes()

	fmt.Println("web server started")

	s.errShutdown = s.App.Listen(":" + strconv.Itoa(int(s.port)))

	if s.errShutdown != nil {
		fmt.Printf("start - stopped now: %s\n", s.errShutdown)

		return
	}

	fmt.Println("start - stopped now: no error")
}

func (s *WebServer) StartWGraphql() {
	s.AddGraphql()

	s.Start()
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

func (s *WebServer) AddGraphql() error {
	graphqlResolver, errGraphql := resolvers.NewResolverWService(s.serv)
	if errGraphql != nil {
		return errGraphql
	}

	server := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
		Resolvers: graphqlResolver,
	}))

	graphqlHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		server.ServeHTTP(w, r)
	})

	s.App.All(_RouteGraphql, func(c *fiber.Ctx) error {
		fasthttpadaptor.NewFastHTTPHandler(graphqlHandler)(c.Context())

		return nil
	})

	return nil
}

func (s *WebServer) AddRESTRoutes() {
	s.App.Post(rest.RouteItem, s.rest.HandlerNewArticle())
	s.App.Get(rest.RouteItem+"/:id", s.rest.HandlerGetArticle())
	s.App.Put(rest.RouteItem+"/:id", s.rest.HandlerUpdateArticle())
	s.App.Delete(rest.RouteItem+"/:id", s.rest.HandlerDeleteArticle())

	s.App.Get(rest.RouteItems+"/all", s.rest.HandlerGetArticles())
	s.App.Get(rest.RouteItems, s.rest.HandlerGetArticlesWithPagination())

	s.App.Get(rest.RouteAlive, s.rest.HandlerAlive())
}
