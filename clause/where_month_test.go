package clause

import (
	"testing"

	"github.com/suryaherdiyanto/sqlbuilder/dialect"
)

func TestWhereMonthParsing(t *testing.T) {
	dialect := dialect.NewMySQL()
	where := WhereMonth{Field: "created_at", Op: OperatorEqual, Value: 1}

	stmt := where.Parse(dialect)
	expected := "MONTH(`created_at`) = ?"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}

	where2 := WhereMonth{Field: "updated_at", Op: OperatorGreaterThan, Value: 1}
	stmt = where2.Parse(dialect)
	expected = "MONTH(`updated_at`) > ?"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}

func TestWhereMonthParsingPG(t *testing.T) {
	dialect := dialect.NewPostgres()
	where := WhereMonth{Field: "created_at", Op: OperatorEqual, Value: 1}

	stmt := where.Parse(dialect)
	expected := "EXTRACT(MONTH FROM \"created_at\") = ?"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}

	where2 := WhereMonth{Field: "updated_at", Op: OperatorGreaterThan, Value: 1}
	stmt = where2.Parse(dialect)
	expected = "EXTRACT(MONTH FROM \"updated_at\") > ?"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}
