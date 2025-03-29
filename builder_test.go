package sqlbuilder

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

var db, _ = sql.Open("sqlite3", ":memory:")

type User struct {
	Id       int    `db:"id"`
	Username string `db:"username"`
	Email    string `db:"email"`
	Age      int    `db:"age"`
}

func setupSuite(tb testing.TB) func(tb testing.TB) {
	db.Exec(`
		CREATE TABLE users(
			id integer primary key,
			username TEXT,
			email TEXT,
			age integer
		)
	`)

	db.Exec(`
		INSERT INTO users values(null, 'johndoe', 'johndoe@example.com', 35);
	`)
	db.Exec(`
		INSERT INTO users values(null, 'daniel', 'daniel@example.com', 32);
	`)
	db.Exec(`
		INSERT INTO users values(null, 'samuel', 'samuel@example.com', 28);
	`)
	db.Exec(`
		INSERT INTO users values(null, 'dirt', 'dirt@example.com', 20);
	`)
	db.Exec(`
		INSERT INTO users values(null, 'chris', 'chris@example.com', 25);
	`)

	return func(tb testing.TB) {
		db.Exec("DROP table users")
	}
}

func TestBuilder(t *testing.T) {
	builder := NewSelect("sqlite", db)

	builder.Table("users", "*")

	if builder.GetSql() != "SELECT * FROM users" {
		t.Errorf("Unexpected SQL result, got: %s", builder.GetSql())
	}
}

func TestWithWhere(t *testing.T) {
	builder := NewSelect("sqlite", db)

	builder.
		Table("users", "*").
		Where("email = ?", "johndoe@gmail.com")

	if builder.GetSql() != "SELECT * FROM users WHERE email = ?" {
		t.Errorf("Unexpected SQL result, got: %s", builder.GetSql())
	}
}

func TestWithMultipleWhere(t *testing.T) {
	builder := NewSelect("sqlite", db)

	builder.
		Table("users", "*").
		Where("email = ? AND access_role < ?", "johndoe@gmail.com", 3)

	if builder.GetSql() != "SELECT * FROM users WHERE email = ? AND access_role < ?" {
		t.Errorf("Unexpected SQL result, got: %s", builder.GetSql())
	}
}

func TestWhereIn(t *testing.T) {
	builder := NewSelect("sqlite", db)

	builder.Table("users", "*").Where("email IN(?, ?)", "johndoe@example.com", "amal@example.com")

	if builder.GetSql() != "SELECT * FROM users WHERE email IN(?, ?)" {
		t.Errorf("Unexpected SQL result, got: %s", builder.GetSql())
	}
}

func TestWhereBetween(t *testing.T) {
	builder := NewSelect("sqlite", db)

	builder.
		Table("users", "*").
		Where("age BETWEEN ? AND ?", 5, 10)

	if builder.GetSql() != "SELECT * FROM users WHERE age BETWEEN ? AND ?" {
		t.Errorf("Unexpected SQL result, got: %s", builder.GetSql())
	}

	builder = NewSelect("sqlite", db)
	builder.
		Table("users", "*").
		Where("dob BETWEEN ? AND ?", "1995-02-01", "2000-01-01")

	if builder.GetSql() != "SELECT * FROM users WHERE dob BETWEEN ? AND ?" {
		t.Errorf("Unexpected SQL result, got: %s", builder.GetSql())
	}
}

func TestWhereOr(t *testing.T) {
	builder := NewSelect("sqlite", db)

	builder.
		Table("users", "*").
		Where("age >= ? OR email = ?", 18, "johndoe@example.com")

	if builder.GetSql() != "SELECT * FROM users WHERE age >= ? OR email = ?" {
		t.Errorf("Unexpected SQL result, got: %s", builder.GetSql())
	}
}

func TestExecute(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	user := new(User)
	builder := NewSelect("sqlite", db)

	err := builder.Table("users", "*").Scan(user, context.Background())

	if err != nil {
		t.Error(err)
	}

	if user.Email != "johndoe@example.com" {
		t.Errorf("Expected johndoe@example.com, but got: %s", user.Email)
	}

}

func TestExecuteWhere(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	user := new(User)
	builder := NewSelect("sqlite", db)

	err := builder.Table("users", "*").Where("email = ?", "daniel@example.com").Limit(1).Scan(user, context.Background())

	if err != nil {
		t.Error(err)
	}

	if user.Email != "daniel@example.com" {
		t.Errorf("Expected daniel@example.com, but got: %s", user.Email)
	}
}

func TestWhereAnd(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	var users []User
	builder := NewSelect("sqlite", db)

	err := builder.
		Table("users", "*").
		Where("age < ? AND email like ?", 30, "%@example.com").
		ScanAll(&users, context.Background())

	if err != nil {
		t.Error(err)
	}

	if len(users) != 3 {
		t.Errorf("Expected return %d of users, but got %d", 3, len(users))
	}

}
