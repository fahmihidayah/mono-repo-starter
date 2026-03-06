package renderer

import "github.com/fahmihidayah/gogema/internal/model"

func SanitizeModel(model *model.Model) {
	model.Imports = []string{}
	for i := range model.Fields {
		if model.Fields[i].Type == "time.Time" {
			model.Imports = append(model.Imports, "time")
		}
	}
}
