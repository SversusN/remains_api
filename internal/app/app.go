package app

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
	"net/http"
	"os"
	"path/filepath"
	"remains_api/config"
	"remains_api/internal/handlers"
	"remains_api/internal/repository"
	mssqlstorage "remains_api/internal/repository/mssql"
	"remains_api/mw"
	"remains_api/pkg/auth"
	"strings"
)

type App struct {
	C  *config.Config
	h  *handlers.Handlers
	s  repository.Storage
	au *jwtauth.JWTAuth
}

func NewApp() *App {
	conf := config.GetConfig()
	au := auth.NewAuth()
	storage := mssqlstorage.InitDatabase(conf)
	h := handlers.NewHandlers(storage, au)

	return &App{C: conf, s: storage, h: h, au: au}
}

func (a *App) CreateRouter() chi.Router {
	r := chi.NewRouter()
	r.Use(mw.Cors(r))
	folder, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(folder, "static_files"))
	FileServer(r, "/", filesDir)
	r.Group(func(r chi.Router) {

		r.Post("/api/login/*", a.h.LoginUser)
		r.Options("/api/login/*", a.h.LoginUser)
		r.Post("/logout", func(w http.ResponseWriter, r *http.Request) {})
	})
	r.Group(func(r chi.Router) {

		r.Use(jwtauth.Verifier(a.au))
		r.Options("/getall/*", a.h.GetAllHandler)
		r.Get("/getall/*", a.h.GetAllHandler)

	})

	return r
}

func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
