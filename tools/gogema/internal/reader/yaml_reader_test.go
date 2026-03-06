package reader

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadYmlSuccess(t *testing.T) {
	result, err := ReadYml("../../example", "project.yml")

	require.NoError(t, err)
	assert.NotEmpty(t, result)
}

func TestReadYmlToProjectSuccess(t *testing.T) {
	project, err := LoadProject("../../example")
	fmt.Printf("project: %+v\n", project)
	require.NoError(t, err)
	require.NotNil(t, project)
	assert.Equal(t, project.Name, "BlogGo")
}

func TestGetListModelFileSuccess(t *testing.T) {
	files := ReadListModelFile("../../example", "model")
	assert.Equal(t, len(files), 5)
	assert.Contains(t, files[0], "category.yml")
}

func TestReadYmlToModelSuccess(t *testing.T) {
	model, err := ReadYmlToModel("../../example/model", "comment.yml")
	require.NoError(t, err)
	require.NotNil(t, model)
	assert.Equal(t, model.Name, "Comment")
}
