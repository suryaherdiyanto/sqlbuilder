package sqlbuilder

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)


func TestBuilder(t *testing.T) {
	db, _ := sql.Open("sqlite3", ":memory:")
	builder := NewSQLBuilder("sqlite", db)

	builder.Table("users").Select([]string{"*"})
}
