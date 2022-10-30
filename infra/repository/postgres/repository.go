package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/TudorHulban/rest-articles/app/apperrors"
	domain "github.com/TudorHulban/rest-articles/domain/article"
	"gorm.io/gorm"
)

type Repository struct {
	DBConn *gorm.DB
}

func NewRepository(db *gorm.DB) (*Repository, error) {
	if db == nil {
		return nil, errors.New("passed DB connection is nil")
	}

	return &Repository{
		DBConn: db,
	}, nil
}

func (repo *Repository) Migration(model any) error {
	return repo.DBConn.AutoMigrate(model)
}

func (repo *Repository) Create(ctx context.Context, item *domain.Article) (int64, error) {
	errInsert := repo.DBConn.Create(item).Error

	return item.ID, errInsert
}

func (repo *Repository) Find(ctx context.Context, id int64) (*domain.Article, error) {
	var item domain.Article

	tx := repo.DBConn.
		Where("deleted_on is null").
		First(&item, id)

	if tx.Error != nil {
		if tx.Error.Error() == "record not found" {
			return nil, apperrors.ErrObjectNotFound{}
		}

		return nil, tx.Error
	}

	if tx.RowsAffected == 1 {
		return &item, nil
	}

	return nil, fmt.Errorf("duplicates found for ID: %d", id)
}

func (repo *Repository) FindAll(ctx context.Context) (*domain.Articles, error) {
	var items domain.Articles

	if errSelect := repo.DBConn.
		Where("deleted_on is null").
		Find(&items).Error; errSelect != nil {
		return nil, fmt.Errorf("FindAll: %w", errSelect)
	}

	return &items, nil
}

func (repo *Repository) FindAllPaginated(ctx context.Context, paginator *Pagination) (*Pagination, error) {
	var items domain.Articles

	if errSelect := repo.DBConn.
		Scopes(paginate(&domain.Article{}, paginator, repo.DBConn)).
		Where("deleted_on is null").
		Find(&items).Error; errSelect != nil {
		return nil, fmt.Errorf("FindAllPaginated: %w", errSelect)
	}

	return paginator, nil
}

func (repo *Repository) Update(ctx context.Context, item *domain.Article) error {
	rows := repo.DBConn.Model(item).Updates(item).RowsAffected

	if rows == 1 {
		return nil
	}

	return apperrors.ErrObjectNotFound{}
}
