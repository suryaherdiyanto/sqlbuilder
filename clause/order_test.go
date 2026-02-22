package clause

import (
	"testing"

	"github.com/suryaherdiyanto/sqlbuilder/dialect"
)

func TestOrderClause(t *testing.T) {
	dialect := dialect.New("?", "`", "`")
	order := Order{
		OrderingFields: []OrderField{
			{Field: "name", Direction: OrderDirectionDESC},
		},
	}

	stmt := order.Parse(dialect)
	expected := "ORDER BY name DESC"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}

func TestOrderClausePG(t *testing.T) {
	dialect := dialect.NewPostgres()
	order := Order{
		OrderingFields: []OrderField{
			{Field: "name", Direction: OrderDirectionDESC},
		},
	}

	stmt := order.Parse(dialect)
	expected := "ORDER BY name DESC"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}
