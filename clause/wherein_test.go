package clause

import (
	"testing"

	"github.com/suryaherdiyanto/sqlbuilder/dialect"
)

func TestWhereInParsing(t *testing.T) {
	dialect := dialect.New("?", "`", "`")
	where := WhereIn{
		Field:  "name",
		Values: []any{"Alice", "Bob"},
	}

	stmt := where.Parse(dialect)
	expected := "`name` IN(?,?)"
	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}

	wherein2 := WhereIn{
		Field:  "age",
		Values: []any{25, 30, 35},
	}
	stmt = wherein2.Parse(dialect)
	expected = "`age` IN(?,?,?)"
	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}

func TestWhereInParsingPG(t *testing.T) {
	dialect := dialect.NewPostgres()
	where := WhereIn{
		Field:  "name",
		Values: []any{"Alice", "Bob"},
	}

	stmt := where.Parse(dialect)
	expected := "\"name\" IN($1,$2)"
	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}

	wherein2 := WhereIn{
		Field:  "age",
		Values: []any{25, 30, 35},
	}
	stmt = wherein2.Parse(dialect)
	expected = "\"age\" IN($3,$4,$5)"
	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}
