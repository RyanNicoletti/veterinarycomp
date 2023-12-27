package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/justinas/nosurf"
	"github.com/ryannicoletti/veterinarycomp/internal/config"
	"github.com/ryannicoletti/veterinarycomp/internal/models"
)

var app *config.AppConfig

func NewRenderer(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.CSRFToken = nosurf.Token(r)
	if app.Session.Exists(r.Context(), "user_id") {
		td.IsAuthenticated = true
	} else {
		td.IsAuthenticated = false
	}
	return td
}

func Template(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) {
	var tc map[string]*template.Template
	if app.UseCache {
		tc = app.TemplateCache
		log.Println("pulling templates from template cache: ", tc)
		log.Println(app.UseCache)
	} else {
		tc, _ = CreateTemplateCache()
	}

	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("Could not get template from template cache")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td, r)
	// TODO: combine by just saying t.Execute(w, td) ???
	err := t.Execute(buf, td)
	if err != nil {
		log.Println(err)
	}
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println(err)
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		log.Println("failed to get names of template files", err)
		return cache, err
	}
	for _, page := range pages {
		name := filepath.Base(page)
		// ts = template set
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			log.Println("failed to parse template files", err)
			return cache, err
		}
		layouts, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			log.Println("failed to get names of layout files", err)
			return cache, err
		}
		if len(layouts) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				log.Println("failed to parse layout templates", err)
				return cache, err
			}
		}
		cache[name] = ts
	}
	return cache, nil
}
