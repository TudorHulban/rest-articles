package rest

import (
	"errors"

	"github.com/TudorHulban/rest-articles/app/service"
	repository "github.com/TudorHulban/rest-articles/infra/repository/postgres"
)

type Rest struct {
	serv *service.Service
}

func NewREST(repo *repository.Repository) (*Rest, error) {
	service, errServ := service.NewService(repo)
	if errServ != nil {
		return nil, errServ
	}

	return &Rest{
		serv: service,
	}, nil
}

func NewRESTWService(service *service.Service) (*Rest, error) {
	if service == nil {
		return nil, errors.New("passed service is nil")
	}

	return &Rest{
		serv: service,
	}, nil
}
