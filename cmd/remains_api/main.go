package main

import (
	"log"
	"net/http"
	"remains_api/internal/app"
)

// @title chi-swagger example APIs
// @version 1.0
// @description chi-swagger example APIs
// @BasePath /
func main() {
	a := app.NewApp()
	r := a.CreateRouter()

	err := http.ListenAndServe(a.C.Port, r)
	if err != nil {
		log.Fatal(err)
	}
}
