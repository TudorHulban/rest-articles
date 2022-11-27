package web

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/TudorHulban/rest-articles/app/apperrors"
	"github.com/TudorHulban/rest-articles/app/service"
	"github.com/TudorHulban/rest-articles/infra/graphql/generated"
	"github.com/TudorHulban/rest-articles/infra/graphql/resolvers"
	"github.com/TudorHulban/rest-articles/infra/rest"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

type WebServer struct {
	App *fiber.App

	serv        *service.Service
	rest        *rest.Rest
	withGraphql bool

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

func NewWebServerWServiceAndGraphql(port uint, serv *service.Service) (*WebServer, error) {
	if serv == nil {
		return nil, errors.New("NewWebServerWServiceAndGraphql: passed service is nil")
	}

	crud, errREST := rest.NewRESTWService(serv)
	if errREST != nil {
		return nil, errREST
	}

	return &WebServer{
		App:         fiber.New(),
		port:        port,
		rest:        crud,
		serv:        serv,
		withGraphql: true,
	}, nil
}

func (s *WebServer) Errors(repoError error) *apperrors.ErrorApplication {
	return &apperrors.ErrorApplication{
		Area: apperrors.Areas[apperrors.ErrorAreaWebServer],
	}
}

func (s *WebServer) ErrorsWCode(code string, repoError error) *apperrors.ErrorApplication {
	return &apperrors.ErrorApplication{
		Area: apperrors.Areas[apperrors.ErrorAreaWebServer],
		Code: code,
	}
}

func (s *WebServer) Start() error {
	s.AddRESTRoutes()

	if s.withGraphql {
		if errGra := s.AddGraphql(); errGra != nil {
			return s.Errors(errGra)
		}
	}

	fmt.Println("web server started")

	s.errShutdown = s.App.Listen(":" + strconv.Itoa(int(s.port)))

	if s.errShutdown != nil {
		fmt.Printf("start - stopped now: %s\n", s.errShutdown)

		return s.Errors(s.errShutdown)
	}

	fmt.Println("start - stopped now: no error")

	return nil
}

// Stop releases web server and service.
// Returns closing errors from web and service.
func (s *WebServer) Stop() (error, error) {
	fmt.Println("stopping Fiber")

	var errWebClose, errServiceClose error

	if errShut := s.App.Shutdown(); errShut != nil {
		fmt.Printf("error Fiber: %s\n", errShut.Error()) // local handling can be improved

		errWebClose = errShut
	}

	return s.Errors(errWebClose), s.Errors(errServiceClose)
}

func (s *WebServer) AddGraphql() error {
	graphqlResolver, errGraphql := resolvers.NewResolverWService(s.serv)
	if errGraphql != nil {
		return s.ErrorsWCode(apperrors.ErrorGraphqlCODE, errGraphql)
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

	fmt.Println("Graphql support added")

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
