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

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/add-comp", handlers.Repo.AddComp)
	return mux
}
