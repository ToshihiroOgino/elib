package secure

import (
	"html/template"
	"path/filepath"
)

func createSecureTemplateEngine(pattern string) *template.Template {
	funcMap := template.FuncMap{
		"escapeHTML": func(text string) template.HTML {
			return escapeHTML(text)
		},
		"sanitizeHTML": func(text string) string {
			return sanitizeHTML(text)
		},
		"safeJSON": func(text string) string {
			return safeJSONString(text)
		},
	}

	template := template.New("").Funcs(funcMap)
	template, err := template.ParseGlob(pattern)
	if err != nil {
		panic(err)
	}

	return template
}

func LoadSecureTemplates() *template.Template {
	pattern := filepath.Join("templates", "*", "*")
	return createSecureTemplateEngine(pattern)
}
