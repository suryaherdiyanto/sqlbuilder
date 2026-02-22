package clause

import (
	"testing"

	"github.com/suryaherdiyanto/sqlbuilder/dialect"
)

func TestJoinClauseParsing(t *testing.T) {
	dialect := dialect.New("?", "`", "`")
	join := Join{
		Type:        LeftJoin,
		FirstTable:  "users",
		SecondTable: "orders",
		On: JoinON{
			Operator:   OperatorEqual,
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

func TestJoinClauseParsingPG(t *testing.T) {
	dialect := dialect.NewPostgres()
	join := Join{
		Type:        LeftJoin,
		FirstTable:  "users",
		SecondTable: "orders",
		On: JoinON{
			Operator:   OperatorEqual,
			LeftField:  "id",
			RightField: "user_id",
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
				FirstTable:  "users",
				SecondTable: "orders",
				On: JoinON{
					LeftField:  "id",
					Operator:   OperatorEqual,
					RightField: "user_id",
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

func TestStatementWithJoinPG(t *testing.T) {
	dialect := dialect.NewPostgres()
	statement := Select{
		Table:   "users",
		Columns: []string{"users.id", "users.email", "orders.total"},
		Joins: []Join{
			{
				Type:        InnerJoin,
				FirstTable:  "users",
				SecondTable: "orders",
				On: JoinON{
					LeftField:  "id",
					Operator:   OperatorEqual,
					RightField: "user_id",
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
