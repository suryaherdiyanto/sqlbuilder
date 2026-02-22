package clause

import (
	"testing"

	"github.com/suryaherdiyanto/sqlbuilder/dialect"
)

func TestInsertStatement(t *testing.T) {
	type User struct {
		Name  string `db:"name"`
		Email string `db:"email"`
		Age   int    `db:"age"`
	}
	dialect := dialect.New("?", "`", "`")
	statement := Insert{
		Table: "users",
		Rows: []map[string]any{
			{
				"name":  "John Doe",
				"email": "johndoe@example.com",
				"age":   30,
			},
		},
	}

	stmt := statement.Parse(dialect)
	expected := "INSERT INTO users(`age`,`email`,`name`) VALUES(?,?,?)"

	if stmt != expected {
		t.Errorf("Expected: %s, but got: %s", expected, stmt)
	}
}

func TestInsertStatementMultiple(t *testing.T) {
	dialect := dialect.New("?", "`", "`")
	statement := Insert{
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
