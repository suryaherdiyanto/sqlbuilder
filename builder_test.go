package sqlbuilder

import (
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
	builder.Select("*").Table("users")

	if sql, _ := builder.GetSql(); sql != "SELECT * FROM users" {
		t.Errorf("Unexpected SQL result, got: %s", sql)
	}
}

func TestWithWhere(t *testing.T) {
	builder = New("sqlite", db)
	builder.Select("*").Table("users").Where("email", OperatorEqual, "johndoe@gmail.com")

	if sql, _ := builder.GetSql(); sql != "SELECT * FROM users WHERE email = ?" {
		t.Errorf("Unexpected SQL result, got: %s", sql)
	}
}

func TestWithMultipleWhere(t *testing.T) {
	builder = New("sqlite", db)
	builder.Select("*")

	builder.
		Table("users").
		Where("email", "=", "johndoe@gmail.com").
		Where("access_role", "<", 3)

	if sql, _ := builder.GetSql(); sql != "SELECT * FROM users WHERE email = ? AND access_role < ?" {
		t.Errorf("Unexpected SQL result, got: %s", sql)
	}
}

// func TestWhereIn(t *testing.T) {
// 	builder = New("sqlite", db)
// 	builder.Select("*")

// 	builder.Table("users").Where("email IN(?, ?)", "johndoe@example.com", "amal@example.com")

// 	if sql, _ := builder.GetSql(); sql != "SELECT * FROM users WHERE email IN(?, ?)" {
// 		t.Errorf("Unexpected SQL result, got: %s", sql)
// 	}
// }

// func TestWhereBetween(t *testing.T) {
// 	builder = New("sqlite", db)
// 	builder.Select("*")

// 	builder.
// 		Table("users").
// 		Where("age BETWEEN ? AND ?", 5, 10)

// 	if sql, _ := builder.GetSql(); sql != "SELECT * FROM users WHERE age BETWEEN ? AND ?" {
// 		t.Errorf("Unexpected SQL result, got: %s", sql)
// 	}

// 	builder = builder.Select("*")
// 	builder.
// 		Table("users").
// 		Where("dob BETWEEN ? AND ?", "1995-02-01", "2000-01-01")

// 	if sql, _ := builder.GetSql(); sql != "SELECT * FROM users WHERE dob BETWEEN ? AND ?" {
// 		t.Errorf("Unexpected SQL result, got: %s", sql)
// 	}
// }

// func TestWhereOr(t *testing.T) {
// 	builder = New("sqlite", db)
// 	builder.Select("*")

// 	builder.
// 		Table("users").
// 		Where("age >= ? OR email = ?", 18, "johndoe@example.com")

// 	if sql, _ := builder.GetSql(); sql != "SELECT * FROM users WHERE age >= ? OR email = ?" {
// 		t.Errorf("Unexpected SQL result, got: %s", sql)
// 	}
// }

// func TestJoin(t *testing.T) {
// 	builder = New("sqlite", db)
// 	builder.Select("*")

// 	builder.
// 		Table("users").
// 		Join("roles", "users.id", "=", "roles.user_id")

// 	if sql, _ := builder.GetSql(); sql != "SELECT * FROM users INNER JOIN roles ON users.id = roles.user_id" {
// 		t.Errorf("Unexpected SQL result, got: %s", sql)
// 	}
// }

// func TestLeftJoin(t *testing.T) {
// 	builder = New("sqlite", db)
// 	builder.Select("*")

// 	builder.
// 		Table("users").
// 		LeftJoin("roles", "users.id", "=", "roles.user_id")

// 	if sql, _ := builder.GetSql(); sql != "SELECT * FROM users LEFT JOIN roles ON users.id = roles.user_id" {
// 		t.Errorf("Unexpected SQL result, got: %s", sql)
// 	}
// }

// func TestRightJoin(t *testing.T) {
// 	builder = New("sqlite", db)
// 	builder.Select("*")

// 	builder.
// 		Table("users").
// 		RightJoin("roles", "users.id", "=", "roles.user_id")

// 	if sql, _ := builder.GetSql(); sql != "SELECT * FROM users RIGHT JOIN roles ON users.id = roles.user_id" {
// 		t.Errorf("Unexpected SQL result, got: %s", sql)
// 	}
// }

// func TestWhereExists(t *testing.T) {
// 	builder = New("sqlite", db)
// 	builder.Select("*")

// 	builder.
// 		Table("users").
// 		WhereExists(func(b Builder) *SQLBuilder {
// 			return b.Select("*").Table("roles").Where("users.id = roles.user_id")
// 		})

// 	if sql, _ := builder.GetSql(); sql != "SELECT * FROM users WHERE EXISTS (SELECT * FROM roles WHERE users.id = roles.user_id)" {
// 		t.Errorf("Unexpected SQL result, got: %s", sql)
// 	}
// }

// func TestWhereFuncSubquery(t *testing.T) {
// 	builder = New("sqlite", db)
// 	builder.Select("*")

// 	builder.
// 		Table("users").
// 		WhereFunc("email =", func(b Builder) *SQLBuilder {
// 			return b.Select("user_id").Table("roles").Where("users.id = roles.user_id")
// 		})

// 	if sql, _ := builder.GetSql(); sql != "SELECT * FROM users WHERE email = (SELECT user_id FROM roles WHERE users.id = roles.user_id)" {
// 		t.Errorf("Unexpected SQL result, got: %s", sql)
// 	}
// }

// func TestExecuteUpdateStatement(t *testing.T) {
// 	teardownSuite := setupSuite(t)
// 	defer teardownSuite(t)

// 	builder = New("sqlite", db)
// 	type UserRequest struct {
// 		Username string `db:"username"`
// 		Age      int    `db:"age"`
// 	}

// 	user := &UserRequest{
// 		Username: "johndoe",
// 		Age:      31,
// 	}
// 	result, err := builder.Table("users").Where("id = ?", 1).Update(user)

// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	expected := "UPDATE users SET username = ?, age = ? WHERE id = ?"
// 	if sql, _ := builder.GetSql(); sql != expected {
// 		t.Errorf("Unexpected SQL result, got: %s", sql)
// 	}

// 	if rowsAffected, err := result.RowsAffected(); err != nil {
// 		t.Error(err)

// 		if rowsAffected <= 0 {
// 			t.Errorf("Expected rows affected to be greater than 0, but got: %d", rowsAffected)
// 		}
// 	}
// }

// func TestDeleteStatement(t *testing.T) {
// 	teardownSuite := setupSuite(t)
// 	defer teardownSuite(t)

// 	builder = New("sqlite", db)
// 	result, err := builder.Table("users").Where("username = ?", "johndoe").Delete()

// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	expected := "DELETE FROM users WHERE username = ?"
// 	if sql, _ := builder.GetSql(); sql != expected {
// 		t.Errorf("Unexpected SQL result, got: %s", sql)
// 	}

// 	if rowsAffected, err := result.RowsAffected(); err != nil {
// 		t.Error(err)

// 		if rowsAffected <= 0 {
// 			t.Errorf("Expected rows affected to be greater than 0, but got: %d", rowsAffected)
// 		}
// 	}
// }

// func TestExecuteTransaction(t *testing.T) {
// 	teardownSuite := setupSuite(t)
// 	defer teardownSuite(t)

// 	builder = New("sqlite", db)
// 	err := builder.Begin(func(b *SQLBuilder) error {
// 		type UserRequest struct {
// 			Username string `db:"username"`
// 			Age      int    `db:"age"`
// 			Email    string `db:"email"`
// 		}
// 		user := &UserRequest{
// 			Username: "johncena",
// 			Email:    "johncena@example.com",
// 			Age:      35,
// 		}
// 		result, err := b.Table("users").Insert(user)

// 		if err != nil {
// 			return err
// 		}

// 		newUser := &User{}
// 		id, err := result.LastInsertId()
// 		if err != nil {
// 			return err
// 		}

// 		if err = b.Select("*").Table("users").Where("id = ?", id).Find(newUser, context.Background()); err != nil {
// 			return err
// 		}

// 		type UpdateRequest struct {
// 			Age int `db:"age"`
// 		}
// 		update := &UpdateRequest{Age: 40}
// 		if _, err = b.Table("users").Where("id = ?", id).Update(update); err != nil {
// 			return err
// 		}

// 		return nil
// 	})

// 	if err != nil {
// 		t.Error(err)
// 	}
// }

// func TestExecute(t *testing.T) {
// 	teardownSuite := setupSuite(t)
// 	defer teardownSuite(t)

// 	user := new(User)
// 	builder = New("sqlite", db)
// 	err := builder.Select("*").Table("users").Find(user, context.Background())

// 	if err != nil {
// 		t.Error(err)
// 	}

// 	if user.Email != "johndoe@example.com" {
// 		t.Errorf("Expected johndoe@example.com, but got: %s", user.Email)
// 	}

// }

// func TestExecuteInsert(t *testing.T) {
// 	teardownSuite := setupSuite(t)
// 	defer teardownSuite(t)

// 	type UserRequest struct {
// 		Username string `db:"username"`
// 		Email    string `db:"email"`
// 		Age      int    `db:"age"`
// 	}

// 	user := &UserRequest{
// 		Username: "johndoe",
// 		Email:    "johndoe@example.com",
// 		Age:      35,
// 	}

// 	builder = New("sqlite", db)
// 	result, err := builder.Table("users").Insert(user)

// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	expected := "INSERT INTO users(username,email,age) VALUES(?,?,?)"

// 	if sql, _ := builder.GetSql(); sql != expected {
// 		t.Errorf("Unexpected SQL result, got: %s", sql)
// 	}

// 	if id, err := result.LastInsertId(); err != nil {
// 		t.Error(err)

// 		if id <= 0 {
// 			t.Errorf("Expected last insert id to be greater than 0, but got: %d", id)
// 		}
// 	}

// }

// func TestExecuteWhere(t *testing.T) {
// 	teardownSuite := setupSuite(t)
// 	defer teardownSuite(t)

// 	user := new(User)
// 	builder = New("sqlite", db)
// 	builder.Select()

// 	err := builder.Select("*").Table("users").Where("email = ?", "daniel@example.com").Limit(1).Find(user, context.Background())

// 	if err != nil {
// 		t.Error(err)
// 	}

// 	if user.Email != "daniel@example.com" {
// 		t.Errorf("Expected daniel@example.com, but got: %s", user.Email)
// 	}
// }

// func TestWhereAnd(t *testing.T) {
// 	teardownSuite := setupSuite(t)
// 	defer teardownSuite(t)

// 	var users []User
// 	builder = New("sqlite", db)
// 	builder.Select()

// 	err := builder.
// 		Select("*").
// 		Table("users").
// 		Where("age < ? AND email like ?", 30, "%@example.com").
// 		Get(&users, context.Background())

// 	if err != nil {
// 		t.Error(err)
// 	}

// 	if len(users) != 3 {
// 		t.Errorf("Expected return %d of users, but got %d", 3, len(users))
// 	}

// }
