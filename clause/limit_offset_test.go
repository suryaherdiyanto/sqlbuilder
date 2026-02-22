package clause

import (
	"testing"

	"github.com/suryaherdiyanto/sqlbuilder/dialect"
)

func TestSimpleStatementWithLimitAndOffset(t *testing.T) {
	dialect := dialect.New("?", "`", "`")
	statement := Select{
		Table:   "users",
		Columns: []string{"id", "email", "name"},
	}
	limit := Limit{Count: 10}
	offset := Offset{Count: 5}

	stmt, _ := statement.Parse(dialect)
	stmt += limit.Parse(dialect)
	stmt += offset.Parse(dialect)
	expected := "SELECT `id`,`email`,`name` FROM `users` LIMIT ? OFFSET ?"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}

func TestSimpleStatementWithLimitAndOffsetPG(t *testing.T) {
	dialect := dialect.NewPostgres()
	statement := Select{
		Table:   "users",
		Columns: []string{"id", "email", "name"},
	}
	limit := Limit{Count: 10}
	offset := Offset{Count: 5}

	stmt, _ := statement.Parse(dialect)
	stmt += limit.Parse(dialect)
	stmt += offset.Parse(dialect)
	expected := "SELECT \"id\",\"email\",\"name\" FROM \"users\" LIMIT $1 OFFSET $2"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}
