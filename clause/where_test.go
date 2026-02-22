package clause

import (
	"testing"

	"github.com/suryaherdiyanto/sqlbuilder/dialect"
)

func TestWhereParsing(t *testing.T) {
	dialect := dialect.New("?", "`", "`")
	where := Where{Field: "name", Op: OperatorEqual, Value: "Alice"}

	stmt := where.Parse(dialect)
	expected := "`name` = ?"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}

	where2 := Where{Field: "age", Op: OperatorGreaterThan, Value: 30}
	stmt = where2.Parse(dialect)
	expected = "`age` > ?"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}

func TestWhereParsingPG(t *testing.T) {
	dialect := dialect.NewPostgres()
	where := Where{Field: "name", Op: OperatorEqual, Value: "Alice"}

	stmt := where.Parse(dialect)
	expected := `"name" = $1`

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}

	where2 := Where{Field: "age", Op: OperatorGreaterThan, Value: 30}
	stmt = where2.Parse(dialect)
	expected = `"age" > $2`

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}
