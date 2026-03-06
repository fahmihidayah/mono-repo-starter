package reader

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/fahmihidayah/gogema/internal/model"
	"gopkg.in/yaml.v3"
)

func isDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func ReadYml(path, filename string) (string, error) {
	file := filepath.Join(path, filename)
	result, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

func ReadListModelFile(path, modelDir string) []string {
	modelPath := filepath.Join(path, modelDir)
	entries, err := os.ReadDir(modelPath)
	if err != nil {
		return nil
	}
	var files []string
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".yml") {
			continue
		}
		files = append(files, entry.Name())
	}
	return files
}

func LoadProject(path string) (*model.Project, error) {
	result, err := ReadYml(path, "project.yml")
	if err != nil {
		return nil, err
	}

	var config model.Project
	if err := yaml.Unmarshal([]byte(result), &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func LoadModels(path string) (*[]model.Model, error) {

	var listModels []model.Model

	files := ReadListModelFile(path, "model")
	for _, file := range files {
		m, err := ReadYmlToModel(filepath.Join(path, "model"), file)
		if err != nil {
			return nil, err
		}
		listModels = append(listModels, *m)
	}
	return &listModels, nil
}

func ReadYmlToModel(path, filename string) (*model.Model, error) {
	result, err := ReadYml(path, filename)
	if err != nil {
		return nil, err
	}

	var config model.Model
	if err := yaml.Unmarshal([]byte(result), &config); err != nil {
		return nil, err
	}
	return &config, nil
}
