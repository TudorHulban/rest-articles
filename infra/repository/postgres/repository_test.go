package repository

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/TudorHulban/rest-articles/domain/article"
	"github.com/TudorHulban/rest-articles/infra/db"
	"github.com/stretchr/testify/require"
)

func TestRepository(t *testing.T) {
	db, errCo := db.Connect(db.NewTestConfigDB())
	require.NoError(t, errCo)
	require.NotNil(t, db)

	repo, errRepo := NewRepository(db)
	require.NoError(t, errRepo)

	noID := 8

	item := article.Article{
		ID:        noID,
		Title:     "The Title " + strconv.Itoa(noID),
		URL:       "The URL",
		CreatedOn: time.Now(),
	}

	ctx := context.Background()

	require.NoError(t, repo.Create(ctx, &item), "insert issues")

	reconstructedArticle, errFind := repo.Find(ctx, noID)
	require.NoError(t, errFind)
	require.Equal(t, noID, reconstructedArticle.ID)

	articles, errAll := repo.FindAll(ctx)
	require.NoError(t, errAll)
	require.GreaterOrEqual(t, len(articles), 2)

	fmt.Println(articles)

	// TODO; test update
	// TODO: test find with not found
}
