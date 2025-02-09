package sqlbuilder

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)


func TestBuilder(t *testing.T) {
	db, _ := sql.Open("sqlite3", ":memory:")

	builder := NewSQLBuilder("sqlite", db)

	builder.Select([]string{"*"}).Table("users")

	if builder.GetSql() != "SELECT * FROM users" {
		t.Errorf("Unexpected SQL result, got: %s", builder.GetSql())
	}
}
