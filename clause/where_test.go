package clause

import (
	"testing"

	"github.com/suryaherdiyanto/sqlbuilder/dialect"
)

func TestWhereParsing(t *testing.T) {
	dialect := dialect.New("?", "`", "`")
	where := Where{Field: "name", Op: OperatorEqual, Value: "Alice"}

	stmt := where.Parse(dialect)
	expected := "`name` = ?"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}

	where2 := Where{Field: "age", Op: OperatorGreaterThan, Value: 30}
	stmt = where2.Parse(dialect)
	expected = "`age` > ?"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}

func TestWhereParsingPG(t *testing.T) {
	dialect := dialect.NewPostgres()
	where := Where{Field: "name", Op: OperatorEqual, Value: "Alice"}

	stmt := where.Parse(dialect)
	expected := `"name" = $1`

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}

	where2 := Where{Field: "age", Op: OperatorGreaterThan, Value: 30}
	stmt = where2.Parse(dialect)
	expected = `"age" > $2`

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}

func TestWhereGroupParsing(t *testing.T) {
	dialect := dialect.New("?", "`", "`")
	whereGroup := WhereGroup{
		Conj: ConjuctionAnd,
		WhereStatements: WhereStatements{
			Where: []Where{
				{Field: "name", Op: OperatorEqual, Conj: ConjuctionAnd, Value: "Alice"},
				{Field: "age", Op: OperatorGreaterThan, Conj: ConjuctionAnd, Value: 30},
			},
		},
	}

	stmt := whereGroup.Parse(dialect)
	expected := "`name` = ? AND `age` > ?"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}

func TestWhereStatementWithGroupsParsing(t *testing.T) {
	dialect := dialect.New("?", "`", "`")
	statements := WhereStatements{
		Where: []Where{
			{Field: "status", Op: OperatorEqual, Conj: ConjuctionAnd, Value: "active"},
		},
	}
	whereGroup1 := WhereGroup{
		Conj: ConjuctionAnd,
		WhereStatements: WhereStatements{
			Where: []Where{
				{Field: "name", Op: OperatorEqual, Conj: ConjuctionAnd, Value: "Alice"},
				{Field: "age", Op: OperatorGreaterThan, Conj: ConjuctionAnd, Value: 30},
			},
		},
	}

	whereGroup2 := WhereGroup{
		Conj: ConjuctionAnd,
		WhereStatements: WhereStatements{
			Where: []Where{
				{Field: "city", Op: OperatorEqual, Conj: ConjuctionOr, Value: "New York"},
				{Field: "country", Op: OperatorEqual, Conj: ConjuctionOr, Value: "USA"},
			},
		},
	}

	where := Where{
		Conj: ConjuctionAnd,
		Groups: []WhereGroup{
			whereGroup1,
			whereGroup2,
		},
	}

	statements.Where = append(statements.Where, where)
	stmt := statements.ParseWhereStatements(dialect)
	expected := "`status` = ? AND (`name` = ? AND `age` > ?) AND (`city` = ? OR `country` = ?)"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}
