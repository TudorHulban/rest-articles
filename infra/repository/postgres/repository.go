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

func (repo *Repository) Errors(repoError error) *apperrors.ErrorApplication {
	return &apperrors.ErrorApplication{
		Area: apperrors.Areas[apperrors.ErrorAreaRepository],
	}
}

func (repo *Repository) ErrorsWCode(code string, repoError error) *apperrors.ErrorApplication {
	return &apperrors.ErrorApplication{
		Area: apperrors.Areas[apperrors.ErrorAreaRepository],
		Code: code,
	}
}

func (repo *Repository) Migration(model any) error {
	if !repo.DBConn.Migrator().HasTable(model) {
		return repo.DBConn.AutoMigrate(model)
	}

	return nil
}

func (repo *Repository) CreateOne(ctx context.Context, item *domain.Article) (int64, error) {
	callDatabase := func() (int64, error) {
		errInsert := repo.DBConn.Create(item).Error

		return item.ID, repo.Errors(errInsert)
	}

	callDone := make(chan bool)
	defer close(callDone)

	var res int64
	var errCall error

	go func() {
		res, errCall = callDatabase()
		callDone <- true
	}()

	select {
	case <-ctx.Done():
		return 0, repo.Errors(fmt.Errorf("Find:%w", errors.New(apperrors.ErrorMsgForContextExpiration)))

	case <-callDone:
		return res, repo.Errors(errCall)
	}
}

// CreateMany
// TODO: assert if to move to transaction
func (repo *Repository) CreateMany(ctx context.Context, items *domain.Articles) (int, error) {
	for ix, item := range *items {
		_, errInsert := repo.CreateOne(ctx, item)
		if errInsert != nil {
			return ix - 1, repo.Errors(fmt.Errorf("CreateMany:%w", errInsert))
		}
	}

	return len(*items), nil
}

func (repo *Repository) Find(ctx context.Context, id int64) (*domain.Article, error) {
	callDatabase := func() (*domain.Article, error) {
		var item domain.Article

		tx := repo.DBConn.
			Where("deleted_on is null").
			First(&item, id)

		if tx.Error != nil {
			if tx.Error.Error() == "record not found" {
				return nil, repo.Errors(apperrors.ErrObjectNotFound{})
			}

			return nil, tx.Error
		}

		if tx.RowsAffected == 1 {
			return &item, nil
		}

		return nil, fmt.Errorf("duplicates found for ID: %d", id)
	}

	callDone := make(chan bool)
	defer close(callDone)

	var res *domain.Article
	var errCall error

	go func() {
		res, errCall = callDatabase()
		callDone <- true
	}()

	select {
	case <-ctx.Done():
		return nil, repo.Errors(fmt.Errorf("Find:%w", errors.New(apperrors.ErrorMsgForContextExpiration)))

	case <-callDone:
		return res, repo.Errors(errCall)
	}
}

func (repo *Repository) FindAll(ctx context.Context) (*domain.Articles, error) {
	callDatabase := func() (*domain.Articles, error) {
		var items domain.Articles

		if errSelect := repo.DBConn.
			Where("deleted_on is null").
			Order("id asc").
			Find(&items).
			Error; errSelect != nil {
			return nil, fmt.Errorf("FindAll: %w", errSelect)
		}

		return &items, nil
	}

	callDone := make(chan bool)
	defer close(callDone)

	var res *domain.Articles
	var errCall error

	go func() {
		res, errCall = callDatabase()
		callDone <- true
	}()

	select {
	case <-ctx.Done():
		return nil, repo.Errors(fmt.Errorf("Find:%w", errors.New(apperrors.ErrorMsgForContextExpiration)))

	case <-callDone:
		return res, repo.Errors(errCall)
	}
}

func (repo *Repository) FindAllPaginated(ctx context.Context, paginator *Pagination) (*domain.Articles, error) {
	callDatabase := func() (*domain.Articles, error) {
		var items domain.Articles

		if errSelect := repo.DBConn.
			Scopes(paginate(&domain.Article{}, paginator, repo.DBConn)).
			Where("deleted_on is null").
			Order("id asc").
			Find(&items).
			Error; errSelect != nil {
			return nil, fmt.Errorf("FindAllPaginated: %w", errSelect)
		}

		return &items, nil
	}

	callDone := make(chan bool)
	defer close(callDone)

	var res *domain.Articles
	var errCall error

	go func() {
		res, errCall = callDatabase()
		callDone <- true
	}()

	select {
	case <-ctx.Done():
		return nil, repo.Errors(fmt.Errorf("Find:%w", errors.New(apperrors.ErrorMsgForContextExpiration)))

	case <-callDone:
		return res, repo.Errors(errCall)
	}
}

func (repo *Repository) Update(ctx context.Context, item *domain.Article) error {
	callDatabase := func() error {
		rows := repo.DBConn.Model(item).Updates(item).RowsAffected

		if rows == 1 {
			return nil
		}

		return apperrors.ErrObjectNotFound{}
	}

	callDone := make(chan bool)
	defer close(callDone)

	var errCall error

	go func() {
		errCall = callDatabase()
		callDone <- true
	}()

	select {
	case <-ctx.Done():
		return repo.Errors(fmt.Errorf("Find:%w", errors.New(apperrors.ErrorMsgForContextExpiration)))

	case <-callDone:
		return repo.Errors(errCall)
	}
}
