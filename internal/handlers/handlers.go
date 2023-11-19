package handlers

import (
	"net/http"

	"github.com/ryannicoletti/veterinarycomp/internal/config"
	"github.com/ryannicoletti/veterinarycomp/internal/render"
	"github.com/ryannicoletti/veterinarycomp/models"
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
	render.RenderTemplate(w, "home.page.tmpl", &models.TemplateData{})
}

func (repo *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "about.page.tmpl", &models.TemplateData{})
}

func (repo *Repository) AddComp(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "add-comp.page.tmpl", &models.TemplateData{})
}
