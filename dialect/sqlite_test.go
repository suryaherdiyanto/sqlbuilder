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
	expected := "name IN(?,?)"
	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}

	wherein2 := clause.WhereIn{
		Field:  "age",
		Values: []any{25, 30, 35},
	}
	stmt = wherein2.Parse(dialect)
	expected = "age IN(?,?,?)"
	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}

func TestWhereParsing(t *testing.T) {
	dialect := NewSQLiteDialect()
	where := clause.Where{Field: "name", Op: clause.OperatorEqual, Value: "Alice"}

	stmt := where.Parse(dialect)
	expected := "name = ?"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}

	where2 := clause.Where{Field: "age", Op: clause.OperatorGreaterThan, Value: 30}
	stmt = where2.Parse(dialect)
	expected = "age > ?"

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
