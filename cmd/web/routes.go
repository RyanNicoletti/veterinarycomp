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
	mux.Route("/", func(mux chi.Router) {
		mux.With(Paginate).Get("/", handlers.Repo.Home)
		mux.With(Paginate).Get("/{page}", handlers.Repo.Home)
	})
	mux.Get("/search-comp", handlers.Repo.SearchComp)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/add-comp", handlers.Repo.CompForm)
	mux.Post("/add-comp", handlers.Repo.PostCompForm)
	mux.Get("/login", handlers.Repo.Login)
	mux.Post("/login", handlers.Repo.PostLogin)
	mux.Get("/logout", handlers.Repo.Logout)
	mux.Route("/admin", func(mux chi.Router) {
		mux.Use(Auth)
		mux.Use(Admin)
		mux.Get("/dashboard", handlers.Repo.AdminDashboard)
		mux.Post("/verify-comp", handlers.Repo.VerifyComp)
		mux.Get("/download-verification", handlers.Repo.DownloadVerification)
	})
	staticFileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", staticFileServer))
	return mux
}
