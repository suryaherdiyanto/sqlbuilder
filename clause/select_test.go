package clause

import (
	"testing"

	"github.com/suryaherdiyanto/sqlbuilder/dialect"
)

func TestSimpleSelectStatement(t *testing.T) {
	dialect := dialect.New("?", "`", "`")
	statement := Select{
		Columns: []string{"*"},
		Table:   "`users`",
	}
	where := Where{
		Field: "email",
		Value: "johndoe@gmail.com",
		Op:    OperatorEqual,
	}

	stmt, _ := statement.Parse(dialect)
	stmt += " WHERE " + where.Parse(dialect)
	expected := "SELECT * FROM `users` WHERE `email` = ?"
	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}

func TestSimpleSelectStatementPG(t *testing.T) {
	dialect := dialect.NewPostgres()
	statement := Select{
		Columns: []string{"*"},
		Table:   "\"users\"",
	}
	where := Where{
		Field: "email",
		Value: "johndoe@gmail.com",
		Op:    OperatorEqual,
	}

	stmt, _ := statement.Parse(dialect)
	stmt += " WHERE " + where.Parse(dialect)
	expected := "SELECT * FROM \"users\" WHERE \"email\" = $1"
	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}
