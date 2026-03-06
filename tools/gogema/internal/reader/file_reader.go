package reader

import (
	"fmt"
	"os"

	"github.com/fahmihidayah/gogema/internal/model"
)

func IsFileAvailable(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// LoadProjectAndModel reads the project configuration and models
func LoadProjectAndModel(path string) (*model.Project, *[]model.Model) {
	project, err := LoadProject(path)
	if err != nil {
		fmt.Printf("Error reading project.yml: %v\n", err)
		return nil, nil
	}

	fmt.Printf("Project loaded: %+v\n", project)

	models, err := LoadModels(path)
	if err != nil {
		fmt.Printf("Error reading models: %v\n", err)
		return nil, nil
	}

	fmt.Printf("Models loaded: %+v\n", models)
	return project, models

}
