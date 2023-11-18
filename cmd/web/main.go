package main

import (
	"log"
	"net/http"
	"veterinarycomp/internal/config"
	"veterinarycomp/internal/handlers"
	"veterinarycomp/internal/render"
)

func main() {
	var app config.AppConfig
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache", err)
	}
	app.TemplateCache = tc
	app.UseCache = false
	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)
	render.NewTemplates(&app)
	http.HandleFunc("/", repo.Home)
	http.HandleFunc("/about", repo.About)
	http.ListenAndServe(":8080", nil)
}
