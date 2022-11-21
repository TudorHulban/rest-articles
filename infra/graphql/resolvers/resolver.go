package resolvers

import (
	"context"
	"errors"

	"github.com/TudorHulban/rest-articles/app/service"
	"github.com/TudorHulban/rest-articles/infra/graphql/generated/models"
	repository "github.com/TudorHulban/rest-articles/infra/repository/postgres"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	serv *service.Service
}

func NewResolver(repo *repository.Repository) (*Resolver, error) {
	service, errServ := service.NewService(repo)
	if errServ != nil {
		return nil, errServ
	}

	return &Resolver{
		serv: service,
	}, nil
}

func NewResolverWService(service *service.Service) (*Resolver, error) {
	if service == nil {
		return nil, errors.New("passed service is nil")
	}

	return &Resolver{
		serv: service,
	}, nil
}

func (r *mutationResolver) CreateArticle(ctx context.Context, input models.Article) (*models.Article, error) {
	paramsArticle := service.ParamsCreateArticle{
		Title: input.Title,
		URL:   input.URL,
	}

	createdArticleID, errCre := r.serv.CreateArticle(ctx, &paramsArticle)
	if errCre != nil {
		return nil, errCre // TODO: add proper area error
	}

	return &models.Article{
		ID:    int(createdArticleID),
		Title: input.Title,
		URL:   input.URL,
	}, nil
}
