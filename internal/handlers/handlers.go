package handlers

import (
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
	render.Template(w, r, "add-comp.page.tmpl", &models.TemplateData{Form: forms.NewForm(nil), Data: data})
}

func (repo *Repository) PostComp(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(20 << 30)
	if err != nil {
		log.Println(err)
		return
	}
	form := forms.NewForm(r.PostForm)
	form.Required("company-name", "location", "job-title", "base-salary")
	baseSalary, _ := form.StringToFloat("base-salary")
	signOnBonus, _ := form.StringToFloat("sign-on-bonus")
	production, _ := form.StringToFloat("production")
	yearsExperience, _ := form.StringToInt("years-experience")
	totalComp := baseSalary + signOnBonus + production

	var document *models.Document

	if files, ok := r.MultipartForm.File["verification-document"]; ok && len(files) > 0 {
		fileHeader := files[0]
		verificationData, _ := form.ProcessFileUpload("verification-document", fileHeader)
		document = &models.Document{
			Data:        verificationData,
			FileName:    "verification",
			ContentType: http.DetectContentType(verificationData),
			CreatedAt:   time.Now(),
		}

	} else {
		document = nil
	}

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

	if !form.Valid() {
		data := make(map[string]interface{})
		data["compensation"] = compensation
		render.Template(w, r, "add-comp.page.tmpl", &models.TemplateData{Form: form, Data: data})
		return
	}
	repo.App.Session.Put(r.Context(), "compensation", compensation) // might not need to put this in a session...
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
