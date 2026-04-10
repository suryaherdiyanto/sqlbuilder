package clause

import (
	"testing"

	"github.com/suryaherdiyanto/sqlbuilder/dialect"
)

func TestWhereNotBetweenParsing(t *testing.T) {
	dialect := dialect.New("?", "`", "`")
	where := WhereNotBetween{Field: "age", Start: 18, End: 30}

	stmt := where.Parse(dialect)
	expected := "`age` NOT BETWEEN ? AND ?"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}

func TestWhereNotBetweenParsingPG(t *testing.T) {
	dialect := dialect.NewPostgres()
	where := WhereNotBetween{Field: "age", Start: 18, End: 30}

	stmt := where.Parse(dialect)
	expected := "\"age\" NOT BETWEEN $1 AND $2"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}
