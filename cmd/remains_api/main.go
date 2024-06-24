package main

import (
	"net/http"
	"remains_api/internal/app"
)

func main() {
	newapp := app.NewApp()
	r := newapp.CreateRouter()
	http.ListenAndServe(":8080", r)
}
