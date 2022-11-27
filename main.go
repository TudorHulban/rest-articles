package main

import (
	"fmt"
	"os"

	"github.com/TudorHulban/rest-articles/app/apperrors"
	"github.com/TudorHulban/rest-articles/infra"
)

func main() {
	web, errWeb := infra.Initialize()
	if errWeb != nil {
		fmt.Println(errWeb)

		if errWeb.OSExit != nil {
			os.Exit(*errWeb.OSExit)
		}
	}
	defer func() {
		errWebStop, errServStop := web.Stop()
		if errWebStop != nil {
			fmt.Println(errWebStop)
		}

		if errServStop != nil {
			fmt.Println(errServStop)
		}
	}()

	if errStart := web.Start(); errStart != nil {
		fmt.Println(errStart)

		os.Exit(apperrors.OSExitForGraphqlIssues)
	}
}
