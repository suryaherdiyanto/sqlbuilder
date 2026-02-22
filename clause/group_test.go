package clause

import (
	"testing"

	"github.com/suryaherdiyanto/sqlbuilder/dialect"
)

func TestWithGroupByStatement(t *testing.T) {
	dialect := dialect.New("?", "`", "`")
	statement := Select{
		Table:   "users",
		Columns: []string{"role", "COUNT(*) as total"},
	}
	grouping := GroupBy{
		Fields: []string{"role"},
	}

	stmt, _ := statement.Parse(dialect)
	stmt += grouping.Parse(dialect)
	expected := "SELECT `role`,COUNT(*) as total FROM `users` GROUP BY `role`"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}
