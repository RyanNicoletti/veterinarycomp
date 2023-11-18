package handlers

import (
	"net/http"

	"veterinarycomp/internal/config"
	"veterinarycomp/internal/render"
)

var Repo *Repository

type Repository struct {
	App *config.AppConfig
}

func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

func NewHandlers(r *Repository) {
	Repo = r
}

func (repo *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "home.page.tmpl")
}

func (repo *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "about.page.tmpl")
}

func (repo *Repository) AddComp(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "add-comp.page.tmpl")
}
