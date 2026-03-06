package renderer

import (
	"fmt"
	"testing"

	"github.com/fahmihidayah/gogema/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createSimpleProject() *model.Project {
	p := &model.Project{
		Name:      "Blog",
		Package:   " github.com/fahmihidayah/go-api",
		Version:   "1.0",
		Author:    "Fahmi",
		Directory: "../../blog-api/internal",
	}
	return p
}

func createSimpleUserModel() *model.Model {
	m := &model.Model{
		Name:        "User",
		Table:       "users",
		Description: "User entity",
		Fields: []model.Field{
			{
				Name:          "ID",
				PrimaryKey:    true,
				Type:          "uint",
				JSON:          "id",
				DB:            "id",
				AutoIncrement: true,
				Nullable:      false, // Pastikan field ini diisi untuk menghindari error evaluate
			},
			{
				Name:       "Email",
				Type:       "string",
				JSON:       "email",
				DB:         "email",
				Required:   true,
				Request:    true,
				Unique:     true,
				Nullable:   false,
				Validation: "required,email",
			},
			{
				Name:       "Name",
				Type:       "string",
				JSON:       "name",
				DB:         "name",
				Required:   true,
				Request:    true,
				Unique:     true,
				Nullable:   false,
				Validation: "required,min=2,max=100",
			},
			{
				Name:     "DeletedAt",
				Type:     "time.Time",
				JSON:     "deleted_at,omitempty",
				DB:       "deleted_at",
				Nullable: true, // Mengetes apakah pointer (*) muncul di hasil render
			},
		},
		// Tambahkan slice kosong agar range .Relations di template tidak error jika diakses
		Relationships: []model.Relationship{},
	}
	return m
}
func TestRenderModelSuccess(t *testing.T) {
	// 1. Siapkan data dummy yang lengkap sesuai dengan definisi struct Field terbaru
	p := createSimpleProject()
	m := createSimpleUserModel()

	// 2. Pastikan path template benar (disarankan gunakan path relatif yang aman)
	templatePath := "../../template/golang/domain/domain.go.tmpl"

	// 3. Eksekusi Render
	result, err := RenderModel(p, m, templatePath)

	// fmt.Printf("Data Result : %+v\n", result)
	// 4. Assertion
	require.NoError(t, err, "RenderModel harusnya tidak return error")
	assert.NotEmpty(t, result)

	// 5. Opsi tambahan: Cek apakah output mengandung string tertentu
	assert.Contains(t, result, "type User struct")
	assert.Contains(t, result, "gorm:\"column:email;uniqueIndex\"")
	assert.Contains(t, result, "*time.Time") // Memastikan logic Nullable bekerja
}

func TestRenderCreateRequestSuccess(t *testing.T) {
	m := createSimpleUserModel()
	p := createSimpleProject()

	templatePath := "../../template/golang/data/create.go.tmpl"
	result, err := RenderModel(p, m, templatePath)

	// fmt.Printf("Data Result : %+v\n", result)
	require.NoError(t, err, "RenderModel harusnya tidak return error")
	assert.NotEmpty(t, result)
}

func TestRenderRepositorySuccess(t *testing.T) {
	m := createSimpleUserModel()
	p := createSimpleProject()

	templatePath := "../../template/golang/repository/repository.go.tmpl"
	result, err := RenderModel(p, m, templatePath)

	// fmt.Printf("Data Result : %+v\n", result)
	require.NoError(t, err, "RenderModel harusnya tidak return error")
	assert.NotEmpty(t, result)
}

func TestRenderServiceSuccess(t *testing.T) {
	m := createSimpleUserModel()
	p := createSimpleProject()

	templatePath := "../../template/golang/service/service.go.tmpl"
	result, err := RenderModel(p, m, templatePath)

	// fmt.Printf("Data Result : %+v\n", result)
	require.NoError(t, err, "RenderModel harusnya tidak return error")
	assert.NotEmpty(t, result)
}

func TestRenderControllerSuccess(t *testing.T) {
	m := createSimpleUserModel()
	p := createSimpleProject()

	templatePath := "../../template/golang/controller/controller.go.tmpl"
	result, err := RenderModel(p, m, templatePath)

	// fmt.Printf("Controller Result : %+v\n", result)
	require.NoError(t, err, "RenderModel harusnya tidak return error")
	assert.NotEmpty(t, result)
}

func TestRenderHandlerSuccess(t *testing.T) {
	m := createSimpleUserModel()
	p := createSimpleProject()

	templatePath := "../../template/golang/server/server.go.tmpl"
	result, err := RenderModel(p, m, templatePath)

	fmt.Printf("Server Result : %+v\n", result)
	require.NoError(t, err, "RenderModel harusnya tidak return error")
	assert.NotEmpty(t, result)
}
