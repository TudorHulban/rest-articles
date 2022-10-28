package repository

import (
	"context"
	"strconv"
	"testing"
	"time"

	domain "github.com/TudorHulban/rest-articles/domain/article"
	"github.com/TudorHulban/rest-articles/infra/db"
	"github.com/stretchr/testify/require"
)

func TestRepository(t *testing.T) {
	dbConn, errCo := db.GetDBConnection()
	require.NoError(t, errCo)
	require.NotNil(t, dbConn)

	repo, errRepo := NewRepository(dbConn)
	require.NoError(t, errRepo)

	item := domain.Article{
		Title:     "The Title " + strconv.Itoa(int(time.Now().Unix())),
		URL:       "The URL",
		CreatedOn: time.Now(),
	}

	ctx := context.Background()

	insertID, errInsert := repo.Create(ctx, &item)
	require.NoError(t, errInsert, "insert issues")

	reconstructedItem, errFind := repo.Find(ctx, insertID)
	require.NoError(t, errFind, errFind)
	require.Equal(t, insertID, reconstructedItem.ID)
	require.Nil(t, reconstructedItem.UpdatedOn)
	require.Nil(t, reconstructedItem.DeletedOn)

	articles, errAll := repo.FindAll(ctx)
	require.NoError(t, errAll)
	require.GreaterOrEqual(t, len(*articles), 1)
	require.Equal(t, insertID, (*articles)[len(*articles)-1].ID)

	updatedTimestamp := time.Now()
	itemUpdated := domain.Article{
		ID:        insertID,
		Title:     "The Title " + strconv.Itoa(int(time.Now().Unix())),
		URL:       "The URL",
		UpdatedOn: &updatedTimestamp,
	}

	errUpd := repo.Update(ctx, &itemUpdated)
	require.NoError(t, errUpd)

	reconstructedItemUpdated, errFindUpdated := repo.Find(ctx, insertID)
	require.NoError(t, errFindUpdated, errFindUpdated)
	require.Equal(t, insertID, reconstructedItemUpdated.ID)
	require.NotNil(t, reconstructedItemUpdated.UpdatedOn)
	require.Nil(t, reconstructedItem.DeletedOn)

	deletedTimestamp := time.Now()
	itemDeleted := domain.Article{
		ID:        insertID,
		Title:     "The Title " + strconv.Itoa(int(time.Now().Unix())),
		URL:       "The URL",
		DeletedOn: &deletedTimestamp,
	}

	errDel := repo.Update(ctx, &itemDeleted)
	require.NoError(t, errDel)

	reconstructedItemDeleted, errFindDeleted := repo.Find(ctx, insertID)
	require.NoError(t, errFindDeleted, errFindDeleted)
	require.Equal(t, insertID, reconstructedItemDeleted.ID)
	require.NotNil(t, reconstructedItemDeleted.UpdatedOn)
	require.NotNil(t, reconstructedItemDeleted.DeletedOn)
}

func TestItemNotFound(t *testing.T) {
	dbConn, errCo := db.GetDBConnection()
	require.NoError(t, errCo)
	require.NotNil(t, dbConn)

	repo, errRepo := NewRepository(dbConn)
	require.NoError(t, errRepo)

	ctx := context.Background()

	_, errFind := repo.Find(ctx, -1)
	require.Error(t, errFind)
	require.ErrorAs(t, errFind, &db.ErrObjectNotFound{})
}
