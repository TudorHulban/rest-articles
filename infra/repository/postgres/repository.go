package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/TudorHulban/rest-articles/domain/article"
	"github.com/TudorHulban/rest-articles/infra/db"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	Db *sqlx.DB
}

func NewRepository(db *sqlx.DB) (*Repository, error) {
	if db == nil {
		return nil, errors.New("passed DB connection is nil")
	}

	return &Repository{
		Db: db,
	}, nil
}

func (repo *Repository) Create(ctx context.Context, item *article.Article) error {
	query := `INSERT INTO articles (title, url, created_on)
	VALUES (:title, :url, :created_on) RETURNING id;`

	rows, errDB := repo.Db.NamedQueryContext(ctx, query, item)
	if errDB != nil {
		return db.HandleError(errDB)
	}

	for rows.Next() {
		errScan := rows.StructScan(item)
		if errScan != nil {
			return db.HandleError(errScan)
		}
	}

	return nil
}

func (repo *Repository) Find(ctx context.Context, id int) (*article.Article, error) {
	var item article.Article

	query := fmt.Sprintf(
		"SELECT * FROM articles WHERE id = $1 AND deleted_on IS NULL",
	)

	errDB := repo.Db.GetContext(ctx, &item, query, id)

	return &item, db.HandleError(errDB)
}

func (repo *Repository) FindAll(ctx context.Context) (article.Articles, error) {
	var items article.Articles

	query := fmt.Sprintf(
		"SELECT * FROM articles WHERE deleted_on IS NULL",
	)

	errDB := repo.Db.SelectContext(ctx, &items, query)

	return items, db.HandleError(errDB)
}

func (repo *Repository) Update(ctx context.Context, item *article.Article) error {
	query := `UPDATE articles
                SET title = :title, 
                    url = :url, 
                    created_on = :created_on, 
                    updated_on = :updated_on, 
                    deleted_on = :deleted_on
                WHERE id = :id;`

	_, errDB := repo.Db.NamedExecContext(ctx, query, item)

	return db.HandleError(errDB)
}
