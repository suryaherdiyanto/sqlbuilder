package clause

import (
	"testing"

	"github.com/suryaherdiyanto/sqlbuilder/dialect"
)

func TestDeleteStatement(t *testing.T) {
	dialect := dialect.New("?", "`", "`")
	statement := Delete{
		Table: "users",
	}
	where := WhereStatements{
		Where: []Where{
			{
				Field: "id",
				Op:    OperatorEqual,
				Value: 1,
			},
		},
	}

	stmt, _ := statement.Parse(dialect)
	stmt += where.Parse(dialect)
	expected := "DELETE FROM `users` WHERE `id` = ?"
	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}
