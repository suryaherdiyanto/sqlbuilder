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

func TestStatementMultipleWhere(t *testing.T) {
	statement := SelectStatement{
		Table:   "users",
		Columns: []string{"id", "email", "name"},
		WhereStatements: []Where{
			{
				Field: "email",
				Value: "johndoe@gmail.com",
				Op:    OperatorEqual,
			},
			{
				Field: "access_role",
				Value: 3,
				Op:    OperatorLessThan,
				Conj:  ConjuctionAnd,
			},
		},
	}

	stmt := statement.Parse()
	expected := "SELECT id,email,name FROM users WHERE email = ? AND access_role < ?"
	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}

func TestStatementWhereConjuctionOr(t *testing.T) {
	statement := SelectStatement{
		Table:   "users",
		Columns: []string{"id", "email", "name"},
		WhereStatements: []Where{
			{
				Field: "email",
				Value: "johndoe@gmail.com",
				Op:    OperatorEqual,
			},
			{
				Field: "role",
				Value: "admin",
				Op:    OperatorLessThan,
				Conj:  ConjuctionOr,
			},
		},
	}

	stmt := statement.Parse()
	expected := "SELECT id,email,name FROM users WHERE email = ? OR role < ?"
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
	expected := "SELECT id,email,name FROM users WHERE email = ?"
	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}

}

func TestJoinClauseParsing(t *testing.T) {
	join := Join{
		Type:       LeftJoin,
		OtherTable: "orders",
		On: JoinON{
			LeftTable:  "users",
			LeftValue:  "id",
			RightTable: "orders",
			RightValue: "user_id",
		},
	}

	stmt := join.Parse()
	expected := "LEFT JOIN orders ON users.id = orders.user_id"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}

func TestSimpleStatementWithLimitAndOffset(t *testing.T) {
	statement := SelectStatement{
		Table:   "users",
		Columns: []string{"id", "email", "name"},
		Limit:   10,
		Offset:  5,
	}

	stmt := statement.Parse()
	expected := "SELECT id,email,name FROM users LIMIT 10 OFFSET 5"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}

func TestStatementWhereIn(t *testing.T) {
	statement := SelectStatement{
		Table:   "users",
		Columns: []string{"id", "email", "name"},
		WhereInStatements: []WhereIn{
			{
				Field:  "email",
				Values: []any{"johndoe@gmail.com", "test@example.com"},
			},
		},
	}

	stmt := statement.Parse()
	expected := "SELECT id,email,name FROM users WHERE email IN(?,?)"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}

func TestStatementWhereInWithConjuction(t *testing.T) {
	statement := SelectStatement{
		Table:   "users",
		Columns: []string{"id", "email", "name"},
		WhereInStatements: []WhereIn{
			{
				Field:  "email",
				Values: []any{"johndoe@gmail.com", "test@example.com"},
				Conj:   ConjuctionAnd,
			},
		},
		WhereStatements: []Where{
			{
				Field: "name",
				Op:    OperatorEqual,
				Value: "John",
			},
		},
	}

	stmt := statement.Parse()
	expected := "SELECT id,email,name FROM users WHERE name = ? AND email IN(?,?)"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}

func TestStatementWhereNotIn(t *testing.T) {
	statement := SelectStatement{
		Table:   "users",
		Columns: []string{"id", "email", "name"},
		WhereNotInStatements: []WhereNotIn{
			{
				Field:  "email",
				Values: []any{"johndoe@gmail.com", "test@example.com"},
			},
		},
	}

	stmt := statement.Parse()
	expected := "SELECT id,email,name FROM users WHERE email NOT IN(?,?)"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}

func TestStatementWithJoin(t *testing.T) {
	statement := SelectStatement{
		Table:   "users",
		Columns: []string{"users.id", "users.email", "orders.total"},
		JoinStatements: []Join{
			{
				Type:       InnerJoin,
				OtherTable: "orders",
				On: JoinON{
					LeftTable:  "users",
					LeftValue:  "id",
					RightTable: "orders",
					RightValue: "user_id",
				},
			},
		},
		WhereStatements: []Where{
			{
				Field: "orders.total",
				Op:    OperatorGreaterThan,
				Value: 10000,
			},
		},
	}

	stmt := statement.Parse()
	expected := "SELECT users.id,users.email,orders.total FROM users INNER JOIN orders ON users.id = orders.user_id WHERE orders.total > ?"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}

func TestWithSubStatementWhere(t *testing.T) {
	statement := SelectStatement{
		Table:   "users",
		Columns: []string{"*"},
		WhereStatements: []Where{
			{
				Field: "roles_id",
				Op:    OperatorEqual,
				SubStatement: SelectStatement{
					Table:   "roles",
					Columns: []string{"id"},
					WhereStatements: []Where{
						{
							Field: "roles.id",
							Op:    OperatorEqual,
							Value: 3,
						},
					},
				},
			},
		},
	}

	stmt := statement.Parse()
	expected := "SELECT * FROM users WHERE roles_id = (SELECT id FROM roles WHERE roles.id = ?)"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}
