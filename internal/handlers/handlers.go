package handlers

import (
	"io"
	"log"
	"net/http"
	"time"

	"github.com/ryannicoletti/veterinarycomp/internal/config"
	"github.com/ryannicoletti/veterinarycomp/internal/driver"
	"github.com/ryannicoletti/veterinarycomp/internal/forms"
	"github.com/ryannicoletti/veterinarycomp/internal/models"
	"github.com/ryannicoletti/veterinarycomp/internal/render"
	"github.com/ryannicoletti/veterinarycomp/internal/repository"
	"github.com/ryannicoletti/veterinarycomp/internal/repository/dbrepo"
)

var Repo *Repository

type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

func NewHandlers(r *Repository) {
	Repo = r
}

func (repo *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIp := r.RemoteAddr
	repo.App.Session.Put(r.Context(), "remote_ip", remoteIp)
	render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
}

func (repo *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "about.page.tmpl", &models.TemplateData{})
}

func (repo *Repository) AddComp(w http.ResponseWriter, r *http.Request) {
	var emptyCompData models.Compensation
	data := make(map[string]interface{})
	data["compensation"] = emptyCompData
	render.Template(w, r, "add-comp.page.tmpl", &models.TemplateData{Form: forms.New(nil), Data: data})
}

func (repo *Repository) PostComp(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}
	form := forms.New(r.PostForm)

	baseSalary, err := form.StringToFloat("base-salary")
	if err != nil {
		log.Println("Error converting base salary to float:", err)
	}
	signOnBonus, err := form.StringToFloat("sign-on-bonus")
	if err != nil {
		log.Println("Error converting base salary to float:", err)
	}
	production, err := form.StringToFloat("production")
	if err != nil {
		log.Println("Error converting base salary to float:", err)
	}
	yearsExperience, err := form.StringToInt("years-experience")
	if err != nil {
		log.Println("Error converting years of experience to int")
	}

	totalComp := baseSalary + signOnBonus + production

	file, _, err := r.FormFile("verification-document")
	if err != nil {
		form.Errors.Add("document", "error getting verification document from form")
		log.Println("error getting verification document from form")
	}
	data, err := io.ReadAll(file)
	if err != nil {
		form.Errors.Add("document", "Error reading file for verification document")
		log.Println("Error reading file for verification document")
	}
	document := models.Document{
		Data:        data,
		FileName:    "verification",
		ContentType: http.DetectContentType(data),
		CreatedAt:   time.Now()}

	compensation := models.Compensation{
		CompanyName:          r.Form.Get("company-name"),
		Location:             r.Form.Get("location"),
		JobTitle:             r.Form.Get("job-title"),
		PracticeType:         r.Form.Get("type-of-practice"),
		BoardCertification:   r.Form.Get("board-certification"),
		YearsExperience:      yearsExperience,
		BaseSalary:           float32(baseSalary),
		SignOnBonus:          float32(signOnBonus),
		Production:           float32(production),
		TotalCompensation:    totalComp,
		VerificationDocument: document,
		Verified:             false,
		CreatedAt:            time.Now(),
	}

	form.Required("company-name", "location", "job-title", "base-salary")
	if !form.Valid() {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	repo.App.Session.Put(r.Context(), "compensation", compensation)
	http.Redirect(w, r, "/home", http.StatusCreated)
}
