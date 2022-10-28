package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	domain "github.com/TudorHulban/rest-articles/domain/article"
	"github.com/TudorHulban/rest-articles/infra/db"
	repository "github.com/TudorHulban/rest-articles/infra/repository/postgres"
	"github.com/asaskevich/govalidator"
)

type Service struct {
	repo repository.Repository
}

func NewService(repo *repository.Repository) (*Service, error) {
	if repo == nil {
		return nil, errors.New("passed repository is nil")
	}

	return &Service{
		repo: *repo,
	}, nil
}

type ParamsCreateArticle struct {
	Title string `valid:"required"`
	URL   string `valid:"required"`
}

func (s *Service) CreateArticle(ctx context.Context, params *ParamsCreateArticle) (int64, error) {
	if _, errVa := govalidator.ValidateStruct(params); errVa != nil {
		return 0, fmt.Errorf(errorArgumentMessage, errVa)
	}

	item := domain.Article{
		Title:     params.Title,
		URL:       params.URL,
		CreatedOn: time.Now(),
	}

	itemID, errCr := s.repo.Create(ctx, &item)
	if errCr != nil {
		return 0, fmt.Errorf("CreateArticle: %w", errCr)
	}

	return itemID, nil
}

func (s *Service) GetArticle(ctx context.Context, id int64) (*domain.Article, error) {
	article, errFind := s.repo.Find(ctx, id)
	switch {
	case errFind == nil:
		return article, nil

	case errors.As(errFind, &db.ErrObjectNotFound{}):
		return nil, db.ErrObjectNotFound{}

	default:
		return nil, errFind
	}
}

func (s *Service) GetArticles(ctx context.Context) (*domain.Articles, error) {
	articles, errAll := s.repo.FindAll(ctx)
	switch {
	case errAll == nil:
		return articles, nil

	case errors.As(errAll, &db.ErrObjectNotFound{}):
		return nil, db.ErrObjectNotFound{}

	default:
		return nil, errAll
	}
}

type ParamsUpdateArticle struct {
	ID    int64 `valid:"required"`
	Title *string
	URL   *string
}

func (s *Service) UpdateArticle(ctx context.Context, params *ParamsUpdateArticle) error {
	if _, errVa := govalidator.ValidateStruct(params); errVa != nil {
		return fmt.Errorf(errorArgumentMessage, errVa)
	}

	article, errDB := s.repo.Find(ctx, params.ID)
	if errDB != nil {
		return errDB
	}

	if params.Title != nil {
		article.Title = *params.Title
	}

	if params.URL != nil {
		article.URL = *params.URL
	}

	return s.repo.Update(ctx, article)
}

func (s *Service) DeleteArticle(ctx context.Context, id int64) error {
	article, errDB := s.repo.Find(ctx, id)
	if errDB != nil {
		return errDB
	}

	now := time.Now()
	article.DeletedOn = &now

	return s.repo.Update(ctx, article)
}
