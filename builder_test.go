package sqlbuilder

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB
var builder *SQLBuilder

type User struct {
	Id       int    `db:"id"`
	Username string `db:"username"`
	Email    string `db:"email"`
	Age      int    `db:"age"`
}

func seed(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE users(
			id integer primary key,
			username TEXT,
			email TEXT,
			age integer
		)
	`)

	if err != nil {
		return err
	}

	_, err = db.Exec(`
		INSERT INTO users values(null, 'johndoe', 'johndoe@example.com', 35);
		INSERT INTO users values(null, 'daniel', 'daniel@example.com', 32);
		INSERT INTO users values(null, 'samuel', 'samuel@example.com', 28);
		INSERT INTO users values(null, 'dirt', 'dirt@example.com', 20);
		INSERT INTO users values(null, 'chris', 'chris@example.com', 25);
	`)

	if err != nil {
		return err
	}

	return nil
}

func TestBuilder(t *testing.T) {
	builder = New("sqlite", db)
	builder.Table("users").Select("*")

	if sql, _ := builder.GetSql(); sql != "SELECT * FROM users" {
		t.Errorf("Unexpected SQL result, got: %s", sql)
	}

	builder = New("sqlite", db)
	builder.Table("users").Select("id", "username", "email")

	if sql, _ := builder.GetSql(); sql != "SELECT id,username,email FROM users" {
		t.Errorf("Unexpected SQL result, got: %s", sql)
	}

}

func TestWithWhere(t *testing.T) {
	builder = New("sqlite", db)
	builder.Table("users").Select("*").Where("email", OperatorEqual, "johndoe@gmail.com")

	if sql, _ := builder.GetSql(); sql != "SELECT * FROM users WHERE email = ?" {
		t.Errorf("Unexpected SQL result, got: %s", sql)
	}
}

func TestWithMultipleWhere(t *testing.T) {
	builder = New("sqlite", db)
	builder.Table("users").Select("*")

	builder.
		Where("email", "=", "johndoe@gmail.com").
		Where("access_role", "<", 3)

	if sql, _ := builder.GetSql(); sql != "SELECT * FROM users WHERE email = ? AND access_role < ?" {
		t.Errorf("Unexpected SQL result, got: %s", sql)
	}
}

func TestWhereIn(t *testing.T) {
	builder = New("sqlite", db)
	builder.Table("users")

	builder.Select("*").
		WhereIn("email", []any{"johndoe@example.com", "amal@example.com"})

	if sql, _ := builder.GetSql(); sql != "SELECT * FROM users WHERE email IN(?,?)" {
		t.Errorf("Unexpected SQL result, got: %s", sql)
	}
}

func TestWhereBetween(t *testing.T) {
	builder = New("sqlite", db)
	builder.Table("users")

	builder.
		Select("*").
		WhereBetween("age", 5, 10)

	if sql, _ := builder.GetSql(); sql != "SELECT * FROM users WHERE age BETWEEN ? AND ?" {
		t.Errorf("Unexpected SQL result, got: %s", sql)
	}

	builder = builder.Table("users")
	builder.
		Select("*").
		WhereBetween("dob", "1995-02-01", "2000-01-01")

	if sql, _ := builder.GetSql(); sql != "SELECT * FROM users WHERE dob BETWEEN ? AND ?" {
		t.Errorf("Unexpected SQL result, got: %s", sql)
	}
}

func TestWhereOr(t *testing.T) {
	builder = New("sqlite", db)
	builder.Table("users").Select("*")

	builder.
		Where("age", OperatorGreatherThanEqual, 18).
		WhereOr("email", OperatorEqual, "johndoe@example.com")

	if sql, _ := builder.GetSql(); sql != "SELECT * FROM users WHERE age >= ? OR email = ?" {
		t.Errorf("Unexpected SQL result, got: %s", sql)
	}
}

func TestJoin(t *testing.T) {
	builder = New("sqlite", db)
	builder.Table("users").Select("*")

	builder.
		Join("roles", "id", "=", "user_id")

	if sql, _ := builder.GetSql(); sql != "SELECT * FROM users INNER JOIN roles ON users.id = roles.user_id" {
		t.Errorf("Unexpected SQL result, got: %s", sql)
	}
}

func TestLeftJoin(t *testing.T) {
	builder = New("sqlite", db)
	builder.Table("users").Select("*")

	builder.
		LeftJoin("roles", "id", OperatorEqual, "user_id")

	if sql, _ := builder.GetSql(); sql != "SELECT * FROM users LEFT JOIN roles ON users.id = roles.user_id" {
		t.Errorf("Unexpected SQL result, got: %s", sql)
	}
}

func TestRightJoin(t *testing.T) {
	builder = New("sqlite", db)
	builder.Table("users").Select("*")

	builder.
		RightJoin("roles", "id", "=", "user_id")

	if sql, _ := builder.GetSql(); sql != "SELECT * FROM users RIGHT JOIN roles ON users.id = roles.user_id" {
		t.Errorf("Unexpected SQL result, got: %s", sql)
	}
}

func TestWhereExists(t *testing.T) {
	builder = New("sqlite", db)
	builder.Table("users").Select("*")

	builder.
		WhereExists(func(b Builder) *SQLBuilder {
			return b.Table("roles").Select("*").Where("users.id", "=", "roles.user_id")
		})

	if sql, _ := builder.GetSql(); sql != "SELECT * FROM users WHERE EXISTS (SELECT * FROM roles WHERE users.id = ?)" {
		t.Errorf("Unexpected SQL result, got: %s", sql)
	}
}

func TestWhereFuncSubquery(t *testing.T) {
	builder = New("sqlite", db)
	builder.Table("users").Select("*")

	builder.
		WhereFunc("email", "=", func(b Builder) *SQLBuilder {
			return b.Table("roles").Select("user_id").Where("users.id", "=", "roles.user_id")
		})

	if sql, _ := builder.GetSql(); sql != "SELECT * FROM users WHERE email = (SELECT user_id FROM roles WHERE users.id = ?)" {
		t.Errorf("Unexpected SQL result, got: %s", sql)
	}
}

func TestGroupBy(t *testing.T) {
	builder = New("sqlite", db)
	builder.Table("users").Select("*")

	builder.
		GroupBy("age", "role")

	if sql, _ := builder.GetSql(); sql != "SELECT * FROM users GROUP BY age,role" {
		t.Errorf("Unexpected SQL result, got: %s", sql)
	}
}

func TestExecuteSelectStatement(t *testing.T) {

	dba, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	err = seed(dba)
	if err != nil {
		t.Fatal(err)
	}

	var users []User
	builder := New("sqlite", dba)
	err = builder.Table("users").Select("*").Get(&users)

	if err != nil {
		t.Error(err)
	}

	if len(users) != 5 {
		t.Errorf("Expected return %d of users, but got %d", 5, len(users))
	}
}

func TestExecuteWithWhereStatement(t *testing.T) {

	user := new(User)
	dba, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	err = seed(dba)

	if err != nil {
		t.Fatal(err)
	}

	builder := New("sqlite", dba)
	err = builder.Table("users").Select("*").Where("id", "=", 1).Limit(1).Get(user)

	if err != nil {
		arguments := builder.GetArguments()
		stmt, _ := builder.GetSql()

		t.Errorf("SQL: %s", stmt)
		t.Errorf("Arguments: %v", arguments)
		t.Error(err)
	}

	if user.Email != "johndoe@example.com" {
		t.Errorf("Expected user email is johndoe@example.com, but got: %s", user.Email)
	}
}

func TestExecuteSubQuery(t *testing.T) {
	user := new(User)
	dba, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	err = seed(dba)

	if err != nil {
		t.Fatal(err)
	}

	builder := New("sqlite", dba)
	builder = builder.
		Table("users").
		Select("*").
		WhereFunc("age", "=", func(b Builder) *SQLBuilder {
			return b.Table("users").Select("MIN(age)")
		}).
		Limit(1)

	err = builder.Get(user)

	if err != nil {
		arguments := builder.GetArguments()
		stmt, _ := builder.GetSql()

		t.Errorf("SQL: %s", stmt)
		t.Errorf("Arguments: %v", arguments)
		t.Error(err)
	}

	if user.Age != 20 {
		t.Errorf("Expected user age is 20, but got: %d", user.Age)
	}
}

func TestExecuteInsert(t *testing.T) {
	dba, err := sql.Open("sqlite3", ":memory:")

	if err != nil {
		t.Fatal(err)
	}

	err = seed(dba)

	if err != nil {
		t.Fatal(err)
	}

	builder := New("sqlite", dba)

	res, err := builder.Table("users").Insert([]map[string]any{
		{
			"username": "alice",
			"email":    "alice@example.com",
			"age":      29,
		},
	}).Exec()

	if err != nil {
		t.Fatal(err)
	}

	rows, err := res.RowsAffected()

	if err != nil {
		t.Fatal(err)
	}

	if rows <= 0 {
		t.Errorf("Expected rows affected to be greater than 0, but got: %d", rows)
	}

}

func TestExecuteUpdateStatement(t *testing.T) {
	dba, err := sql.Open("sqlite3", ":memory:")

	if err != nil {
		t.Fatal(err)
	}

	err = seed(dba)

	builder = New("sqlite", dba)

	result, err := builder.Table("users").Where("id", OperatorEqual, 1).Update(map[string]any{
		"username": "john_doe_updated",
		"age":      36,
	}).Exec()

	if sql, err := builder.GetSql(); err != nil {
		t.Fatal(err)
		t.Errorf("Unexpected SQL result, got: %s", sql)
	}

	if err != nil {
		t.Fatal(err)
	}

	if rowsAffected, err := result.RowsAffected(); err != nil {
		t.Error(err)

		if rowsAffected <= 0 {
			t.Errorf("Expected rows affected to be greater than 0, but got: %d", rowsAffected)
		}
	}
}

func TestExecuteDeleteStatement(t *testing.T) {
	dba, err := sql.Open("sqlite3", ":memory:")

	if err != nil {
		t.Fatal(err)
	}

	err = seed(dba)
	if err != nil {
		t.Fatal(err)
	}

	builder = New("sqlite", dba)
	result, err := builder.Table("users").Where("username", OperatorEqual, "johndoe").Delete().Exec()

	if err != nil {
		t.Fatal(err)
	}

	if sql, err := builder.GetSql(); err != nil {
		t.Fatal(err)
		t.Errorf("Unexpected SQL result, got: %s", sql)
	}

	if rowsAffected, err := result.RowsAffected(); err != nil {
		t.Error(err)

		if rowsAffected <= 0 {
			t.Errorf("Expected rows affected to be greater than 0, but got: %d", rowsAffected)
		}
	}
}

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
