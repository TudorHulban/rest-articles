package service

import (
	"context"
	"testing"

	"github.com/TudorHulban/rest-articles/infra/db"
	repository "github.com/TudorHulban/rest-articles/infra/repository/postgres"
	"github.com/stretchr/testify/require"
)

func TestCreateArticle(t *testing.T) {
	params := ParamsCreateArticle{
		Title: "xxx 1",
		URL:   "url 1",
	}

	ctx := context.Background()

	db, errCo := db.GetDBConnection()
	require.NoError(t, errCo)
	require.NotNil(t, db)

	repo, errRepo := repository.NewRepository(db)
	require.NoError(t, errRepo)

	serv := NewService(repo)

	rowID, errCr := serv.CreateArticle(ctx, &params)
	require.NoError(t, errCr)
	require.NotNil(t, rowID)
}
