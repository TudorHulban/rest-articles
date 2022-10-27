package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	domain "github.com/TudorHulban/rest-articles/domain/article"
	"github.com/TudorHulban/rest-articles/infra/db"
)

type Repository struct {
	Db *sql.DB
}

func NewRepository(db *sql.DB) (*Repository, error) {
	if db == nil {
		return nil, errors.New("passed DB connection is nil")
	}

	return &Repository{
		Db: db,
	}, nil
}

func (repo *Repository) Create(ctx context.Context, item *domain.Article) (int64, error) {
	stmt, errPrep := repo.Db.PrepareContext(ctx, `INSERT INTO articles(title, url, created_on) values(?,?,?)`)
	if errPrep != nil {
		return 0, db.HandleError(errPrep)
	}
	defer stmt.Close()

	result, errDB := stmt.ExecContext(ctx, item.Title, item.URL, item.CreatedOn)
	if errDB != nil {
		return 0, db.HandleError(errDB)
	}

	return result.LastInsertId()
}

func (repo *Repository) Find(ctx context.Context, id int64) (*domain.Article, error) {
	var item domain.Article

	query := fmt.Sprintf("SELECT * FROM articles WHERE id = $1 AND deleted_on IS NULL")

	_ = repo.Db.QueryRowContext(ctx, query, id).Scan(&item)

	if item == nil {
		return nil, db.ErrObjectNotFound{}
	}

	return &item, nil
}

func (repo *Repository) FindAll(ctx context.Context) (*domain.Articles, error) {
	var items domain.Articles

	query := fmt.Sprintf("SELECT * FROM articles WHERE deleted_on IS NULL")

	rows, errDB := repo.Db.QueryContext(ctx, query)
	if errDB != nil {
		return nil, db.HandleError(errDB)
	}
	defer rows.Close()

	for rows.Next() {
		var item domain.Article

		errScan := rows.Scan(&item)
		if errScan != nil {
			return nil, db.HandleError(errScan)
		}

		items = append(items, &item)
	}

	return &items, nil
}

func (repo *Repository) Update(ctx context.Context, item *domain.Article) error {
	query := `UPDATE articles
                SET title = :title, 
                    url = :url, 
                    created_on = :created_on, 
                    updated_on = :updated_on, 
                    deleted_on = :deleted_on
                WHERE id = :id;`

	_, errDB := repo.Db.ExecContext(ctx, query, item)

	return db.HandleError(errDB)
}
