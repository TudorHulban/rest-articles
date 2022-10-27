package repository

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	domain "github.com/TudorHulban/rest-articles/domain/article"
	"github.com/TudorHulban/rest-articles/infra/db"
	"github.com/stretchr/testify/require"
)

func TestRepository(t *testing.T) {
	db, errCo := db.GetDBConnection()
	require.NoError(t, errCo)
	require.NotNil(t, db)

	repo, errRepo := NewRepository(db)
	require.NoError(t, errRepo)

	item := domain.Article{
		Title:     "The Title " + strconv.Itoa(int(time.Now().Unix())),
		URL:       "The URL",
		CreatedOn: time.Now(),
	}

	ctx := context.Background()

	insertID, errInsert := repo.Create(ctx, &item)
	require.NoError(t, errInsert, "insert issues")

	reconstructedArticle, errFind := repo.Find(ctx, insertID)
	require.NoError(t, errFind, errFind)
	require.Equal(t, insertID, reconstructedArticle.ID)

	articles, errAll := repo.FindAll(ctx)
	require.NoError(t, errAll)
	require.GreaterOrEqual(t, len(*articles), 2)

	fmt.Println(articles)

	// TODO; test update
	// TODO: test find with not found
}
