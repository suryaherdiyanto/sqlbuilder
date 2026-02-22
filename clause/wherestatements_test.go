package clause

import (
	"testing"

	"github.com/suryaherdiyanto/sqlbuilder/dialect"
)

func TestStatementMultipleWhere(t *testing.T) {
	dialect := dialect.New("?", "`", "`")
	statement := Select{
		Columns: []string{"*"},
		Table:   "users",
	}
	where := WhereStatements{
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
	}

	stmt, _ := statement.Parse(dialect)
	stmt += where.Parse(dialect)
	expected := "SELECT * FROM `users` WHERE `email` = ? AND `access_role` < ?"
	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}

func TestStatementWhereConjuctionOr(t *testing.T) {
	dialect := dialect.New("?", "`", "`")
	statement := Select{
		Columns: []string{"*"},
		Table:   "users",
	}
	where := WhereStatements{
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
	}

	stmt, _ := statement.Parse(dialect)
	stmt += where.Parse(dialect)
	expected := "SELECT * FROM `users` WHERE `email` = ? OR `role` < ?"
	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}

func TestStatementWhereIn(t *testing.T) {
	dialect := dialect.New("?", "`", "`")
	statement := Select{
		Table:   "users",
		Columns: []string{"id", "email", "name"},
	}
	where := WhereStatements{
		WhereIn: []WhereIn{
			{
				Field:  "email",
				Values: []any{"johndoe@gmail.com", "test@example.com"},
			},
		},
	}

	stmt, _ := statement.Parse(dialect)
	stmt += where.Parse(dialect)
	expected := "SELECT `id`,`email`,`name` FROM `users` WHERE `email` IN(?,?)"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}

func TestStatementWhereInWithConjuction(t *testing.T) {
	dialect := dialect.New("?", "`", "`")
	statement := Select{
		Table:   "users",
		Columns: []string{"id", "email", "name"},
	}
	where := WhereStatements{
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
	}

	stmt, _ := statement.Parse(dialect)
	stmt += where.Parse(dialect)
	expected := "SELECT `id`,`email`,`name` FROM `users` WHERE `name` = ? AND `email` IN(?,?)"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}

func TestStatementWhereNotIn(t *testing.T) {
	dialect := dialect.New("?", "`", "`")
	statement := Select{
		Table:   "users",
		Columns: []string{"id", "email", "name"},
	}
	where := WhereStatements{
		WhereNotIn: []WhereNotIn{
			{
				Field:  "email",
				Values: []any{"johndoe@gmail.com", "test@example.com"},
			},
		},
	}

	stmt, _ := statement.Parse(dialect)
	stmt += where.Parse(dialect)
	expected := "SELECT `id`,`email`,`name` FROM `users` WHERE `email` NOT IN(?,?)"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}

func TestWithSubStatementWhere(t *testing.T) {
	dialect := dialect.New("?", "`", "`")
	statement := Select{
		Table:   "users",
		Columns: []string{"*"},
	}
	where := WhereStatements{
		Where: []Where{
			{
				Field: "roles_id",
				Op:    OperatorEqual,
				SubStatement: SubStatement{
					Select: Select{
						Table:   "roles",
						Columns: []string{"id"},
					},
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
	}

	stmt, _ := statement.Parse(dialect)
	stmt += where.Parse(dialect)
	expected := "SELECT * FROM `users` WHERE `roles_id` = (SELECT `id` FROM `roles` WHERE `roles`.`id` = ?)"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}

func TestWhereInParsingWithSubquery(t *testing.T) {
	dialect := dialect.New("?", "`", "`")
	statement := Select{
		Table:   "users",
		Columns: []string{"*"},
	}
	where := WhereStatements{
		WhereIn: []WhereIn{
			{
				Field: "id",
				SubStatement: SubStatement{
					Select: Select{
						Table:   "orders",
						Columns: []string{"user_id"},
					},
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
	}

	stmt, _ := statement.Parse(dialect)
	stmt += where.Parse(dialect)
	expected := "SELECT * FROM `users` WHERE `id` IN (SELECT `user_id` FROM `orders` WHERE `total` > ?)"
	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}

func TestWhereNotInParsingWithSubquery(t *testing.T) {
	dialect := dialect.New("?", "`", "`")
	statement := Select{
		Table:   "users",
		Columns: []string{"*"},
	}
	where := WhereStatements{
		WhereNotIn: []WhereNotIn{
			{
				Field: "id",
				SubStatement: Select{
					Table:   "banned_users",
					Columns: []string{"user_id"},
				},
			},
		},
	}

	stmt, _ := statement.Parse(dialect)
	stmt += where.Parse(dialect)

	expected := "SELECT * FROM `users` WHERE `id` NOT IN (SELECT `user_id` FROM `banned_users`)"
	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}
