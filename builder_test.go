package sqlbuilder

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

var db, _ = sql.Open("sqlite3", ":memory:")
var builder *SQLBuilder

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
	builder = New("sqlite", db)
	builder.NewSelect()

	builder.Table("users", "*")

	if builder.GetSql() != "SELECT * FROM users" {
		t.Errorf("Unexpected SQL result, got: %s", builder.GetSql())
	}
}

func TestWithWhere(t *testing.T) {
	builder = New("sqlite", db)
	builder.NewSelect()

	builder.
		Table("users", "*").
		Where("email = ?", "johndoe@gmail.com")

	if builder.GetSql() != "SELECT * FROM users WHERE email = ?" {
		t.Errorf("Unexpected SQL result, got: %s", builder.GetSql())
	}
}

func TestWithMultipleWhere(t *testing.T) {
	builder = New("sqlite", db)
	builder.NewSelect()

	builder.
		Table("users", "*").
		Where("email = ? AND access_role < ?", "johndoe@gmail.com", 3)

	if builder.GetSql() != "SELECT * FROM users WHERE email = ? AND access_role < ?" {
		t.Errorf("Unexpected SQL result, got: %s", builder.GetSql())
	}
}

func TestWhereIn(t *testing.T) {
	builder = New("sqlite", db)
	builder.NewSelect()

	builder.Table("users", "*").Where("email IN(?, ?)", "johndoe@example.com", "amal@example.com")

	if builder.GetSql() != "SELECT * FROM users WHERE email IN(?, ?)" {
		t.Errorf("Unexpected SQL result, got: %s", builder.GetSql())
	}
}

func TestWhereBetween(t *testing.T) {
	builder = New("sqlite", db)
	builder.NewSelect()

	builder.
		Table("users", "*").
		Where("age BETWEEN ? AND ?", 5, 10)

	if builder.GetSql() != "SELECT * FROM users WHERE age BETWEEN ? AND ?" {
		t.Errorf("Unexpected SQL result, got: %s", builder.GetSql())
	}

	builder = builder.NewSelect()
	builder.
		Table("users", "*").
		Where("dob BETWEEN ? AND ?", "1995-02-01", "2000-01-01")

	if builder.GetSql() != "SELECT * FROM users WHERE dob BETWEEN ? AND ?" {
		t.Errorf("Unexpected SQL result, got: %s", builder.GetSql())
	}
}

func TestWhereOr(t *testing.T) {
	builder = New("sqlite", db)
	builder.NewSelect()

	builder.
		Table("users", "*").
		Where("age >= ? OR email = ?", 18, "johndoe@example.com")

	if builder.GetSql() != "SELECT * FROM users WHERE age >= ? OR email = ?" {
		t.Errorf("Unexpected SQL result, got: %s", builder.GetSql())
	}
}

func TestJoin(t *testing.T) {
	builder = New("sqlite", db)
	builder.NewSelect()

	builder.
		Table("users", "*").
		Join("roles", "users.id", "=", "roles.user_id")

	if builder.GetSql() != "SELECT * FROM users INNER JOIN roles ON users.id = roles.user_id" {
		t.Errorf("Unexpected SQL result, got: %s", builder.GetSql())
	}
}

func TestLeftJoin(t *testing.T) {
	builder = New("sqlite", db)
	builder.NewSelect()

	builder.
		Table("users", "*").
		LeftJoin("roles", "users.id", "=", "roles.user_id")

	if builder.GetSql() != "SELECT * FROM users LEFT JOIN roles ON users.id = roles.user_id" {
		t.Errorf("Unexpected SQL result, got: %s", builder.GetSql())
	}
}

func TestRightJoin(t *testing.T) {
	builder = New("sqlite", db)
	builder.NewSelect()

	builder.
		Table("users", "*").
		RightJoin("roles", "users.id", "=", "roles.user_id")

	if builder.GetSql() != "SELECT * FROM users RIGHT JOIN roles ON users.id = roles.user_id" {
		t.Errorf("Unexpected SQL result, got: %s", builder.GetSql())
	}
}

func TestWhereExists(t *testing.T) {
	builder = New("sqlite", db)
	builder.NewSelect()

	builder.
		Table("users", "*").
		WhereExists(func(b Builder) *SQLBuilder {
			return b.Table("roles", "user_id").Where("users.id = roles.user_id")
		})

	if builder.GetSql() != "SELECT * FROM users WHERE EXISTS (SELECT user_id FROM roles WHERE users.id = roles.user_id)" {
		t.Errorf("Unexpected SQL result, got: %s", builder.GetSql())
	}
}

func TestWhereFuncSubquery(t *testing.T) {
	builder = New("sqlite", db)
	builder.NewSelect()

	builder.
		Table("users", "*").
		WhereFunc("email = ", func(b Builder) *SQLBuilder {
			return b.Table("roles", "user_id").Where("users.id = roles.user_id")
		})

	if builder.GetSql() != "SELECT * FROM users WHERE email = (SELECT user_id FROM roles WHERE users.id = roles.user_id)" {
		t.Errorf("Unexpected SQL result, got: %s", builder.GetSql())
	}
}

func TestExecute(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	user := new(User)
	builder = New("sqlite", db)
	builder.NewSelect()

	err := builder.Table("users", "*").Find(user, context.Background())

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
	builder = New("sqlite", db)
	builder.NewSelect()

	err := builder.Table("users", "*").Where("email = ?", "daniel@example.com").Limit(1).Find(user, context.Background())

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
	builder = New("sqlite", db)
	builder.NewSelect()

	err := builder.
		Table("users", "*").
		Where("age < ? AND email like ?", 30, "%@example.com").
		Get(&users, context.Background())

	if err != nil {
		t.Error(err)
	}

	if len(users) != 3 {
		t.Errorf("Expected return %d of users, but got %d", 3, len(users))
	}

}
