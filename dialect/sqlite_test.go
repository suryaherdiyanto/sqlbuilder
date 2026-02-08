package dialect

import (
	"testing"

	"github.com/suryaherdiyanto/sqlbuilder/clause"
)

func TestWhereInParsing(t *testing.T) {
	dialect := NewSQLiteDialect()
	where := clause.WhereIn{
		Field:  "name",
		Values: []any{"Alice", "Bob"},
	}

	stmt := where.Parse(dialect)
	expected := "`name` IN(?,?)"
	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}

	wherein2 := clause.WhereIn{
		Field:  "age",
		Values: []any{25, 30, 35},
	}
	stmt = wherein2.Parse(dialect)
	expected = "`age` IN(?,?,?)"
	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}

func TestWhereParsing(t *testing.T) {
	dialect := NewSQLiteDialect()
	where := clause.Where{Field: "name", Op: clause.OperatorEqual, Value: "Alice"}

	stmt := where.Parse(dialect)
	expected := "`name` = ?"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}

	where2 := clause.Where{Field: "age", Op: clause.OperatorGreaterThan, Value: 30}
	stmt = where2.Parse(dialect)
	expected = "`age` > ?"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}

func TestOrderClause(t *testing.T) {
	dialect := NewSQLiteDialect()
	order := clause.Order{
		OrderingFields: []clause.OrderField{
			{Field: "name", Direction: clause.OrderDirectionDESC},
		},
	}

	stmt := order.Parse(dialect)
	expected := "ORDER BY name DESC"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}

func TestJoinClauseParsing(t *testing.T) {
	dialect := NewSQLiteDialect()
	join := clause.Join{
		Type:        clause.LeftJoin,
		FirstTable:  "users",
		SecondTable: "orders",
		On: clause.JoinON{
			Operator:   clause.OperatorEqual,
			LeftField:  "id",
			RightField: "user_id",
		},
	}

	stmt := join.Parse(dialect)
	expected := "LEFT JOIN `orders` ON `users`.`id` = `orders`.`user_id`"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}

func TestStatementMultipleWhere(t *testing.T) {
	dialect := NewSQLiteDialect()
	statement := clause.Select{
		Table:   "users",
		Columns: []string{"id", "email", "name"},
		WhereStatements: clause.WhereStatements{
			Where: []clause.Where{
				{
					Field: "email",
					Value: "johndoe@gmail.com",
					Op:    clause.OperatorEqual,
				},
				{
					Field: "access_role",
					Value: 3,
					Op:    clause.OperatorLessThan,
					Conj:  clause.ConjuctionAnd,
				},
			},
		},
	}

	stmt, _ := statement.Parse(dialect)
	expected := "SELECT `id`,`email`,`name` FROM `users` WHERE `email` = ? AND `access_role` < ?"
	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}

func TestStatementWhereConjuctionOr(t *testing.T) {
	dialect := NewSQLiteDialect()
	statement := clause.Select{
		Table:   "users",
		Columns: []string{"id", "email", "name"},
		WhereStatements: clause.WhereStatements{
			Where: []clause.Where{
				{
					Field: "email",
					Value: "johndoe@gmail.com",
					Op:    clause.OperatorEqual,
				},
				{
					Field: "role",
					Value: "admin",
					Op:    clause.OperatorLessThan,
					Conj:  clause.ConjuctionOr,
				},
			},
		},
	}

	stmt, _ := statement.Parse(dialect)
	expected := "SELECT `id`,`email`,`name` FROM `users` WHERE `email` = ? OR `role` < ?"
	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}

func TestSimpleSelectStatement(t *testing.T) {
	dialect := NewSQLiteDialect()
	statement := clause.Select{
		Table:   "users",
		Columns: []string{"id", "email", "name"},
		WhereStatements: clause.WhereStatements{
			Where: []clause.Where{
				{
					Field: "email",
					Value: "johndoe@gmail.com",
					Op:    clause.OperatorEqual,
				},
			},
		},
	}

	stmt, s := statement.Parse(dialect)
	expected := "SELECT `id`,`email`,`name` FROM `users` WHERE `email` = ?"
	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}

	arguments := s.GetArguments()

	if len(arguments) != 1 {
		t.Errorf("Expected 1 argument, but got: %d", len(arguments))
	}

}

func TestSimpleStatementWithLimitAndOffset(t *testing.T) {
	dialect := NewSQLiteDialect()
	statement := clause.Select{
		Table:   "users",
		Columns: []string{"id", "email", "name"},
		Limit:   clause.Limit{Count: 10},
		Offset:  clause.Offset{Count: 5},
	}

	stmt, _ := statement.Parse(dialect)
	expected := "SELECT `id`,`email`,`name` FROM `users` LIMIT ? OFFSET ?"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}

func TestStatementWhereIn(t *testing.T) {
	dialect := NewSQLiteDialect()
	statement := clause.Select{
		Table:   "users",
		Columns: []string{"id", "email", "name"},
		WhereStatements: clause.WhereStatements{
			WhereIn: []clause.WhereIn{
				{
					Field:  "email",
					Values: []any{"johndoe@gmail.com", "test@example.com"},
				},
			},
		},
	}

	stmt, _ := statement.Parse(dialect)
	expected := "SELECT `id`,`email`,`name` FROM `users` WHERE `email` IN(?,?)"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}

func TestStatementWhereInWithConjuction(t *testing.T) {
	dialect := NewSQLiteDialect()
	statement := clause.Select{
		Table:   "users",
		Columns: []string{"id", "email", "name"},
		WhereStatements: clause.WhereStatements{
			Where: []clause.Where{
				{
					Field: "name",
					Op:    clause.OperatorEqual,
					Value: "John",
				},
			},
			WhereIn: []clause.WhereIn{
				{
					Field:  "email",
					Values: []any{"johndoe@gmail.com", "test@example.com"},
					Conj:   clause.ConjuctionAnd,
				},
			},
		},
	}

	stmt, _ := statement.Parse(dialect)
	expected := "SELECT `id`,`email`,`name` FROM `users` WHERE `name` = ? AND `email` IN(?,?)"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}

func TestStatementWhereNotIn(t *testing.T) {
	dialect := NewSQLiteDialect()
	statement := clause.Select{
		Table:   "users",
		Columns: []string{"id", "email", "name"},
		WhereStatements: clause.WhereStatements{
			WhereNotIn: []clause.WhereNotIn{
				{
					Field:  "email",
					Values: []any{"johndoe@gmail.com", "test@example.com"},
				},
			},
		},
	}

	stmt, _ := statement.Parse(dialect)
	expected := "SELECT `id`,`email`,`name` FROM `users` WHERE `email` NOT IN(?,?)"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}

func TestStatementWithJoin(t *testing.T) {
	dialect := NewSQLiteDialect()
	statement := clause.Select{
		Table:   "users",
		Columns: []string{"users.id", "users.email", "orders.total"},
		Joins: []clause.Join{
			{
				Type:        clause.InnerJoin,
				FirstTable:  "users",
				SecondTable: "orders",
				On: clause.JoinON{
					LeftField:  "id",
					Operator:   clause.OperatorEqual,
					RightField: "user_id",
				},
			},
		},
		WhereStatements: clause.WhereStatements{
			Where: []clause.Where{
				{
					Field: "orders.total",
					Op:    clause.OperatorGreaterThan,
					Value: 10000,
				},
			},
		},
	}

	stmt, _ := statement.Parse(dialect)
	expected := "SELECT `users`.`id`,`users`.`email`,`orders`.`total` FROM `users` INNER JOIN `orders` ON `users`.`id` = `orders`.`user_id` WHERE `orders`.`total` > ?"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}

func TestWithSubStatementWhere(t *testing.T) {
	dialect := NewSQLiteDialect()
	statement := clause.Select{
		Table:   "users",
		Columns: []string{"*"},
		WhereStatements: clause.WhereStatements{
			Where: []clause.Where{
				{
					Field: "roles_id",
					Op:    clause.OperatorEqual,
					SubStatement: clause.Select{
						Table:   "roles",
						Columns: []string{"id"},
						WhereStatements: clause.WhereStatements{
							Where: []clause.Where{
								{
									Field: "roles.id",
									Op:    clause.OperatorEqual,
									Value: 3,
								},
							},
						},
					},
				},
			},
		},
	}

	stmt, _ := statement.Parse(dialect)
	expected := "SELECT * FROM `users` WHERE `roles_id` = (SELECT `id` FROM `roles` WHERE `roles`.`id` = ?)"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}

func TestWithGroupByStatement(t *testing.T) {
	dialect := NewSQLiteDialect()
	statement := clause.Select{
		Table:   "users",
		Columns: []string{"role", "COUNT(*) as total"},
		GroupBy: clause.GroupBy{
			Fields: []string{"role"},
		},
	}

	stmt, _ := statement.Parse(dialect)
	expected := "SELECT `role`,COUNT(*) as total FROM `users` GROUP BY `role`"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}

func TestWhereInParsingWithSubquery(t *testing.T) {
	dialect := NewSQLiteDialect()
	statement := clause.Select{
		Table:   "users",
		Columns: []string{"*"},
		WhereStatements: clause.WhereStatements{
			WhereIn: []clause.WhereIn{
				{
					Field: "id",
					SubStatement: clause.Select{
						Table:   "orders",
						Columns: []string{"user_id"},
						WhereStatements: clause.WhereStatements{
							Where: []clause.Where{
								{
									Field: "total",
									Op:    clause.OperatorGreaterThan,
									Value: 100,
								},
							},
						},
					},
				},
			},
		},
	}

	stmt, _ := statement.Parse(dialect)
	expected := "SELECT * FROM `users` WHERE `id` IN (SELECT `user_id` FROM `orders` WHERE `total` > ?)"
	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}

func TestWhereNotInParsingWithSubquery(t *testing.T) {
	dialect := NewSQLiteDialect()
	statement := clause.Select{
		Table:   "users",
		Columns: []string{"*"},
		WhereStatements: clause.WhereStatements{
			WhereNotIn: []clause.WhereNotIn{
				{
					Field: "id",
					SubStatement: clause.Select{
						Table:   "banned_users",
						Columns: []string{"user_id"},
					},
				},
			},
		},
	}

	stmt, _ := statement.Parse(dialect)
	expected := "SELECT * FROM `users` WHERE `id` NOT IN (SELECT `user_id` FROM `banned_users`)"
	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}

func TestInsertStatement(t *testing.T) {
	type User struct {
		Name  string `db:"name"`
		Email string `db:"email"`
		Age   int    `db:"age"`
	}
	dialect := NewSQLiteDialect()
	statement := clause.Insert{
		Table: "users",
		Rows: []map[string]any{
			{
				"name":  "John Doe",
				"email": "johndoe@example.com",
				"age":   30,
			},
		},
	}

	stmt := dialect.ParseInsert(statement)
	expected := "INSERT INTO users(`age`,`email`,`name`) VALUES(?,?,?)"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}

func TestInsertStatementMultiple(t *testing.T) {

	dialect := NewSQLiteDialect()
	statement := clause.Insert{
		Table: "users",
		Rows: []map[string]any{
			{
				"name":  "John Doe",
				"email": "johndoe@example.com",
				"age":   30,
			},
			{
				"name":  "Foo barr",
				"email": "foobarr@example.com",
				"age":   21,
			},
		},
	}

	stmt := statement.Parse(dialect)
	expected := "INSERT INTO users(`age`,`email`,`name`) VALUES(?,?,?),(?,?,?)"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}

func TestUpdateStatement(t *testing.T) {
	dialect := NewSQLiteDialect()
	statement := clause.Update{
		Table: "users",
		Rows: map[string]any{
			"name": "test",
			"age":  25,
		},
		WhereStatements: clause.WhereStatements{
			Where: []clause.Where{
				{
					Field: "id",
					Op:    clause.OperatorEqual,
					Value: 1,
				},
			},
		},
	}

	stmt := dialect.ParseUpdate(statement)
	expected := "UPDATE users SET `age` = ?, `name` = ? WHERE `id` = ?"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}

func TestDeleteStatement(t *testing.T) {
	dialect := NewSQLiteDialect()
	statement := clause.Delete{
		Table: "users",
		WhereStatements: clause.WhereStatements{
			Where: []clause.Where{
				{
					Field: "id",
					Op:    clause.OperatorEqual,
					Value: 1,
				},
			},
		},
	}

	stmt := statement.Parse(dialect)
	expected := "DELETE FROM `users` WHERE `id` = ?"
	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}
