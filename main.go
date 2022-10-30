package main

import (
	"fmt"
	"os"

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
	defer web.Stop()

	web.Start()
}
