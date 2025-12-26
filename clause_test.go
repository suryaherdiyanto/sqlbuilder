package sqlbuilder

import "testing"

func TestWhereInParsing(t *testing.T) {
	whereIn := WhereIn{
		Field:  "name",
		Values: []any{"john", "doe"},
	}

	stmt := whereIn.Parse()
	expected := "name IN(?,?)"
	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}

	whereIn = WhereIn{
		Field:  "age",
		Values: []any{19, 20, 31},
	}
	stmt = whereIn.Parse()
	expected = "age IN(?,?,?)"
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
	expected := "name = ?"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}

	where = Where{
		Field: "age",
		Op:    OperatorGreaterThan,
		Value: 17,
	}

	stmt = where.Parse()
	expected = "age > ?"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}

func TestOrderClause(t *testing.T) {
	order := Order{
		OrderingFields: []OrderField{
			{
				Field:     "name",
				Direction: OrderDirectionDESC,
			},
		},
	}

	stmt := order.Parse()
	expected := "ORDER BY name desc"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}

func TestJoinClauseParsing(t *testing.T) {
	join := Join{
		Type:        LeftJoin,
		FirstTable:  "users",
		SecondTable: "orders",
		On: JoinON{
			LeftValue:  "id",
			Operator:   OperatorEqual,
			RightValue: "user_id",
		},
	}

	stmt := join.Parse()
	expected := "LEFT JOIN orders ON users.id = orders.user_id"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}
