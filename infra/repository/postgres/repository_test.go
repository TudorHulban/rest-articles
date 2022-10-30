package repository

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/TudorHulban/rest-articles/app/apperrors"
	domain "github.com/TudorHulban/rest-articles/domain/article"
	"github.com/TudorHulban/rest-articles/infra/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepository(t *testing.T) {
	dbConn, errCo := db.GetDBConnection()
	require.NoError(t, errCo)
	require.NotNil(t, dbConn)

	repo, errRepo := NewRepository(dbConn)
	require.NoError(t, errRepo)
	require.NoError(t, repo.Migration(&domain.Article{}))

	ctx := context.Background()

	item := domain.Article{
		Title:     "The Title " + strconv.Itoa(int(time.Now().Unix())),
		URL:       "The URL",
		CreatedOn: time.Now(),
	}

	insertID, errInsert := repo.CreateOne(&item)
	require.NoError(t, errInsert, "insert issues")

	reconstructedItem, errFind := repo.Find(ctx, insertID)
	require.NoError(t, errFind, errFind)
	require.Equal(t, insertID, reconstructedItem.ID)
	require.Nil(t, reconstructedItem.UpdatedOn)
	require.Nil(t, reconstructedItem.DeletedOn)

	articles, errAll := repo.FindAll(ctx)
	require.NoError(t, errAll)
	require.GreaterOrEqual(t, len(*articles), 1)

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
	require.ErrorAs(t, errFindDeleted, &apperrors.ErrObjectNotFound{})
	require.Zero(t, reconstructedItemDeleted)
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
	require.ErrorAs(t, errFind, &apperrors.ErrObjectNotFound{})
}

func TestPaginationScopes(t *testing.T) {
	dbConn, errCo := db.GetDBConnection()
	require.NoError(t, errCo)
	require.NotNil(t, dbConn)

	repo, errRepo := NewRepository(dbConn)
	require.NoError(t, errRepo)
	require.NoError(t, repo.Migration(&domain.Article{}))

	ctx := context.Background()

	item1 := domain.Article{
		Title:     "The Title 1 - " + strconv.Itoa(int(time.Now().Unix())),
		URL:       "http://paginationscopes.url1.eu",
		CreatedOn: time.Now(),
	}

	item2 := domain.Article{
		Title:     "The Title 2 - " + strconv.Itoa(int(time.Now().Unix())),
		URL:       "http://paginationscopes.url2.eu",
		CreatedOn: time.Now(),
	}

	item3 := domain.Article{
		Title:     "The Title 3 - " + strconv.Itoa(int(time.Now().Unix())),
		URL:       "http://paginationscopes.url3.eu",
		CreatedOn: time.Now(),
	}

	item4 := domain.Article{
		Title:     "The Title 4 - " + strconv.Itoa(int(time.Now().Unix())),
		URL:       "http://paginationscopes.url4.eu",
		CreatedOn: time.Now(),
	}

	howManyCreated, errCreateMany := repo.CreateMany(ctx, &domain.Articles{&item1, &item2, &item3, &item4})
	require.NoError(t, errCreateMany)
	require.Equal(t, 4, howManyCreated)

	testCases := []struct {
		description string
		input       Pagination
		want        int
	}{
		{"all in one page", Pagination{Limit: 4, Page: 1}, 4},
		{"one per page", Pagination{Limit: 1, Page: 1}, 1},
		{"outside of scope", Pagination{Limit: 100000, Page: 999999}, 0},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			items, _ := repo.FindAllPaginated(ctx, &tc.input)

			assert.Equal(t, tc.want, len(*items))
		})
	}
}
