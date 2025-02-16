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

func TestWithWhere(t *testing.T) {
	db, _ := sql.Open("sqlite3", ":memory:")

	builder := NewSQLBuilder("sqlite", db)

	builder.
		Select([]string{"*"}).
		Table("users").
		Where("email", Eq, "johndoe@gmail.com")

	if builder.GetSql() != "SELECT * FROM users WHERE email = 'johndoe@gmail.com'" {
		t.Errorf("Unexpected SQL result, got: %s", builder.GetSql())
	}
}

func TestWithMultipleWhere(t *testing.T) {
	db, _ := sql.Open("sqlite3", ":memory:")

	builder := NewSQLBuilder("sqlite", db)

	builder.
		Select([]string{"*"}).
		Table("users").
		Where("email", Eq, "johndoe@gmail.com").
		Where("access_role", Lt, 3)

	if builder.GetSql() != "SELECT * FROM users WHERE email = 'johndoe@gmail.com' AND access_role < 3" {
		t.Errorf("Unexpected SQL result, got: %s", builder.GetSql())
	}
}
