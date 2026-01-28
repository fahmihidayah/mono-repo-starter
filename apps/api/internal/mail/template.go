package mail

import (
	"bytes"
	"html/template"
	"path/filepath"
)

type TemplateEngine struct {
	basePath string
}

func TemplateEngineProvider() *TemplateEngine {
	// Default to templates directory in mail package
	basePath := "internal/mail/templates"
	return &TemplateEngine{
		basePath: basePath,
	}
}

func (te *TemplateEngine) Render(name string, data any) (string, error) {
	path := filepath.Join(te.basePath, name)

	tmpl, err := template.ParseFiles(path)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer

	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
