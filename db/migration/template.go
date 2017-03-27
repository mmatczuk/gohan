package migration

import (
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"text/template"
)

var goMigrationTemplate = template.Must(template.New("goose.go-migration").Parse(`
package migration

import (
    "database/sql"

    "github.com/cloudwan/gohan/db/migration"
)

func init() {
    migration.AddMigration(Up_{{.}}, Down_{{.}})
}

func Up_{{.}}(tx *sql.Tx) error {
    return nil
}

func Down_{{.}}(tx *sql.Tx) error {
    return nil
}
`))

var sqlMigrationTemplate = template.Must(template.New("goose.sql-migration").Parse(`
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

`))

func createTemplate(w io.Writer, migrationType string, version int64) error {
	var tmpl *template.Template
	switch migrationType {
	case "go":
		tmpl = goMigrationTemplate
	case "sql":
		tmpl = sqlMigrationTemplate
	default:
		return errors.New("unsupported migration type")
	}

	return tmpl.Execute(w, version)
}

func templateFile(dir, name, migrationType string, version int64) (path string) {
	filename := fmt.Sprintf("%d_%s.%s", version, name, migrationType)
	return filepath.Join(dir, filename)
}
