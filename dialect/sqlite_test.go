package dialect

import (
	"testing"

	"github.com/suryaherdiyanto/sqlbuilder/clause"
)

func TestWhereInParsing(t *testing.T) {
	dialect := NewSQLiteDialect()
	dialect.NewWhereIn("name", []any{"Alice", "Bob"}, clause.ConjuctionAnd)

	stmt := dialect.ParseWhereIn()
	expected := "name IN(?,?)"
	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}

	dialect = NewSQLiteDialect()
	dialect.NewWhereIn("age", []any{19, 20, 31}, clause.ConjuctionAnd)
	stmt = dialect.ParseWhereIn()
	expected = "age IN(?,?,?)"
	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}

func TestWhereParsing(t *testing.T) {
	dialect := NewSQLiteDialect()
	dialect.NewWhere("name", clause.OperatorEqual, "John Doe")

	stmt := dialect.ParseWhere()
	expected := "name = ?"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}

	dialect = NewSQLiteDialect()
	dialect.NewWhere("age", clause.OperatorGreaterThan, 17)

	stmt = dialect.ParseWhere()
	expected = "age > ?"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}

func TestOrderClause(t *testing.T) {
	dialect := NewSQLiteDialect()
	dialect.NewOrder([]clause.OrderField{
		{Field: "name", Direction: clause.OrderDirectionDESC},
	})

	stmt := dialect.ParseOrder()
	expected := "ORDER BY name DESC"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}

func TestJoinClauseParsing(t *testing.T) {
	dialect := NewSQLiteDialect()
	dialect.NewJoin(clause.LeftJoin, "users", "orders", clause.JoinON{
		LeftValue:  "id",
		Operator:   clause.OperatorEqual,
		RightValue: "user_id",
	})

	stmt := dialect.ParseJoin()
	expected := "LEFT JOIN orders ON users.id = orders.user_id"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}
