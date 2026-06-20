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
	where := Where{
		Field: "id",
		Op:    OperatorEqual,
		Value: 1,
	}

	stmt, _ := statement.Parse(dialect, 2)
	stmt += " WHERE " + where.Parse(dialect)
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
	where := Where{
		Field: "id",
		Op:    OperatorEqual,
		Value: 1,
	}

	w := where.Parse(dialect)
	stmt, _ := statement.Parse(dialect, 2)
	stmt = stmt + " WHERE " + w
	expected := "UPDATE users SET \"age\" = $2, \"name\" = $3 WHERE \"id\" = $1"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}
