package clause

import (
	"testing"

	"github.com/suryaherdiyanto/sqlbuilder/dialect"
)

func TestWhereNotInParsing(t *testing.T) {
	dialect := dialect.New("?", "`", "`")
	where := WhereNotIn{Field: "id", Values: []any{1, 2}}

	stmt := where.Parse(dialect)
	expected := "`id` NOT IN(?,?)"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}

func TestWhereNotInParsingPG(t *testing.T) {
	dialect := dialect.NewPostgres()
	where := WhereNotIn{Field: "id", Values: []any{1, 2}}

	stmt := where.Parse(dialect)
	expected := "\"id\" NOT IN($1,$2)"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}
