package handlers

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/ryannicoletti/veterinarycomp/internal/config"
	"github.com/ryannicoletti/veterinarycomp/internal/driver"
	"github.com/ryannicoletti/veterinarycomp/internal/forms"
	"github.com/ryannicoletti/veterinarycomp/internal/helpers"
	"github.com/ryannicoletti/veterinarycomp/internal/models"
	"github.com/ryannicoletti/veterinarycomp/internal/render"
	"github.com/ryannicoletti/veterinarycomp/internal/repository"
	"github.com/ryannicoletti/veterinarycomp/internal/repository/repositoryimpl"
)

var Repo *Repository

type Repository struct {
	App                *config.AppConfig
	CompensationDBRepo repository.CompensationRepo
	UserDBRepo         repository.UserRepo
}

func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App:                a,
		CompensationDBRepo: repositoryimpl.NewPostgresCompensationRepo(db.SQL, a),
		UserDBRepo:         repositoryimpl.NewPostgresUserRepo(db.SQL, a),
	}
}

func NewHandlers(r *Repository) {
	Repo = r
}

func (repo *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIp := r.RemoteAddr
	repo.App.Session.Put(r.Context(), "remote_ip", remoteIp)
	data := make(map[string]interface{})
	page, ok := r.Context().Value("page").(int)
	if !ok || page == 0 {
		page = 1
	}
	rowPerPage := 10

	c, err := Repo.CompensationDBRepo.GetAllApprovedCompensations()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data["compensations"] = c

	total, err := Repo.CompensationDBRepo.GetApprovedCompensationsCount()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	t := total / rowPerPage
	remainder := total % rowPerPage
	var totalPages int
	if remainder == 0 {
		totalPages = t
	} else {
		totalPages = t + 1
	}

	paginationData := models.Pagination{
		Next:          page + 1,
		Previous:      page - 1,
		RecordPerPage: rowPerPage,
		CurrentPage:   page,
		TotalPage:     totalPages,
	}

	data["page"] = paginationData

	if _, ok := repo.App.Session.Get(r.Context(), "compensation").(models.Compensation); ok {
		repo.App.Session.Put(r.Context(), "flash", "Thank you, your submission will be reviewed as soon as possible.")
		repo.App.Session.Remove(r.Context(), "compensation")
	}

	render.Template(w, r, "home.page.tmpl", &models.TemplateData{Data: data})
}

func (repo *Repository) SearchComp(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}
	rowPerPage := 10
	offset := rowPerPage * (page - 1)
	locationOrHospital := r.URL.Query().Get("location-hospital")
	data := make(map[string]interface{})
	c, err := Repo.CompensationDBRepo.SearchCompensation(locationOrHospital, rowPerPage, offset)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data["compensations"] = c
	total, err := Repo.CompensationDBRepo.GetTotalSearchCompensationsCount(locationOrHospital)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	t := total / rowPerPage
	remainder := total % rowPerPage
	var totalPages int
	if remainder == 0 {
		totalPages = t
	} else {
		totalPages = t + 1
	}

	paginationData := models.Pagination{
		Next:          page + 1,
		Previous:      page - 1,
		RecordPerPage: rowPerPage,
		CurrentPage:   page,
		TotalPage:     totalPages,
	}

	data["page"] = paginationData
	data["locationOrHospital"] = locationOrHospital
	render.Template(w, r, "home.page.tmpl", &models.TemplateData{Data: data, IsSearchPerformed: true})
}

func (repo *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "about.page.tmpl", &models.TemplateData{})
}
func (repo *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.page.tmpl", &models.TemplateData{})
}

func (repo *Repository) Login(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "login.page.tmpl", &models.TemplateData{Form: forms.NewForm(nil)})
}

func (repo *Repository) PostLogin(w http.ResponseWriter, r *http.Request) {
	_ = repo.App.Session.RenewToken(r.Context())
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}
	email := r.Form.Get("email")
	password := r.Form.Get("password")
	form := forms.NewForm(r.PostForm)
	form.Required("email", "password")
	form.IsEmail(email)
	if !form.Valid() {
		render.Template(w, r, "login.page.tmpl", &models.TemplateData{Form: form})
		return
	}
	id, _, err := repo.UserDBRepo.Authenticate(email, password)
	if err != nil {
		repo.App.Session.Put(r.Context(), "error", "Invalid login credentials")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	repo.App.Session.Put(r.Context(), "user_id", id)

	isAdmin, _ := repo.UserDBRepo.IsAdmin(id)
	if isAdmin {
		repo.App.Session.Put(r.Context(), "is_admin", isAdmin)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (repo *Repository) Logout(w http.ResponseWriter, r *http.Request) {
	_ = repo.App.Session.Destroy(r.Context())
	_ = repo.App.Session.RenewToken(r.Context())
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (repo *Repository) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	c, err := Repo.CompensationDBRepo.GetAllUnapprovedCompensations(100, 0)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data := make(map[string]interface{})
	// sort comp data so rows with verification document appear at the top
	sort.Slice(c, func(i, j int) bool {
		if c[i].VerificationDocument != nil && c[j].VerificationDocument == nil {
			return true
		} else if c[i].VerificationDocument == nil && c[j].VerificationDocument != nil {
			return false
		}
		return c[i].CreatedAt.After(c[j].CreatedAt)
	})
	data["compensations"] = c
	render.Template(w, r, "admin-dashboard.page.tmpl", &models.TemplateData{Data: data})
}

func (repo *Repository) DownloadVerification(w http.ResponseWriter, r *http.Request) {
	ID, _ := strconv.Atoi(r.URL.Query().Get("ID"))
	c, err := Repo.CompensationDBRepo.GetDocumentMetaDataById(ID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	fp := c.VerificationDocument.FilePath
	w.Header().Set("Content-Type", c.VerificationDocument.ContentType)
	w.Header().Set("Content-Disposition", "attachment; filename="+c.VerificationDocument.FileName)
	f, err := os.ReadFile(fp)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	_, err = w.Write(f)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
}

func (repo *Repository) VerifyComp(w http.ResponseWriter, r *http.Request) {
	ID, _ := strconv.Atoi(r.URL.Query().Get("ID"))
	err := Repo.CompensationDBRepo.VerifyComp(ID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	w.Write([]byte(`<span class="checkmark">&#9745;</span>`))
}

func (repo *Repository) ApproveComp(w http.ResponseWriter, r *http.Request) {
	ID, _ := strconv.Atoi(r.URL.Query().Get("ID"))
	err := Repo.CompensationDBRepo.ApproveComp(ID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	return
}

func (repo *Repository) DeleteComp(w http.ResponseWriter, r *http.Request) {
	ID, _ := strconv.Atoi(r.URL.Query().Get("ID"))
	err := Repo.CompensationDBRepo.DeleteCompensationByID(ID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
}

func (repo *Repository) DeleteCompDocument(w http.ResponseWriter, r *http.Request) {
	ID, _ := strconv.Atoi(r.URL.Query().Get("ID"))
	c, err := Repo.CompensationDBRepo.GetCompensationByID(ID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	if c.VerificationDocument == nil {
		return
	}
	// delete metadata
	err = Repo.CompensationDBRepo.DeleteCompensationDocumentByID(ID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	// delete file
	err = os.Remove(c.VerificationDocument.FilePath)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	w.Write([]byte(`<span>No verification provided</span>`))
}

func (repo *Repository) CompForm(w http.ResponseWriter, r *http.Request) {
	var emptyCompData models.Compensation
	data := make(map[string]interface{})
	data["compensation"] = emptyCompData
	render.Template(w, r, "add-comp.page.tmpl", &models.TemplateData{Form: forms.NewForm(nil), Data: data})
}

func (repo *Repository) PostCompForm(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(20 << 30)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	form := forms.NewForm(r.PostForm)
	form.TrimMoneyvalue("base-salary", "production", "sign-on-bonus", "hourly-rate")
	form.Required("company-name", "job-title", "is-veterinarian", "years-experience", "country", "state")
	signOnBonus, _ := form.StringToFloat("sign-on-bonus")
	production, _ := form.StringToFloat("production")
	yearsExperience, _ := form.StringToInt("years-experience")

	isVetStr := form.Get("is-veterinarian")
	isVet, err := strconv.ParseBool(isVetStr)
	if err != nil {
		isVet = false
	}

	isHourlyStr := form.Get("is-hourly")
	isHourly, err := strconv.ParseBool(isHourlyStr)
	if err != nil {
		isHourly = false
	}

	var hourlyRate float64
	var baseSalary float64
	var totalComp float64
	if isHourly {
		hourlyRate, _ = form.StringToFloat("hourly-rate")
		form.Required("hourly-rate")
		baseSalary = 0.0
	} else {
		form.Required("base-salary")
		baseSalary, _ = form.StringToFloat("base-salary")
		totalComp = baseSalary + signOnBonus + production
		hourlyRate = 0.0
	}

	var document *models.Document

	if files, ok := r.MultipartForm.File["verification-document"]; ok && len(files) > 0 {
		fileHeader := files[0]
		fileData, err := form.ProcessFileUpload("verification-document", fileHeader)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}
		fileName := uuid.NewString()
		filePath := filepath.Join("uploads", fileName)
		err = os.WriteFile(filePath, fileData, 0644)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}
		document = &models.Document{
			ID:          uuid.NewString(),
			FileName:    fileHeader.Filename,
			ContentType: http.DetectContentType(fileData),
			FilePath:    filePath,
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
		Approved:             false,
		IsHourly:             isHourly,
		HourlyRate:           float32(hourlyRate),
		CreatedAt:            time.Now(),
		IsVeterinarian:       isVet,
		Country:              r.Form.Get("country"),
		State:                r.Form.Get("state"),
		City:                 r.Form.Get("city"),
	}

	if !form.Valid() {
		data := make(map[string]interface{})
		data["compensation"] = compensation
		render.Template(w, r, "add-comp.page.tmpl", &models.TemplateData{Form: form, Data: data})
		return
	}
	err = Repo.CompensationDBRepo.InsertCompensation(compensation)
	if err != nil {
		helpers.ServerError(w, err)
		// return here? needs testing
	}
	repo.App.Session.Put(r.Context(), "compensation", compensation) // might not need to put this in a session...
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
