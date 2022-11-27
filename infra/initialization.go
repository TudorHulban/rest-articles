package infra

import (
	"fmt"

	"github.com/TudorHulban/rest-articles/app/apperrors"
	"github.com/TudorHulban/rest-articles/app/service"
	domain "github.com/TudorHulban/rest-articles/domain/article"
	"github.com/TudorHulban/rest-articles/infra/db"
	repository "github.com/TudorHulban/rest-articles/infra/repository/postgres"
	"github.com/TudorHulban/rest-articles/infra/web"
)

func Initialize() (*web.WebServer, *apperrors.ErrorApplication) {
	errorApp := apperrors.ErrorApplication{
		Area: apperrors.Areas[apperrors.ErrorAreaInfra],
	}

	dbConn, errCo := db.GetDBConnection()
	if errCo != nil {
		errorApp.AreaError = fmt.Errorf(apperrors.ErrorMsgConnectionCreation, errCo)
		errorApp.OSExit = &apperrors.OSExitForDatabaseIssues

		return nil, &errorApp
	}

	repo, errRepo := repository.NewRepository(dbConn)
	if errRepo != nil {
		errorApp.AreaError = fmt.Errorf(apperrors.ErrorMsgRepositoryCreation, errRepo)
		errorApp.OSExit = &apperrors.OSExitForRepositoryIssues

		return nil, &errorApp
	}

	errMi := repo.Migration(&domain.Article{})
	if errMi != nil {
		errorApp.AreaError = fmt.Errorf(apperrors.ErrorMsgRepositoryMigrationsRun, errMi)
		errorApp.OSExit = &apperrors.OSExitForRepositoryMigrationsIssues

		return nil, &errorApp
	}

	service, errServ := service.NewService(repo)
	if errServ != nil {
		errorApp.AreaError = fmt.Errorf(apperrors.ErrorMsgServiceCreation, errServ)
		errorApp.OSExit = &apperrors.OSExitForServiceIssues

		return nil, &errorApp
	}

	web, errWeb := web.NewWebServerWServiceAndGraphql(3000, service)
	if errWeb != nil {
		errorApp.AreaError = fmt.Errorf(apperrors.ErrorMsgWebServerCreation, errServ)
		errorApp.OSExit = &apperrors.OSExitForWebServerIssues

		return nil, &errorApp
	}

	return web, nil
}
