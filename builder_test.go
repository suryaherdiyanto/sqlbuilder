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

func TestWhereIn(t *testing.T) {
	db, _ := sql.Open("sqlite3", ":memory:")

	builder := NewSQLBuilder("sqlite", db)

	builder.
		Select([]string{"*"}).
		Table("users").
		WhereIn("email", []string{"john@gmail.com", "admin@example.com"})

	if builder.GetSql() != "SELECT * FROM users WHERE email IN('john@gmail.com', 'admin@example.com')" {
		t.Errorf("Unexpected SQL result, got: %s", builder.GetSql())
	}
}

func TestWhereInWithNumber(t *testing.T) {
	db, _ := sql.Open("sqlite3", ":memory:")

	builder := NewSQLBuilder("sqlite", db)

	builder.
		Select([]string{"*"}).
		Table("users").
		WhereIn("some_column", []int{1, 2, 3})

	if builder.GetSql() != "SELECT * FROM users WHERE some_column IN(1, 2, 3)" {
		t.Errorf("Unexpected SQL result, got: %s", builder.GetSql())
	}
}

func TestWhereBetween(t *testing.T) {
	db, _ := sql.Open("sqlite3", ":memory:")

	builder := NewSQLBuilder("sqlite", db)

	builder.
		Select([]string{"*"}).
		Table("users").
		WhereBetween("age", 5, 10)

	if builder.GetSql() != "SELECT * FROM users WHERE age BETWEEN 5 AND 10" {
		t.Errorf("Unexpected SQL result, got: %s", builder.GetSql())
	}
}
