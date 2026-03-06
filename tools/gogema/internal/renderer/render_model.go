package renderer

import (
	"bytes"
	"fmt"
	"go/format"
	"html/template"
	"path/filepath"
	"strings"

	"github.com/fahmihidayah/gogema/internal/model"
	"github.com/fahmihidayah/gogema/internal/reader"
)

func RenderModel(p *model.Project, m *model.Model, templateFile string) (string, error) {
	if !reader.IsFileAvailable(templateFile) {
		return "", fmt.Errorf("template file not found")
	}

	// 1. Definisikan Map fungsi yang akan digunakan di dalam template
	funcMap := template.FuncMap{
		"lower": strings.ToLower,
		"upper": strings.ToUpper,
		"removeWhiteSpace": func(s string) string {
			return strings.ReplaceAll(s, " ", "")
		},
		// Anda bisa menambah fungsi lain seperti snake_case jika perlu
	}

	SanitizeModel(m)

	m.Project = p

	var buf bytes.Buffer

	// 2. Gunakan New().Funcs().ParseFiles() untuk mendaftarkan funcMap sebelum parsing
	tmpl, err := template.New(filepath.Base(templateFile)).Funcs(funcMap).ParseFiles(templateFile)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %v", err)
	}

	if err := tmpl.Execute(&buf, m); err != nil {
		return "", err
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return buf.String(), fmt.Errorf("failed to format code: %v", err)
	}

	return string(formatted), nil
}
