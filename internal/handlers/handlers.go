package handlers

import (
	"net/http"

	"github.com/ryannicoletti/veterinarycomp/internal/config"
	"github.com/ryannicoletti/veterinarycomp/internal/models"
	"github.com/ryannicoletti/veterinarycomp/internal/render"
)

var Repo *Repository

type Repository struct {
	App *config.AppConfig
}

func CreateNewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

func NewHandlers(r *Repository) {
	Repo = r
}

func (repo *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIp := r.RemoteAddr
	repo.App.Session.Put(r.Context(), "remote_ip", remoteIp)
	render.RenderTemplate(w, "home.page.tmpl", &models.TemplateData{})
}

func (repo *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "about.page.tmpl", &models.TemplateData{})
}

func (repo *Repository) AddComp(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "add-comp.page.tmpl", &models.TemplateData{})
}
