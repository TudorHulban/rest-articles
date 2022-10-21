package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/TudorHulban/rest-articles/domain/article"
	"github.com/TudorHulban/rest-articles/infra/db"
	repository "github.com/TudorHulban/rest-articles/infra/repository/postgres"
	"github.com/asaskevich/govalidator"
)

type Service struct {
	repo repository.Repository
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		repo: *repo,
	}
}

type ParamsCreateArticle struct {
	Title string `valid:"required"`
	URL   string `valid:"required"`
}

func (s *Service) CreateArticle(ctx context.Context, params *ParamsCreateArticle) (int, error) {
	if _, errVa := govalidator.ValidateStruct(params); errVa != nil {
		return 0, fmt.Errorf(errorArgumentMessage, errVa)
	}

	tx, errTx := s.repo.Db.BeginTxx(ctx, nil)
	if errTx != nil {
		return 0, errTx
	}
	defer tx.Rollback()

	item := article.Article{
		Title:     params.Title,
		URL:       params.URL,
		CreatedOn: time.Now(),
	}

	if errCr := s.repo.Create(ctx, &item); errCr != nil {
		return 0, errCr
	}

	if errCo := tx.Commit(); errCo != nil {
		return 0, errCo
	}

	return item.ID, nil
}

func (s *Service) GetArticle(ctx context.Context, id int) (*article.Article, error) {
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

func (s *Service) GetArticles(ctx context.Context) (*article.Articles, error) {
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
	ID    int `valid:"required"`
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

	tx, errTx := s.repo.Db.BeginTxx(ctx, nil)
	if errTx != nil {
		return errTx
	}
	defer tx.Rollback()

	if errDB = s.repo.Update(ctx, article); errDB != nil {
		return errDB
	}

	return tx.Commit()
}

func (s *Service) Delete(ctx context.Context, id int) error {
	article, errDB := s.repo.Find(ctx, id)
	if errDB != nil {
		return errDB
	}

	tx, errTx := s.repo.Db.BeginTxx(ctx, nil)
	if errTx != nil {
		return errTx
	}
	defer tx.Rollback()

	now := time.Now()
	article.DeletedOn = &now

	errUpd := s.repo.Update(ctx, article)
	if errUpd != nil {
		return errUpd
	}

	return tx.Commit()
}
