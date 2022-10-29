package rest

import (
	"fmt"
	"io/ioutil"
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
	_itemUpdate = `{"id":"%s","url":"%s"}`
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
	urlCreate := "http://abc.eu"

	resp, errBadReq := web.app.Test(httptest.NewRequest(http.MethodGet, _routeItem+"/abc", nil))
	utils.AssertEqual(t, nil, errBadReq)
	utils.AssertEqual(t, http.StatusBadRequest, resp.StatusCode)

	reqCreate := httptest.NewRequest(http.MethodPost, _routeItem, strings.NewReader(fmt.Sprintf(_itemCreate, title, urlCreate)))
	reqCreate.Header.Set("Content-type", "application/json")

	respPOST, errPOST := web.app.Test(reqCreate)
	utils.AssertEqual(t, nil, errPOST)
	utils.AssertEqual(t, http.StatusOK, respPOST.StatusCode)

	defer respPOST.Body.Close()

	bodyID, errID := ioutil.ReadAll(respPOST.Body)
	t.Log(string(bodyID))

	require.NoError(errID)

	insertID := gjson.Get(string(bodyID), "id").String()
	t.Logf("insertID:%s", insertID)

	respItem, errGET := web.app.Test(httptest.NewRequest(http.MethodGet, _routeItem+"/"+insertID, nil))
	t.Logf("GET:\n%v", respItem)

	utils.AssertEqual(t, nil, errGET)
	utils.AssertEqual(t, http.StatusOK, respItem.StatusCode)

	defer respItem.Body.Close()

	bodyItem, errItem := ioutil.ReadAll(respItem.Body)
	t.Logf("ITEM:\n%s", bodyItem)

	require.NoError(errItem)
	require.Equal(title, gjson.Get(string(bodyItem), "article.title").String())
	require.Equal(urlCreate, gjson.Get(string(bodyItem), "article.url").String())

	urlUpdate := "http://xyz.eu"

	reqUpdate := httptest.NewRequest(http.MethodPut, _routeItem, strings.NewReader(fmt.Sprintf(_itemUpdate, insertID, urlUpdate)))
	reqCreate.Header.Set("Content-type", "application/json")

	respPUT, errPUT := web.app.Test(reqUpdate)
	t.Logf("PUT:\n%v", *respPUT)

	defer respPUT.Body.Close()

	bodyPUT, errPUTBody := ioutil.ReadAll(respPUT.Body)
	t.Logf("Body PUT response:\n%s", bodyPUT)
	t.Logf("Error PUT response:\n%s", errPUTBody)

	require.NoError(errPUTBody)

	utils.AssertEqual(t, nil, errPUT)
	utils.AssertEqual(t, http.StatusOK, respPUT.StatusCode)
}