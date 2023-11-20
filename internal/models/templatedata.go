package models

// holds data send from handlers to templates
type TemplateData struct {
	Data         map[string]interface{}
	FlashMessage string
	Warning      string
	Error        string
	CSRFToken    string
}
