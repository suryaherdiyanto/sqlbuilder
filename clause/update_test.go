package clause

import (
	"testing"

	"github.com/suryaherdiyanto/sqlbuilder/dialect"
)

func TestUpdateStatement(t *testing.T) {
	dialect := dialect.New("?", "`", "`")
	statement := Update{
		Table: "users",
		Rows: map[string]any{
			"name": "test",
			"age":  25,
		},
	}
	where := WhereStatements{
		Where: []Where{
			{
				Field: "id",
				Op:    OperatorEqual,
				Value: 1,
			},
		},
	}

	stmt, _ := statement.Parse(dialect)
	stmt += where.Parse(dialect)
	expected := "UPDATE users SET `age` = ?, `name` = ? WHERE `id` = ?"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}

func TestUpdateStatementPG(t *testing.T) {
	dialect := dialect.NewPostgres()
	statement := Update{
		Table: "users",
		Rows: map[string]any{
			"name": "test",
			"age":  25,
		},
	}
	where := WhereStatements{
		Where: []Where{
			{
				Field: "id",
				Op:    OperatorEqual,
				Value: 1,
			},
		},
	}

	stmt, _ := statement.Parse(dialect)
	stmt += where.Parse(dialect)
	expected := "UPDATE users SET \"age\" = $1, \"name\" = $2 WHERE \"id\" = $3"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}
