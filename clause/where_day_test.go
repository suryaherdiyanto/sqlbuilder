package clause

import (
	"testing"

	"github.com/suryaherdiyanto/sqlbuilder/dialect"
)

func TestWhereDayParsing(t *testing.T) {
	dialect := dialect.NewMySQL()
	where := WhereDay{Field: "created_at", Op: OperatorEqual, Value: 1}

	stmt := where.Parse(dialect)
	expected := "DAY(`created_at`) = ?"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}

	where2 := WhereDay{Field: "updated_at", Op: OperatorGreaterThan, Value: 1}
	stmt = where2.Parse(dialect)
	expected = "DAY(`updated_at`) > ?"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}

func TestWhereDayParsingPG(t *testing.T) {
	dialect := dialect.NewPostgres()
	where := WhereDay{Field: "created_at", Op: OperatorEqual, Value: 1}

	stmt := where.Parse(dialect)
	expected := "EXTRACT(DAY FROM \"created_at\") = ?"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}

	where2 := WhereDay{Field: "updated_at", Op: OperatorGreaterThan, Value: 1}
	stmt = where2.Parse(dialect)
	expected = "EXTRACT(DAY FROM \"updated_at\") > ?"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}
