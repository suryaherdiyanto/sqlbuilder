package sqlbuilder

import (
	"testing"
)

func TestWhereInParsing(t *testing.T) {
	whereIn := WhereIn{
		Field:  "name",
		Values: []any{"john", "doe"},
	}

	stmt := whereIn.Parse()
	expected := "name IN('john','doe')"
	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}

	whereIn = WhereIn{
		Field:  "age",
		Values: []any{19, 20, 31},
	}
	stmt = whereIn.Parse()
	expected = "age IN(19,20,31)"
	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}

func TestWhereParsing(t *testing.T) {
	where := Where{
		Field: "name",
		Op:    OperatorEqual,
		Value: "John Doe",
	}

	stmt := where.Parse()
	expected := "name = 'John Doe'"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}

	where = Where{
		Field: "age",
		Op:    OperatorGreaterThan,
		Value: 17,
	}

	stmt = where.Parse()
	expected = "age > 17"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}

func TestOrderClause(t *testing.T) {
	order := Order{
		Field:     "name",
		Direction: OrderDirectionDESC,
	}

	stmt := order.Parse()
	expected := "ORDER BY name desc"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}

func TestSimpleSelectStatement(t *testing.T) {
	statement := SelectStatement{
		Table:   "users",
		Columns: []string{"id", "email", "name"},
		WhereStatements: []Where{
			{
				Field: "email",
				Value: "johndoe@gmail.com",
				Op:    OperatorEqual,
			},
		},
	}

	stmt := statement.Parse()
	expected := "SELECT id,email,name FROM users WHERE email = 'johndoe@gmail.com'"
	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}

}
