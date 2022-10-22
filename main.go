package main

import "github.com/TudorHulban/rest-articles/infra/rest"

func main() {
	web := rest.NewWebServer(3000)
	web.Start()
}
