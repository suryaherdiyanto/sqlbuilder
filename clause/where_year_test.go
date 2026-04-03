package clause

import (
	"testing"

	"github.com/suryaherdiyanto/sqlbuilder/dialect"
)

func TestWhereYearParsing(t *testing.T) {
	dialect := dialect.NewMySQL()
	where := WhereYear{Field: "created_at", Op: OperatorEqual, Value: 2024}

	stmt := where.Parse(dialect)
	expected := "YEAR(`created_at`) = ?"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}

	where2 := WhereYear{Field: "updated_at", Op: OperatorGreaterThan, Value: 2024}
	stmt = where2.Parse(dialect)
	expected = "YEAR(`updated_at`) > ?"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}

func TestWhereYearParsingPG(t *testing.T) {
	dialect := dialect.NewPostgres()
	where := WhereYear{Field: "created_at", Op: OperatorEqual, Value: 2024}

	stmt := where.Parse(dialect)
	expected := "EXTRACT(YEAR FROM \"created_at\") = ?"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}

	where2 := WhereYear{Field: "updated_at", Op: OperatorGreaterThan, Value: 2024}
	stmt = where2.Parse(dialect)
	expected = "EXTRACT(YEAR FROM \"updated_at\") > ?"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}
