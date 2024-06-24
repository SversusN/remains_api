package app

import (
	"github.com/go-chi/chi/v5"
	"remains_api/config"
	"remains_api/internal/handlers"
	"remains_api/internal/repository"
	mssqlstorage "remains_api/internal/repository/mssql"
)

type App struct {
	c *config.Config
	h *handlers.Handlers
	s repository.Storage
}

func NewApp() *App {
	conf := config.GetConfig()
	storage := mssqlstorage.InitDatabase(conf)
	h := handlers.NewHandlers(storage)
	return &App{c: conf, s: storage, h: h}
}

func (a *App) CreateRouter() chi.Router {
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Get("/", a.h.GetAllHandler)
		r.Post("/", a.h.GetFilteredHandle)
		r.Post("/{group}", a.h.GetGroupHandler)
	})
	return r
}
