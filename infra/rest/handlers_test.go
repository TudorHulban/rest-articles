package rest

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/TudorHulban/rest-articles/app/service"
	domain "github.com/TudorHulban/rest-articles/domain/article"
	"github.com/TudorHulban/rest-articles/infra/db"
	repository "github.com/TudorHulban/rest-articles/infra/repository/postgres"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

var (
	_itemCreate = `{"title":"%s","url":"%s"}`
	_itemUpdate = `{"url":"%s"}`
)

func TestHandlers(t *testing.T) {
	require := require.New(t)

	dbConn, errCo := db.GetDBConnection()
	require.NoError(errCo)

	repo, errNew := repository.NewRepository(dbConn)
	require.NoError(errNew)

	require.NoError(repo.Migration(&domain.Article{}))

	service, errServ := service.NewService(repo)
	require.NoError(errServ)

	web, errWeb := NewWebServer(3000, service)
	require.NoError(errWeb)

	defer web.Stop()

	web.addRoutes()

	title := "The Title"
	urlCreate := "http://initial.abc.eu"

	resp, errBadReq := web.app.Test(httptest.NewRequest(http.MethodGet, _routeItem+"/abc", nil))
	utils.AssertEqual(t, nil, errBadReq)
	utils.AssertEqual(t, http.StatusBadRequest, resp.StatusCode)

	reqCreate := httptest.NewRequest(http.MethodPost, _routeItem, strings.NewReader(fmt.Sprintf(_itemCreate, title, urlCreate)))
	reqCreate.Header.Set("Content-type", "application/json")

	respPOST, errPOST := web.app.Test(reqCreate)
	utils.AssertEqual(t, nil, errPOST)
	utils.AssertEqual(t, http.StatusOK, respPOST.StatusCode)

	defer respPOST.Body.Close()

	bodyID, errID := io.ReadAll(respPOST.Body)
	t.Log(string(bodyID))

	require.NoError(errID)

	insertID := gjson.Get(string(bodyID), "id").String()
	t.Logf("insertID:%s", insertID)

	respItem, errGET := web.app.Test(httptest.NewRequest(http.MethodGet, _routeItem+"/"+insertID, nil))
	t.Logf("GET:\n%v", respItem)

	utils.AssertEqual(t, nil, errGET)
	utils.AssertEqual(t, http.StatusOK, respItem.StatusCode)

	defer respItem.Body.Close()

	bodyItem, errItem := io.ReadAll(respItem.Body)
	t.Logf("ITEM:\n%s", bodyItem)

	require.NoError(errItem)
	require.Equal(title, gjson.Get(string(bodyItem), "article.title").String())
	require.Equal(urlCreate, gjson.Get(string(bodyItem), "article.url").String())

	urlUpdate := "http://updated.abc.eu"

	reqUpdate := httptest.NewRequest(http.MethodPut, _routeItem+"/"+insertID, strings.NewReader(fmt.Sprintf(_itemUpdate, urlUpdate)))

	reqCreate.Header.Set("Content-type", "application/json")

	respPUT, errPUT := web.app.Test(reqUpdate)
	t.Logf("PUT:\n%v", *respPUT)

	defer respPUT.Body.Close()

	bodyPUT, errPUTBody := io.ReadAll(respPUT.Body)
	t.Logf("Body PUT response:\n%s", bodyPUT)
	t.Logf("Error PUT response:\n%s", errPUTBody)

	require.NoError(errPUTBody)

	utils.AssertEqual(t, nil, errPUT)
	utils.AssertEqual(t, http.StatusOK, respPUT.StatusCode)

	respUpdated, errUpdated := web.app.Test(httptest.NewRequest(http.MethodGet, _routeItem+"/"+insertID, nil))
	t.Logf("GET:\n%v", respUpdated)

	utils.AssertEqual(t, nil, errUpdated)
	utils.AssertEqual(t, http.StatusOK, respUpdated.StatusCode)

	defer respUpdated.Body.Close()

	bodyUpdated, errBodyUpdated := io.ReadAll(respUpdated.Body)
	t.Logf("ITEM Updated:\n%s", bodyUpdated)

	require.NoError(errBodyUpdated)
	require.Equal(title, gjson.Get(string(bodyUpdated), "article.title").String())
	require.Equal(urlUpdate, gjson.Get(string(bodyUpdated), "article.url").String())

	respDel, errDel := web.app.Test(httptest.NewRequest(http.MethodDelete, _routeItem+"/"+insertID, nil))
	utils.AssertEqual(t, nil, errDel)
	utils.AssertEqual(t, 200, respDel.StatusCode)

	respItemDeleted, errGETDeleted := web.app.Test(httptest.NewRequest(http.MethodGet, _routeItem+"/"+insertID, nil))
	t.Logf("GET:\n%v", respItemDeleted)
	t.Logf("Error DELETE response:\n%s", errGETDeleted)

	utils.AssertEqual(t, nil, errGETDeleted)
	utils.AssertEqual(t, http.StatusOK, respItemDeleted.StatusCode)

	defer respItemDeleted.Body.Close()

	bodyDeleted, errBodyDeleted := io.ReadAll(respItemDeleted.Body)
	t.Logf("ITEM Deleted:\n%s", string(bodyDeleted))

	require.NoError(errBodyDeleted, errBodyDeleted)
}
