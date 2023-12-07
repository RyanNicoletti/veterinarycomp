package models

import (
	"github.com/ryannicoletti/veterinarycomp/internal/forms"
)

// holds data send from handlers to templates
type TemplateData struct {
	Data              map[string]interface{}
	IsSearchPerformed bool
	FlashMessage      string
	Warning           string
	Error             string
	CSRFToken         string
	Form              *forms.Form
	IsAuthenticated   bool
}
