package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/ryannicoletti/veterinarycomp/internal/config"
	"github.com/ryannicoletti/veterinarycomp/internal/handlers"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()
	// middlewares
	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)
	mux.Use(Paginate)
	mux.Route("/", func(mux chi.Router) {
		mux.With(Paginate).Get("/", handlers.Repo.Home)
		mux.With(Paginate).Get("/{page}", handlers.Repo.Home)
	})
	mux.Post("/search-comp", handlers.Repo.SearchComp)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/add-comp", handlers.Repo.CompForm)
	mux.Post("/add-comp", handlers.Repo.PostCompForm)
	staticFileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", staticFileServer))
	return mux
}
