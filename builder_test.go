package sqlbuilder

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/suryaherdiyanto/sqlbuilder/clause"
	"github.com/suryaherdiyanto/sqlbuilder/dialect"
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
	dialect := dialect.New("?", "`", "`")
	builder = New(dialect, db)
	builder.Table("users").Select("*")

	if sql, _ := builder.GetSql(); sql != "SELECT * FROM `users`" {
		t.Errorf("Unexpected SQL result, got: %s", sql)
	}

	builder = New(dialect, db)
	builder.Table("users").Select("id", "username", "email")

	if sql, _ := builder.GetSql(); sql != "SELECT `id`,`username`,`email` FROM `users`" {
		t.Errorf("Unexpected SQL result, got: %s", sql)
	}

}

func TestWithWhere(t *testing.T) {
	dialect := dialect.New("?", "`", "`")
	builder = New(dialect, db)
	builder.Table("users").Select("*").Where("email", clause.OperatorEqual, "johndoe@gmail.com")

	if sql, _ := builder.GetSql(); sql != "SELECT * FROM `users` WHERE `email` = ?" {
		t.Errorf("Unexpected SQL result, got: %s", sql)
	}
}

func TestWithMultipleWhere(t *testing.T) {
	dialect := dialect.New("?", "`", "`")
	builder = New(dialect, db)
	builder.Table("users").Select("*")

	builder.
		Where("email", "=", "johndoe@gmail.com").
		Where("access_role", "<", 3)

	if sql, _ := builder.GetSql(); sql != "SELECT * FROM `users` WHERE `email` = ? AND `access_role` < ?" {
		t.Errorf("Unexpected SQL result, got: %s", sql)
	}
}

func TestWhereIn(t *testing.T) {
	dialect := dialect.New("?", "`", "`")
	builder = New(dialect, db)
	builder.Table("users")

	builder.Select("*").
		WhereIn("email", []any{"johndoe@example.com", "amal@example.com"})

	if sql, _ := builder.GetSql(); sql != "SELECT * FROM `users` WHERE `email` IN(?,?)" {
		t.Errorf("Unexpected SQL result, got: %s", sql)
	}
}

func TestWhereBetween(t *testing.T) {
	dialect := dialect.New("?", "`", "`")
	builder = New(dialect, db)
	builder.Table("users")

	builder.
		Select("*").
		WhereBetween("age", 5, 10)

	if sql, _ := builder.GetSql(); sql != "SELECT * FROM `users` WHERE `age` BETWEEN ? AND ?" {
		t.Errorf("Unexpected SQL result, got: %s", sql)
	}

	builder = builder.Table("users")
	builder.
		Select("*").
		WhereBetween("dob", "1995-02-01", "2000-01-01")

	if sql, _ := builder.GetSql(); sql != "SELECT * FROM `users` WHERE `dob` BETWEEN ? AND ?" {
		t.Errorf("Unexpected SQL result, got: %s", sql)
	}
}

func TestWhereOr(t *testing.T) {
	dialect := dialect.New("?", "`", "`")
	builder = New(dialect, db)
	builder.Table("users").Select("*")

	builder.
		Where("age", clause.OperatorGreatherThanEqual, 18).
		WhereOr("email", clause.OperatorEqual, "johndoe@example.com")

	if sql, _ := builder.GetSql(); sql != "SELECT * FROM `users` WHERE `age` >= ? OR `email` = ?" {
		t.Errorf("Unexpected SQL result, got: %s", sql)
	}
}

func TestJoin(t *testing.T) {
	dialect := dialect.New("?", "`", "`")
	builder = New(dialect, db)
	builder.Table("users").Select("*")

	builder.
		Join("roles", "id", "=", "user_id")

	if sql, _ := builder.GetSql(); sql != "SELECT * FROM `users` INNER JOIN `roles` ON `users`.`id` = `roles`.`user_id`" {
		t.Errorf("Unexpected SQL result, got: %s", sql)
	}
}

func TestLeftJoin(t *testing.T) {
	dialect := dialect.New("?", "`", "`")
	builder = New(dialect, db)
	builder.Table("users").Select("*")

	builder.
		LeftJoin("roles", "id", clause.OperatorEqual, "user_id")
	if sql, _ := builder.GetSql(); sql != "SELECT * FROM `users` LEFT JOIN `roles` ON `users`.`id` = `roles`.`user_id`" {
		t.Errorf("Unexpected SQL result, got: %s", sql)
	}
}

func TestRightJoin(t *testing.T) {
	dialect := dialect.New("?", "`", "`")
	builder = New(dialect, db)
	builder.Table("users").Select("*")

	builder.
		RightJoin("roles", "id", clause.OperatorEqual, "user_id")
	if sql, _ := builder.GetSql(); sql != "SELECT * FROM `users` RIGHT JOIN `roles` ON `users`.`id` = `roles`.`user_id`" {
		t.Errorf("Unexpected SQL result, got: %s", sql)
	}
}

func TestWhereExists(t *testing.T) {
	dialect := dialect.New("?", "`", "`")
	builder = New(dialect, db)
	builder.Table("users").Select("*")

	builder.
		WhereExists(func(b Builder) *SQLBuilder {
			return b.Table("roles").Select("*").Where("users.id", "=", "roles.user_id")
		})

	if sql, _ := builder.GetSql(); sql != "SELECT * FROM `users` WHERE EXISTS (SELECT * FROM `roles` WHERE `users`.`id` = ?)" {
		t.Errorf("Unexpected SQL result, got: %s", sql)
	}
}

func TestWhereFuncSubquery(t *testing.T) {
	dialect := dialect.New("?", "`", "`")
	builder = New(dialect, db)
	builder.Table("users").Select("*")

	builder.
		WhereFunc("email", "=", func(b Builder) *SQLBuilder {
			return b.Table("roles").Select("user_id").Where("users.id", "=", "roles.user_id")
		})

	if sql, _ := builder.GetSql(); sql != "SELECT * FROM `users` WHERE `email` = (SELECT `user_id` FROM `roles` WHERE `users`.`id` = ?)" {
		t.Errorf("Unexpected SQL result, got: %s", sql)
	}
}

func TestGroupBy(t *testing.T) {
	dialect := dialect.New("?", "`", "`")
	builder = New(dialect, db)
	builder.Table("users").Select("*")

	builder.
		GroupBy("age", "role")

	if sql, _ := builder.GetSql(); sql != "SELECT * FROM `users` GROUP BY `age`,`role`" {
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
	dialect := dialect.New("?", "`", "`")
	builder := New(dialect, dba)
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

	dialect := dialect.New("?", "`", "`")
	builder := New(dialect, dba)
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

	dialect := dialect.New("?", "`", "`")
	builder := New(dialect, dba)
	builder = builder.
		Table("users").
		Select("*").
		WhereFunc("age", "=", func(b Builder) *SQLBuilder {
			return b.Table("users").Select("MIN(age) as min")
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

	dialect := dialect.New("?", "`", "`")
	builder := New(dialect, dba)

	id, err := builder.Table("users").Insert(map[string]any{
		"username": "alice",
		"email":    "alice@example.com",
		"age":      29,
	})

	if err != nil {
		t.Fatal(err)
	}

	if id == 0 {
		t.Error("Expected id not to be 0")
	}

}

func TestExecuteInsertWithStructData(t *testing.T) {
	dba, err := sql.Open("sqlite3", ":memory:")
	type UserData struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Age      uint64 `json:"age"`
	}

	if err != nil {
		t.Fatal(err)
	}

	err = seed(dba)

	if err != nil {
		t.Fatal(err)
	}

	dialect := dialect.New("?", "`", "`")
	builder := New(dialect, dba)
	data := UserData{
		Username: "foo",
		Email:    "foobar@gmail.com",
		Age:      23,
	}
	id, err := builder.Table("users").Insert(data)

	if err != nil {
		t.Fatal(err)
	}

	if id == 0 {
		t.Error("Expected id not to be 0")
	}

}

func TestInsertMultipleRows(t *testing.T) {
	dba, err := sql.Open("sqlite3", ":memory:")

	if err != nil {
		t.Fatal(err)
	}

	err = seed(dba)

	if err != nil {
		t.Fatal(err)
	}

	dialect := dialect.New("?", "`", "`")
	builder := New(dialect, dba)

	res, err := builder.Table("users").InsertMany([]map[string]any{
		{
			"username": "alice",
			"email":    "alice@example.com",
			"age":      29,
		},
		{
			"username": "john doe",
			"email":    "johndoe@example.com",
			"age":      20,
		},
	}).Exec()

	if err != nil {
		t.Fatal(err)
	}

	rows, err := res.RowsAffected()

	if err != nil {
		t.Fatal(err)
	}

	if rows != 2 {
		t.Errorf("Expected rows affected to be greater than 0, but got: %d", rows)
	}
}

func TestExecuteUpdateStatement(t *testing.T) {
	dba, err := sql.Open("sqlite3", ":memory:")

	if err != nil {
		t.Fatal(err)
	}

	err = seed(dba)

	dialect := dialect.New("?", "`", "`")
	builder = New(dialect, dba)

	builder = builder.Table("users").Where("id", clause.OperatorEqual, 1).Update(map[string]any{
		"username": "john_doe_updated",
		"age":      36,
	})
	result, err := builder.Exec()

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

	dialect := dialect.New("?", "`", "`")
	builder = New(dialect, dba)
	result, err := builder.Table("users").Where("username", clause.OperatorEqual, "johndoe").Delete().Exec()

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

func TestExecuteTransaction(t *testing.T) {
	dba, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	err = seed(dba)

	if err != nil {
		t.Fatal(err)
	}

	dialect := dialect.New("?", "`", "`")
	builder = New(dialect, dba)
	err = builder.Begin(func(b *SQLBuilder) error {
		type UserRequest struct {
			Username string `json:"username"`
			Age      int    `json:"age"`
			Email    string `json:"email"`
		}
		user := UserRequest{
			Username: "johncena",
			Email:    "johncena@example.com",
			Age:      35,
		}
		lastInsertId, err := b.Table("users").Insert(user)

		if err != nil {
			return err
		}

		newUser := &User{}

		if err = b.Select("*").Table("users").Where("id", clause.OperatorEqual, lastInsertId).Get(&newUser); err != nil {
			return err
		}

		type UpdateRequest struct {
			Age int `db:"age"`
		}
		update := UpdateRequest{Age: 40}
		updateMap := map[string]any{}
		toMap(update, updateMap)

		if _, err = b.Table("users").Where("id", clause.OperatorEqual, lastInsertId).Update(updateMap).Exec(); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		t.Error(err)
	}
}

func TestExecute(t *testing.T) {
	dba, err := sql.Open("sqlite3", ":memory:")

	if err != nil {
		t.Fatal(err)
	}

	err = seed(dba)
	if err != nil {
		t.Fatal(err)
	}

	user := new(User)
	dialect := dialect.New("?", "`", "`")
	builder = New(dialect, dba)
	err = builder.Table("users").Select("*").Get(user)

	if err != nil {
		sql, _ := builder.GetSql()
		t.Errorf("Failed to execute select statement: %s", sql)
		t.Error(err)
	}

	if user.Email != "johndoe@example.com" {
		t.Errorf("Expected johndoe@example.com, but got: %s", user.Email)
	}

}

func TestExecuteWhere(t *testing.T) {
	dba, err := sql.Open("sqlite3", ":memory:")

	if err != nil {
		t.Fatal(err)
	}

	err = seed(dba)
	if err != nil {
		t.Fatal(err)
	}

	user := new(User)
	dialect := dialect.New("?", "`", "`")
	builder = New(dialect, dba)

	err = builder.Table("users").Select("*").Where("email", clause.OperatorEqual, "daniel@example.com").Limit(1).Get(user)

	if err != nil {
		t.Error(err)
	}

	if user.Email != "daniel@example.com" {
		t.Errorf("Expected daniel@example.com, but got: %s", user.Email)
	}
}

func TestWhereAnd(t *testing.T) {
	dba, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	err = seed(dba)
	if err != nil {
		t.Fatal(err)
	}

	var users []User
	dialect := dialect.New("?", "`", "`")
	builder = New(dialect, dba)
	builder.Select()

	err = builder.
		Table("users").
		Select("*").
		Where("age", clause.OperatorLessThan, 30).
		Where("email", clause.OperatorLike, "%@example.com").
		Get(&users)

	if err != nil {
		t.Error(err)
	}

	if len(users) != 3 {
		t.Errorf("Expected return %d of users, but got %d", 3, len(users))
	}

}
