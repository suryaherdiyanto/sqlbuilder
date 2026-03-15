package clause

import (
	"testing"

	"github.com/suryaherdiyanto/sqlbuilder/dialect"
)

func TestJoinClauseParsing(t *testing.T) {
	dialect := dialect.New("?", "`", "`")
	join := Join{
		Type:        LeftJoin,
		SecondTable: "orders",
		On: JoinON{
			Operator:   OperatorEqual,
			LeftField:  "users.id",
			RightField: "orders.user_id",
		},
	}

	stmt := join.Parse(dialect)
	expected := "LEFT JOIN `orders` ON `users`.`id` = `orders`.`user_id`"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}

func TestJoinClauseParsingPG(t *testing.T) {
	dialect := dialect.NewPostgres()
	join := Join{
		Type:        LeftJoin,
		SecondTable: "orders",
		On: JoinON{
			Operator:   OperatorEqual,
			LeftField:  "users.id",
			RightField: "orders.user_id",
		},
	}

	stmt := join.Parse(dialect)
	expected := "LEFT JOIN \"orders\" ON \"users\".\"id\" = \"orders\".\"user_id\""

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}

func TestStatementWithJoin(t *testing.T) {
	dialect := dialect.New("?", "`", "`")
	statement := Select{
		Table:   "users",
		Columns: []string{"users.id", "users.email", "orders.total"},
		Joins: []Join{
			{
				Type:        InnerJoin,
				SecondTable: "orders",
				On: JoinON{
					LeftField:  "users.id",
					Operator:   OperatorEqual,
					RightField: "orders.user_id",
				},
			},
		},
	}
	where := WhereStatements{
		Where: []Where{
			{
				Field: "orders.total",
				Op:    OperatorGreaterThan,
				Value: 10000,
			},
		},
	}

	stmt, _ := statement.Parse(dialect)
	stmt += where.Parse(dialect)

	expected := "SELECT `users`.`id`,`users`.`email`,`orders`.`total` FROM `users` INNER JOIN `orders` ON `users`.`id` = `orders`.`user_id` WHERE `orders`.`total` > ?"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}

func TestStatementMultipleJoin(t *testing.T) {
	dialect := dialect.New("?", "`", "`")
	statement := Select{
		Table:   "users",
		Columns: []string{"users.*", "user_orders.total", "products.name"},
		Joins: []Join{
			{
				Type:        InnerJoin,
				SecondTable: "user_orders",
				On: JoinON{
					LeftField:  "users.id",
					Operator:   OperatorEqual,
					RightField: "user_orders.user_id",
				},
			},
			{
				Type:        InnerJoin,
				SecondTable: "products",
				On: JoinON{
					LeftField:  "user_orders.product_id",
					Operator:   OperatorEqual,
					RightField: "products.id",
				},
			},
		},
	}
	where := WhereStatements{
		Where: []Where{
			{
				Field: "user_orders.total",
				Op:    OperatorGreaterThan,
				Value: 10000,
			},
		},
	}

	stmt, _ := statement.Parse(dialect)
	stmt += where.Parse(dialect)

	expected := "SELECT `users`.*,`user_orders`.`total`,`products`.`name` FROM `users` INNER JOIN `user_orders` ON `users`.`id` = `user_orders`.`user_id` INNER JOIN `products` ON `user_orders`.`product_id` = `products`.`id` WHERE `user_orders`.`total` > ?"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}

func TestStatementWithJoinPG(t *testing.T) {
	dialect := dialect.NewPostgres()
	statement := Select{
		Table:   "users",
		Columns: []string{"users.id", "users.email", "orders.total"},
		Joins: []Join{
			{
				Type:        InnerJoin,
				SecondTable: "orders",
				On: JoinON{
					LeftField:  "users.id",
					Operator:   OperatorEqual,
					RightField: "orders.user_id",
				},
			},
		},
	}
	where := WhereStatements{
		Where: []Where{
			{
				Field: "orders.total",
				Op:    OperatorGreaterThan,
				Value: 10000,
			},
		},
	}

	stmt, _ := statement.Parse(dialect)
	stmt += where.Parse(dialect)

	expected := "SELECT \"users\".\"id\",\"users\".\"email\",\"orders\".\"total\" FROM \"users\" INNER JOIN \"orders\" ON \"users\".\"id\" = \"orders\".\"user_id\" WHERE \"orders\".\"total\" > $1"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}
