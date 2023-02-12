package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/rsan92/teste-vibbra/app/application"
)

func main() {
	app, err := application.NewApplication()
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", app.Configuration.Application.ApiPort), app.HttpRouter))
}
