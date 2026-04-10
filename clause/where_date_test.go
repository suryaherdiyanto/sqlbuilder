package clause

import (
	"testing"

	"github.com/suryaherdiyanto/sqlbuilder/dialect"
)

func TestWhereDateParsing(t *testing.T) {
	dialect := dialect.New("?", "`", "`")
	where := WhereDate{Field: "created_at", Op: OperatorEqual, Value: "2024-01-01"}

	stmt := where.Parse(dialect)
	expected := "DATE(`created_at`) = ?"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}

	where2 := WhereDate{Field: "updated_at", Op: OperatorGreaterThan, Value: "2024-01-01"}
	stmt = where2.Parse(dialect)
	expected = "DATE(`updated_at`) > ?"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}

func TestWhereDateParsingPG(t *testing.T) {
	dialect := dialect.NewPostgres()
	where := WhereDate{Field: "created_at", Op: OperatorEqual, Value: "2024-01-01"}

	stmt := where.Parse(dialect)
	expected := "CAST(\"created_at\" AS DATE) = ?"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}

	where2 := WhereDate{Field: "updated_at", Op: OperatorGreaterThan, Value: "2024-01-01"}
	stmt = where2.Parse(dialect)
	expected = "CAST(\"updated_at\" AS DATE) > ?"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}
