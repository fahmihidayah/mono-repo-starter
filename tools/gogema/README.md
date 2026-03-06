# Gogema

A YAML-driven code generator for Go applications. Gogema reads project and model definitions from YAML files and generates boilerplate Go code including structs, database models, and more.

## Installation

```bash
go install github.com/fahmihidayah/gogema@latest
```

Or build from source:

```bash
git clone https://github.com/fahmihidayah/gogema.git
cd gogema
go build -o gogema .
```

## Usage

```bash
gogema generate --path ./my-project --framework golang
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--path` | `.` | Path to project configuration directory |
| `--framework` | `golang` | Target framework for code generation |

## Project Structure

Gogema expects the following directory structure:

```
my-project/
├── project.yml          # Project configuration
└── model/               # Model definitions
    ├── user.yml
    ├── post.yml
    └── ...
```

## Configuration

### project.yml

Define your project metadata:

```yaml
name: my-app
package: github.com/username/my-app
version: 1.0.0
author: Your Name
directory: ./output
```

### Model Definition

Create YAML files in the `model/` directory for each model:

```yaml
name: User
table: users
description: User account model

fields:
  - name: ID
    type: uint
    json: id
    db: id
    primary_key: true
    auto_increment: true

  - name: Email
    type: string
    json: email
    db: email
    required: true
    unique: true
    length: 255
    validation: email

  - name: Name
    type: string
    json: name
    db: name
    required: true
    length: 100

  - name: CreatedAt
    type: time.Time
    json: created_at
    db: created_at
    auto_now_add: true

  - name: UpdatedAt
    type: time.Time
    json: updated_at
    db: updated_at
    auto_now: true

indexes:
  - name: idx_users_email
    columns:
      - email
    unique: true

relationships:
  - name: Posts
    type: has_many
    model: Post
    foreign_key: user_id
    references: id
    on_delete: CASCADE
```

### Field Options

| Option | Type | Description |
|--------|------|-------------|
| `name` | string | Field name in Go struct |
| `type` | string | Go type (string, int, uint, bool, time.Time, etc.) |
| `json` | string | JSON tag name |
| `db` | string | Database column name |
| `primary_key` | bool | Mark as primary key |
| `auto_increment` | bool | Auto-increment field |
| `required` | bool | Field is required |
| `unique` | bool | Unique constraint |
| `nullable` | bool | Allow NULL values |
| `length` | int | Field length (for varchar) |
| `type_override` | string | Custom database type (text, jsonb, etc.) |
| `validation` | string | Validation rule |
| `default` | any | Default value |
| `auto_now_add` | bool | Set timestamp on creation |
| `auto_now` | bool | Update timestamp on modification |

### Foreign Keys

```yaml
fields:
  - name: AuthorID
    type: uint
    json: author_id
    db: author_id
    foreign_key:
      model: User
      field: id
      on_delete: CASCADE
      on_update: CASCADE
```

### Relationships

| Type | Description |
|------|-------------|
| `has_one` | One-to-one relationship |
| `has_many` | One-to-many relationship |
| `belongs_to` | Inverse of has_one/has_many |
| `many2many` | Many-to-many with join table |

```yaml
relationships:
  - name: Profile
    type: has_one
    model: Profile
    foreign_key: user_id
    references: id

  - name: Tags
    type: many2many
    model: Tag
    join_table: post_tags
```

## Requirements

- Go 1.24+

## Dependencies

- [cobra](https://github.com/spf13/cobra) - CLI framework
- [yaml.v3](https://gopkg.in/yaml.v3) - YAML parsing
- [color](https://github.com/fatih/color) - Terminal colors

## License

MIT
