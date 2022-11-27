package service

import (
	"context"
	"testing"

	"github.com/TudorHulban/rest-articles/app/apperrors"
	domain "github.com/TudorHulban/rest-articles/domain/article"
	"github.com/TudorHulban/rest-articles/infra/db"
	repository "github.com/TudorHulban/rest-articles/infra/repository/postgres"
	"github.com/stretchr/testify/require"
)

func TestServiceArticle(t *testing.T) {
	paramsCreate := ParamsCreateArticle{
		Title: "xxx 1",
		URL:   "url 1",
	}

	ctx := context.Background()

	db, errCo := db.GetDBConnection()
	require.NoError(t, errCo)
	require.NotNil(t, db)

	repo, errRepo := repository.NewRepository(db)
	require.NoError(t, errRepo)
	require.NoError(t, repo.Migration(&domain.Article{}))

	serv, errServ := NewService(repo)
	require.NoError(t, errServ)

	rowID, errCr := serv.CreateArticle(ctx, &paramsCreate)
	require.NoError(t, errCr)
	require.NotNil(t, rowID)

	reconstructedItemCreated, errGet := serv.GetArticle(ctx, rowID)
	require.NoError(t, errGet)
	require.Equal(t, paramsCreate.Title, reconstructedItemCreated.Title)

	newValue := "yyy 2"
	paramsUpdate := ParamsUpdateArticle{
		ID:    rowID,
		Title: &newValue,
	}

	errUpd := serv.UpdateArticle(ctx, &paramsUpdate)
	require.NoError(t, errUpd)

	reconstructedItemUpdated, errUpd := serv.GetArticle(ctx, rowID)
	require.NoError(t, errUpd)
	require.Equal(t, paramsUpdate.ID, reconstructedItemUpdated.ID)
	require.Equal(t, *paramsUpdate.Title, reconstructedItemUpdated.Title)
	require.Equal(t, paramsCreate.URL, reconstructedItemUpdated.URL)
	require.NotZero(t, reconstructedItemUpdated.UpdatedOn)

	items, errAll := serv.GetArticles(ctx)
	require.NoError(t, errAll)
	require.GreaterOrEqual(t, len(*items), 1)

	errDel := serv.DeleteArticle(ctx, rowID)
	require.NoError(t, errDel)

	reconstructedItemDeleted, errDel := serv.GetArticle(ctx, rowID)
	require.ErrorAs(t, errDel, &apperrors.ErrObjectNotFound{})
	require.Zero(t, reconstructedItemDeleted)
}
