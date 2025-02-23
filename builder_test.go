package sqlbuilder

import (
	"context"
	"database/sql"
	"sync"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

var db, _ = sql.Open("sqlite3", ":memory:")
type User struct {
	Id int `db:"id"`
	Username string `db:"username"`
	Email string `db:"email"`
	Age int `db:"age"`
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
		db.Exec("DROP table users");
	}
}

func TestBuilder(t *testing.T) {
	builder := NewSQLBuilder("sqlite", db)

	builder.Select("*").Table("users")

	if builder.GetSql() != "SELECT * FROM users" {
		t.Errorf("Unexpected SQL result, got: %s", builder.GetSql())
	}
}

func TestWithWhere(t *testing.T) {
	builder := NewSQLBuilder("sqlite", db)

	builder.
		Select("*").
		Table("users").
		Where("email", Eq, "johndoe@gmail.com")

	if builder.GetSql() != "SELECT * FROM users WHERE email = ?" {
		t.Errorf("Unexpected SQL result, got: %s", builder.GetSql())
	}
}

func TestWithMultipleWhere(t *testing.T) {
	builder := NewSQLBuilder("sqlite", db)

	builder.
		Select("*").
		Table("users").
		Where("email", Eq, "johndoe@gmail.com").
		Where("access_role", Lt, 3)

	if builder.GetSql() != "SELECT * FROM users WHERE email = ? AND access_role < ?" {
		t.Errorf("Unexpected SQL result, got: %s", builder.GetSql())
	}
}

func TestWhereIn(t *testing.T) {
	builder := NewSQLBuilder("sqlite", db)

	builder.
		Select("*").
		Table("users").
		WhereIn("email", []string{"john@gmail.com", "admin@example.com"})

	if builder.GetSql() != "SELECT * FROM users WHERE email IN(?, ?)" {
		t.Errorf("Unexpected SQL result, got: %s", builder.GetSql())
	}
}

func TestWhereInWithNumber(t *testing.T) {
	builder := NewSQLBuilder("sqlite", db)

	builder.
		Select("*").
		Table("users").
		WhereIn("some_column", []int{1, 2, 3})

	if builder.GetSql() != "SELECT * FROM users WHERE some_column IN(?, ?, ?)" {
		t.Errorf("Unexpected SQL result, got: %s", builder.GetSql())
	}
}

func TestWhereBetween(t *testing.T) {
	builder := NewSQLBuilder("sqlite", db)

	builder.
		Select("*").
		Table("users").
		WhereBetween("age", 5, 10)

	if builder.GetSql() != "SELECT * FROM users WHERE age BETWEEN ? AND ?" {
		t.Errorf("Unexpected SQL result, got: %s", builder.GetSql())
	}

	builder = NewSQLBuilder("sqlite", db)
	builder.
		Select("*").
		Table("users").
		WhereBetween("dob", "1995-02-01", "2000-01-01")

	if builder.GetSql() != "SELECT * FROM users WHERE dob BETWEEN ? AND ?" {
		t.Errorf("Unexpected SQL result, got: %s", builder.GetSql())
	}
}

func TestWhereOr(t *testing.T) {
	builder := NewSQLBuilder("sqlite", db)

	builder.
		Select("*").
		Table("users").
		Where("age", Gte, 18).
		OrWhere("email", Eq, "admin@example.com")

	if builder.GetSql() != "SELECT * FROM users WHERE age >= ? OR email = ?" {
		t.Errorf("Unexpected SQL result, got: %s", builder.GetSql())
	}
}

func TestExecute(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	user := new(User)
	builder := NewSQLBuilder("sqlite", db)

	err := builder.Select("*").Table("users").Scan(user, context.Background())

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
	builder := NewSQLBuilder("sqlite", db)

	err := builder.Select("*").Table("users").Where("email", Eq, "daniel@example.com").Limit(1).Scan(user, context.Background())

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
	builder := NewSQLBuilder("sqlite", db)

	err := builder.
				Select("*").
				Table("users").
				Where("age", Lt, 30).
				Where("email", "like", "%@example.com").
				ScanAll(&users, context.Background())

	if err != nil {
		t.Error(err)
	}

	if len(users) != 3 {
		t.Errorf("Expected return %d of users, but got %d", 3, len(users))
	}

}

func TestGoRoutineSafe(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	var users []User
	user := new(User)
	var wg sync.WaitGroup

	wg.Add(1)
	go func(wg *sync.WaitGroup, builder *SQLBuilder) {
		err := builder.
				Select("*").
				Table("users").
				Where("age", Lt, 30).
				Where("email", "like", "%@example.com").
				ScanAll(&users, context.Background())

		if err != nil {
			t.Error(err)
		}
		wg.Done()
	}(&wg, NewSQLBuilder("sqlite", db))
	wg.Add(1)
	go func(wg *sync.WaitGroup, builder *SQLBuilder) {
		err := builder.
				Select("*").
				Table("users").
				Where("email", Eq, "johndoe@gmail.com").
				Scan(user, context.Background())

		if err != nil {
			t.Error(err)
		}
		wg.Done()
	}(&wg, NewSQLBuilder("sqlite", db))
	wg.Wait()
}
