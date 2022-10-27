package main

import (
	"fmt"
	"os"

	"github.com/TudorHulban/rest-articles/app/service"
	domain "github.com/TudorHulban/rest-articles/domain/article"
	"github.com/TudorHulban/rest-articles/infra/db"
	repository "github.com/TudorHulban/rest-articles/infra/repository/postgres"
	"github.com/TudorHulban/rest-articles/infra/rest"
)

func main() {
	dbConn, errCo := db.GetDBConnection()
	if errCo != nil {
		fmt.Printf("DB connection creation: %s", errCo)
		os.Exit(1)
	}

	repo, errRepo := repository.NewRepository(dbConn)
	if errRepo != nil {
		fmt.Printf("repository creation: %s", errRepo)
		os.Exit(2)
	}

	errMi := repo.Migration(&domain.Article{})
	if errMi != nil {
		fmt.Printf("repository migrations: %s", errMi)
		os.Exit(3)
	}

	service := service.NewService(repo)

	web := rest.NewWebServer(3000, service)
	web.Start()
}
