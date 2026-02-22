package clause

import (
	"testing"

	"github.com/suryaherdiyanto/sqlbuilder/dialect"
)

func TestSimpleSelectStatement(t *testing.T) {
	dialect := dialect.New("?", "`", "`")
	statement := Select{
		Columns: []string{"*"},
		Table:   "users",
	}
	where := WhereStatements{
		Where: []Where{
			{
				Field: "email",
				Value: "johndoe@gmail.com",
				Op:    OperatorEqual,
			},
		},
	}

	stmt, _ := statement.Parse(dialect)
	stmt += where.Parse(dialect)
	expected := "SELECT * FROM `users` WHERE `email` = ?"
	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}
