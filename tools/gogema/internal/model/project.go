package model

// Project represents the project configuration
type Project struct {
	Name      string `yaml:"name"`
	Package   string `yaml:"package"`
	Version   string `yaml:"version"`
	Author    string `yaml:"author"`
	Directory string `yaml:"directory"`
}

// Field represents a model field
type Field struct {
	Name          string      `yaml:"name"`
	Type          string      `yaml:"type"`
	JSON          string      `yaml:"json"`
	DB            string      `yaml:"db"`
	Request       bool        `yaml:"request"`
	PrimaryKey    bool        `yaml:"primary_key"`
	AutoIncrement bool        `yaml:"auto_increment"`
	Required      bool        `yaml:"required"`
	Unique        bool        `yaml:"unique"`
	Nullable      bool        `yaml:"nullable"`      // Tambahkan ini (Penyebab error template)
	Length        int         `yaml:"length"`        // Tambahkan ini untuk varchar(n)
	TypeOverride  string      `yaml:"type_override"` // Tambahkan ini untuk custom type (text, jsonb)
	Validation    string      `yaml:"validation"`
	Default       interface{} `yaml:"default"`
	AutoNowAdd    bool        `yaml:"auto_now_add"` // Tambahkan ini (CreatedAt)
	AutoNow       bool        `yaml:"auto_now"`     // Tambahkan ini (UpdatedAt)
	ForeignKey    *ForeignKey `yaml:"foreign_key"`
}

// ForeignKey represents a foreign key relationship (Inline Field Constraint)
type ForeignKey struct {
	Model    string `yaml:"model"`
	Field    string `yaml:"field"`
	OnDelete string `yaml:"on_delete"`
	OnUpdate string `yaml:"on_update"` // Tambahkan ini agar lebih lengkap
}

// Index represents a database index
type Index struct {
	Name    string   `yaml:"name"`
	Columns []string `yaml:"columns"` // Sesuaikan dengan YAML sebelumnya (tadi Anda pakai 'columns' di YAML)
	Unique  bool     `yaml:"unique"`
}

// Relationship represents model relationships (Association)
type Relationship struct {
	Name       string `yaml:"name"` // Tambahkan ini untuk nama field di Go (e.g., Author)
	Type       string `yaml:"type"`
	Model      string `yaml:"model"`
	ForeignKey string `yaml:"foreign_key"`
	References string `yaml:"references"`
	JoinTable  string `yaml:"join_table"` // Tambahkan ini untuk many2many
	OnDelete   string `yaml:"on_delete"`  // Tambahkan agar bisa diakses template relasi
	OnUpdate   string `yaml:"on_update"`  // Tambahkan agar bisa diakses template relasi
}

// Model represents a complete model definition
type Model struct {
	Name          string         `yaml:"name"`
	Table         string         `yaml:"table"`
	Description   string         `yaml:"description"`
	Fields        []Field        `yaml:"fields"`
	Indexes       []Index        `yaml:"indexes"`
	Relationships []Relationship `yaml:"relationships"`
	Imports       []string       `yaml:"imports"`
	Project       *Project       `yaml:"project"`
}
