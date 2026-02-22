package clause

import (
	"testing"

	"github.com/suryaherdiyanto/sqlbuilder/dialect"
)

func TestWhereBetweenParsing(t *testing.T) {
	dialect := dialect.New("?", "`", "`")
	where := WhereBetween{Field: "age", Start: 18, End: 30}

	stmt := where.Parse(dialect)
	expected := "`age` BETWEEN ? AND ?"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}

func TestWhereBetweenParsingPG(t *testing.T) {
	dialect := dialect.NewPostgres()
	where := WhereBetween{Field: "age", Start: 18, End: 30}

	stmt := where.Parse(dialect)
	expected := "\"age\" BETWEEN $1 AND $2"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}
