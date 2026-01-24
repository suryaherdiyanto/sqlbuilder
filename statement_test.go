package sqlbuilder

import (
	"testing"
)

func TestStatementMultipleWhere(t *testing.T) {
	statement := SelectStatement{
		Table:   "users",
		Columns: []string{"id", "email", "name"},
		WhereStatements: WhereStatements{
			Where: []Where{
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
		WhereStatements: WhereStatements{
			Where: []Where{
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
		},
	}

	stmt := statement.Parse()
	expected := "SELECT id,email,name FROM users WHERE email = ? OR role < ?"
	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}

func TestSimpleSelectStatement(t *testing.T) {
	statement := SelectStatement{
		Table:   "users",
		Columns: []string{"id", "email", "name"},
		WhereStatements: WhereStatements{
			Where: []Where{
				{
					Field: "email",
					Value: "johndoe@gmail.com",
					Op:    OperatorEqual,
				},
			},
		},
	}

	stmt := statement.Parse()
	expected := "SELECT id,email,name FROM users WHERE email = ?"
	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}

	arguments := statement.GetArguments()

	if len(arguments) != 1 {
		t.Errorf("Expected 1 argument, but got: %d", len(arguments))
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
	expected := "SELECT id,email,name FROM users LIMIT ? OFFSET ?"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}

func TestStatementWhereIn(t *testing.T) {
	statement := SelectStatement{
		Table:   "users",
		Columns: []string{"id", "email", "name"},
		WhereStatements: WhereStatements{
			WhereIn: []WhereIn{
				{
					Field:  "email",
					Values: []any{"johndoe@gmail.com", "test@example.com"},
				},
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
		WhereStatements: WhereStatements{
			Where: []Where{
				{
					Field: "name",
					Op:    OperatorEqual,
					Value: "John",
				},
			},
			WhereIn: []WhereIn{
				{
					Field:  "email",
					Values: []any{"johndoe@gmail.com", "test@example.com"},
					Conj:   ConjuctionAnd,
				},
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
		WhereStatements: WhereStatements{
			WhereNotIn: []WhereNotIn{
				{
					Field:  "email",
					Values: []any{"johndoe@gmail.com", "test@example.com"},
				},
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
				Type:        InnerJoin,
				FirstTable:  "users",
				SecondTable: "orders",
				On: JoinON{
					LeftValue:  "id",
					Operator:   OperatorEqual,
					RightValue: "user_id",
				},
			},
		},
		WhereStatements: WhereStatements{
			Where: []Where{
				{
					Field: "orders.total",
					Op:    OperatorGreaterThan,
					Value: 10000,
				},
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
		WhereStatements: WhereStatements{
			Where: []Where{
				{
					Field: "roles_id",
					Op:    OperatorEqual,
					SubStatement: SelectStatement{
						Table:   "roles",
						Columns: []string{"id"},
						WhereStatements: WhereStatements{
							Where: []Where{
								{
									Field: "roles.id",
									Op:    OperatorEqual,
									Value: 3,
								},
							},
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

func TestWithGroupByStatement(t *testing.T) {
	statement := SelectStatement{
		Table:   "users",
		Columns: []string{"role", "COUNT(*) as total"},
		GroupByStatement: GroupBy{
			Fields: []string{"role"},
		},
	}

	stmt := statement.Parse()
	expected := "SELECT role,COUNT(*) as total FROM users GROUP BY role"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}

func TestWhereInParsingWithSubquery(t *testing.T) {
	statement := SelectStatement{
		Table:   "users",
		Columns: []string{"*"},
		WhereStatements: WhereStatements{
			WhereIn: []WhereIn{
				{
					Field: "id",
					SubStatement: SelectStatement{
						Table:   "orders",
						Columns: []string{"user_id"},
						WhereStatements: WhereStatements{
							Where: []Where{
								{
									Field: "total",
									Op:    OperatorGreaterThan,
									Value: 100,
								},
							},
						},
					},
				},
			},
		},
	}

	stmt := statement.Parse()
	expected := "SELECT * FROM users WHERE id IN (SELECT user_id FROM orders WHERE total > ?)"
	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}

func TestWhereNotInParsingWithSubquery(t *testing.T) {
	statement := SelectStatement{
		Table:   "users",
		Columns: []string{"*"},
		WhereStatements: WhereStatements{
			WhereNotIn: []WhereNotIn{
				{
					Field: "id",
					SubStatement: SelectStatement{
						Table:   "banned_users",
						Columns: []string{"user_id"},
					},
				},
			},
		},
	}

	stmt := statement.Parse()
	expected := "SELECT * FROM users WHERE id NOT IN (SELECT user_id FROM banned_users)"
	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}

func TestInsertStatement(t *testing.T) {
	statement := InsertStatement{
		Table: "users",
		Rows: []map[string]any{
			{
				"name":  "John Doe",
				"email": "johndoe@example.com",
				"age":   30,
			},
		},
	}

	stmt := statement.Parse()
	expected := "INSERT INTO users(`name`,`email`,`age`) VALUES(?,?,?)"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}

func TestUpdateStatement(t *testing.T) {
	statement := UpdateStatement{
		Table: "users",
		Rows: map[string]any{
			"name": "test",
			"age":  25,
		},
		WhereStatements: WhereStatements{
			Where: []Where{
				{
					Field: "id",
					Op:    OperatorEqual,
					Value: 1,
				},
			},
		},
	}

	stmt := statement.Parse()
	expected := "UPDATE users SET `name` = ?, `age` = ? WHERE id = ?"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}

func TestDeleteStatement(t *testing.T) {
	statement := DeleteStatement{
		Table: "users",
		WhereStatements: WhereStatements{
			Where: []Where{
				{
					Field: "id",
					Op:    OperatorEqual,
					Value: 1,
				},
			},
		},
	}

	stmt := statement.Parse()
	expected := "DELETE FROM users WHERE id = ?"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}
