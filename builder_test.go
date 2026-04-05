package sqlbuilder

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/suryaherdiyanto/sqlbuilder/clause"
	"github.com/suryaherdiyanto/sqlbuilder/dialect"
)

var db *sql.DB
var builder *SQLBuilder

type User struct {
	Id        int       `db:"id"`
	Username  string    `db:"username"`
	Email     string    `db:"email"`
	Age       int       `db:"age"`
	CreatedAt time.Time `db:"created_at"`
}

func seed(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE users(
			id integer primary key,
			username TEXT,
			email TEXT,
			age integer,
			created_at datetime default current_timestamp
		)
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE TABLE roles(
			id integer primary key,
			name TEXT
		)
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE TABLE user_roles(
			id integer primary key,
			user_id integer,
			role_id integer
		)
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		INSERT INTO users values(null, 'johndoe', 'johndoe@example.com', 35, '2023-01-01 10:00:00');
		INSERT INTO users values(null, 'daniel', 'daniel@example.com', 32, '2023-01-02 10:00:00');
		INSERT INTO users values(null, 'samuel', 'samuel@example.com', 28, '2023-01-03 10:00:00');
		INSERT INTO users values(null, 'dirt', 'dirt@example.com', 20, '2023-01-04 10:00:00');
		INSERT INTO users values(null, 'chris', 'chris@example.com', 25, '2023-01-05 10:00:00');
		INSERT INTO users values(null, 'alice', 'alice@example.com', 29, '2023-01-06 10:00:00');
		INSERT INTO users values(null, 'bob', 'bob@example.com', 31, '2023-01-07 10:00:00');
		INSERT INTO users values(null, 'carol', 'carol@example.com', 27, '2023-01-08 10:00:00');
		INSERT INTO users values(null, 'eve', 'eve@example.com', 24, '2023-01-09 10:00:00');
		INSERT INTO users values(null, 'frank', 'frank@example.com', 38, '2023-01-10 10:00:00');
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		INSERT INTO roles values(1, 'admin');
		INSERT INTO roles values(2, 'staff');
		INSERT INTO roles values(3, 'manager');
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		INSERT INTO user_roles values(null, 1, 1);
		INSERT INTO user_roles values(null, 2, 2);
		INSERT INTO user_roles values(null, 3, 3);
		INSERT INTO user_roles values(null, 4, 2);
		INSERT INTO user_roles values(null, 5, 2);
		INSERT INTO user_roles values(null, 6, 2);
		INSERT INTO user_roles values(null, 7, 3);
		INSERT INTO user_roles values(null, 8, 2);
		INSERT INTO user_roles values(null, 9, 2);
		INSERT INTO user_roles values(null, 10, 1);
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

	if sql := builder.GetSql(); sql != "SELECT * FROM `users`" {
		t.Errorf("Unexpected SQL result, got: %s", sql)
	}

	builder = New(dialect, db)
	builder.Table("users").Select("id", "username", "email")

	if sql := builder.GetSql(); sql != "SELECT `id`,`username`,`email` FROM `users`" {
		t.Errorf("Unexpected SQL result, got: %s", sql)
	}

}

func TestWithWhere(t *testing.T) {
	dialect := dialect.New("?", "`", "`")
	builder = New(dialect, db)
	builder.Table("users").Select("*").Where("email", clause.OperatorEqual, "johndoe@gmail.com")

	if sql := builder.GetSql(); sql != "SELECT * FROM `users` WHERE `email` = ?" {
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

	if sql := builder.GetSql(); sql != "SELECT * FROM `users` WHERE `email` = ? AND `access_role` < ?" {
		t.Errorf("Unexpected SQL result, got: %s", sql)
	}
}

func TestWhereIn(t *testing.T) {
	dialect := dialect.New("?", "`", "`")
	builder = New(dialect, db)
	builder.Table("users")

	builder.Select("*").
		WhereIn("email", []any{"johndoe@example.com", "amal@example.com"})

	if sql := builder.GetSql(); sql != "SELECT * FROM `users` WHERE `email` IN(?,?)" {
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

	if sql := builder.GetSql(); sql != "SELECT * FROM `users` WHERE `age` BETWEEN ? AND ?" {
		t.Errorf("Unexpected SQL result, got: %s", sql)
	}

	builder = builder.Table("users")
	builder.
		Select("*").
		WhereBetween("dob", "1995-02-01", "2000-01-01")

	if sql := builder.GetSql(); sql != "SELECT * FROM `users` WHERE `dob` BETWEEN ? AND ?" {
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

	if sql := builder.GetSql(); sql != "SELECT * FROM `users` WHERE `age` >= ? OR `email` = ?" {
		t.Errorf("Unexpected SQL result, got: %s", sql)
	}
}

func TestJoin(t *testing.T) {
	dialect := dialect.New("?", "`", "`")
	builder = New(dialect, db)
	builder.Table("users").Select("*")

	builder.
		Join("roles", "users.id", clause.OperatorEqual, "roles.user_id")

	if sql := builder.GetSql(); sql != "SELECT * FROM `users` INNER JOIN `roles` ON `users`.`id` = `roles`.`user_id`" {
		t.Errorf("Unexpected SQL result, got: %s", sql)
	}
}

func TestLeftJoin(t *testing.T) {
	dialect := dialect.New("?", "`", "`")
	builder = New(dialect, db)
	builder.Table("users").Select("*")

	builder.
		LeftJoin("roles", "users.id", clause.OperatorEqual, "roles.user_id")
	if sql := builder.GetSql(); sql != "SELECT * FROM `users` LEFT JOIN `roles` ON `users`.`id` = `roles`.`user_id`" {
		t.Errorf("Unexpected SQL result, got: %s", sql)
	}
}

func TestRightJoin(t *testing.T) {
	dialect := dialect.New("?", "`", "`")
	builder = New(dialect, db)
	builder.Table("users").Select("*")

	builder.
		RightJoin("roles", "users.id", clause.OperatorEqual, "roles.user_id")
	if sql := builder.GetSql(); sql != "SELECT * FROM `users` RIGHT JOIN `roles` ON `users`.`id` = `roles`.`user_id`" {
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

	if sql := builder.GetSql(); sql != "SELECT * FROM `users` WHERE EXISTS (SELECT * FROM `roles` WHERE `users`.`id` = ?)" {
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

	if sql := builder.GetSql(); sql != "SELECT * FROM `users` WHERE `email` = (SELECT `user_id` FROM `roles` WHERE `users`.`id` = ?)" {
		t.Errorf("Unexpected SQL result, got: %s", sql)
	}
}

func TestGroupBy(t *testing.T) {
	dialect := dialect.New("?", "`", "`")
	builder = New(dialect, db)
	builder.Table("users").Select("*")

	builder.
		GroupBy("age", "role")

	if sql := builder.GetSql(); sql != "SELECT * FROM `users` GROUP BY `age`,`role`" {
		t.Errorf("Unexpected SQL result, got: %s", sql)
	}
}

func TestAdvancedQueryWithMultipleClauses(t *testing.T) {
	dialect := dialect.New("?", "`", "`")
	builder = New(dialect, db)
	builder.Table("users").Select("*").
		WhereBetween("age", 18, 30).
		WhereIn("email", []any{""}).
		WhereDate("created_at", clause.OperatorGreaterThan, "2023-01-01").
		LockForUpdate()

	expected := "SELECT * FROM `users` WHERE `age` BETWEEN ? AND ? AND `email` IN(?) AND DATE(`created_at`) > ? FOR UPDATE"

	if sql := builder.GetSql(); sql != expected {
		t.Errorf("Unexpected SQL result, got: %s", sql)
	}
}

func TestForUpdate(t *testing.T) {
	dialect := dialect.New("?", "`", "`")
	builder = New(dialect, db)
	builder.Table("users").Select("*").LockForUpdate()

	if sql := builder.GetSql(); sql != "SELECT * FROM `users` FOR UPDATE" {
		t.Errorf("Unexpected SQL result, got: %s", sql)
	}
}

func TestForShare(t *testing.T) {
	dialect := dialect.New("?", "`", "`")
	builder = New(dialect, db)
	builder.Table("users").Select("*").LockForShare()

	if sql := builder.GetSql(); sql != "SELECT * FROM `users` FOR SHARE" {
		t.Errorf("Unexpected SQL result, got: %s", sql)
	}
}

func TestForUpdateQuery(t *testing.T) {
	dialect := dialect.New("?", "`", "`")
	builder = New(dialect, db)
	builder.Table("users").Select("*").Where("age", clause.OperatorGreaterThan, 30).LockForUpdate()

	if sql := builder.GetSql(); sql != "SELECT * FROM `users` WHERE `age` > ? FOR UPDATE" {
		t.Errorf("Unexpected SQL result, got: %s", sql)
	}
}

func TestForShareQuery(t *testing.T) {
	dialect := dialect.New("?", "`", "`")
	builder = New(dialect, db)
	builder.Table("users").Select("*").Where("age", clause.OperatorGreaterThan, 30).LockForShare()

	if sql := builder.GetSql(); sql != "SELECT * FROM `users` WHERE `age` > ? FOR SHARE" {
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

	if len(users) != 10 {
		t.Errorf("Expected return %d of users, but got %d", 10, len(users))
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
		stmt := builder.GetSql()

		t.Errorf("SQL: %s", stmt)
		t.Errorf("Arguments: %v", arguments)
		t.Error(err)
	}

	if user.Email != "johndoe@example.com" {
		t.Errorf("Expected user email is johndoe@example.com, but got: %s", user.Email)
	}
}

func TestExecuteWhereStatementWithJoin(t *testing.T) {

	type UserWithRole struct {
		Id       int    `db:"id"`
		Username string `db:"username"`
		Email    string `db:"email"`
		Age      int    `db:"age"`
		RoleName string `db:"name"`
	}

	user := new(UserWithRole)
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
	err = builder.Table("users").
		Select("users.id", "users.username", "users.email", "users.age", "roles.name").
		Join("user_roles", "users.id", "=", "user_roles.user_id").
		Join("roles", "user_roles.role_id", "=", "roles.id").
		Where("roles.id", "=", 1).
		Get(user)

	if err != nil {
		arguments := builder.GetArguments()
		stmt := builder.GetSql()

		t.Errorf("SQL: %s", stmt)
		t.Errorf("Arguments: %v", arguments)
		t.Error(err)
	}

	if user.Email != "johndoe@example.com" {
		t.Errorf("Expected user email is johndoe@example.com, but got: %s", user.Email)
	}

	if user.RoleName != "admin" {
		t.Errorf("Expected user role is admin, but got: %s", user.RoleName)
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
		stmt := builder.GetSql()

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

func TestStatementMustBeDifferentForSameBuilderInstance(t *testing.T) {
	dialect := dialect.New("?", "`", "`")
	builder = New(dialect, db)
	builder.Table("users").Select("*").Where("email", clause.OperatorEqual, "test@xample.com")

	if sql := builder.GetSql(); sql != "SELECT * FROM `users` WHERE `email` = ?" {
		t.Fatalf("Unexpected SQL result, got: %s", sql)
	}

	builder.Table("users").Select("id", "username").
		Where("age", clause.OperatorGreaterThan, 30).
		Where("email", clause.OperatorEqual, "test@example.com").
		Limit(1)

	if sql := builder.GetSql(); sql != "SELECT `id`,`username` FROM `users` WHERE `age` > ? AND `email` = ? LIMIT ?" {
		t.Fatalf("Unexpected SQL result, got: %s", sql)
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

	if err != nil {
		t.Fatal(err)
	}

	dialect := dialect.New("?", "`", "`")
	builder = New(dialect, dba)

	result, err := builder.Table("users").Where("id", clause.OperatorEqual, 1).Update(map[string]any{
		"username": "john_doe_updated",
		"age":      36,
	})

	if sql := builder.GetSql(); err != nil {
		t.Errorf("Unexpected SQL result, got: %s", sql)
		t.Fatal(err)
	}

	if rowsAffected, err := result.RowsAffected(); err != nil {
		t.Error(err)

		if rowsAffected <= 0 {
			t.Errorf("Expected rows affected to be greater than 0, but got: %d", rowsAffected)
		}

		if rowsAffected > 1 {
			t.Errorf("Expected rows affected to be 1, but got: %d", rowsAffected)
		}
	}

	var user User
	err = builder.Table("users").Select("*").Where("id", clause.OperatorEqual, 1).Get(&user)
	if err != nil {
		t.Errorf("Unexpected SQL Result, got: %s", builder.GetSql())
		t.Fatal(err)
	}

	if user.Username != "john_doe_updated" {
		t.Errorf("Expected username to be john_doe_updated, but got: %s", user.Username)
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
	result, err := builder.Table("users").Where("username", clause.OperatorEqual, "johndoe").Delete()

	if err != nil {
		t.Errorf("SQL: %s", builder.GetSql())
		t.Fatal(err)
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

		if err = b.Table("users").Select("*").Where("id", clause.OperatorEqual, lastInsertId).Get(&newUser); err != nil {
			return errors.New("failed to retrieve user after insert: " + err.Error())
		}

		type UpdateRequest struct {
			Age int `db:"age"`
		}
		update := UpdateRequest{Age: 40}
		updateMap := map[string]any{}
		err = toMap(update, &updateMap)

		if err != nil {
			return errors.New("failed to convert update struct to map: " + err.Error())
		}

		if _, err = b.Table("users").Where("id", clause.OperatorEqual, lastInsertId).Update(updateMap); err != nil {
			return errors.New("failed to update user: " + err.Error())
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
		sql := builder.GetSql()
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
		t.Fatal(err)
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

	if len(users) != 6 {
		t.Errorf("Expected return %d of users, but got %d", 6, len(users))
	}

}

func TestRawStatement(t *testing.T) {
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

	err = builder.Raw("SELECT * FROM users WHERE age > ?", 30).Get(&users)

	if err != nil {
		t.Error(err)
	}

	if len(users) != 4 {
		t.Errorf("Expected return %d of users, but got %d", 4, len(users))
	}
}

func TestRetrieveUserCreatedOnSpecificDate(t *testing.T) {
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

	err = builder.Table("users").Select("*").WhereDate("created_at", clause.OperatorEqual, "2023-01-01").Get(&users)

	if err != nil {
		sql := builder.GetSql()
		t.Error(err)
		t.Errorf("Query: %s", sql)
	}

	if len(users) != 1 {
		t.Errorf("Expected return %d of users, but got %d", 1, len(users))
	}
}

func TestRetrieveUserCreatedBeforeSpecificDate(t *testing.T) {
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

	err = builder.Table("users").Select("*").WhereDate("created_at", clause.OperatorLessThan, "2023-01-05").Get(&users)

	if err != nil {
		sql := builder.GetSql()
		t.Error(err)
		t.Errorf("Query: %s", sql)
	}

	if len(users) != 4 {
		t.Errorf("Expected return %d of users, but got %d", 4, len(users))
	}
}

func TestRetrieveUserCreatedAfterSpecificDate(t *testing.T) {
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

	err = builder.Table("users").Select("*").WhereDate("created_at", clause.OperatorGreaterThan, "2023-01-05").Get(&users)

	if err != nil {
		sql := builder.GetSql()
		t.Error(err)
		t.Errorf("Query: %s", sql)
	}

	if len(users) != 5 {
		t.Errorf("Expected return %d of users, but got %d", 5, len(users))
	}
}

func TestRetrieveUserCreatedBetweenSpecificDates(t *testing.T) {
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

	err = builder.Table("users").Select("*").WhereDate("created_at", clause.OperatorGreatherThanEqual, "2023-01-03").WhereDate("created_at", clause.OperatorLessThanEqual, "2023-01-07").Get(&users)

	if err != nil {
		sql := builder.GetSql()
		arguments := builder.GetArguments()
		t.Error(err)
		t.Errorf("Query: %s", sql)
		t.Errorf("Arguments: %v", arguments)
	}

	if len(users) != 5 {
		t.Errorf("Expected return %d of users, but got %d", 5, len(users))
	}
}

func TestRetrieveUserCreatedOnSpecificYear(t *testing.T) {
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

	err = builder.Table("users").Select("*").WhereYear("created_at", clause.OperatorEqual, 2023).Get(&users)

	if err != nil {
		sql := builder.GetSql()
		t.Error(err)
		t.Errorf("Query: %s", sql)
	}

	if len(users) != 10 {
		sql := builder.GetSql()
		t.Errorf("Expected return %d of users, but got %d", 10, len(users))
		t.Errorf("Query: %s", sql)
	}

	var users2024 []User
	err = builder.Table("users").Select("*").WhereYear("created_at", clause.OperatorEqual, 2024).Get(&users2024)

	if err != nil {
		sql := builder.GetSql()
		t.Error(err)
		t.Errorf("Query: %s", sql)
	}

	if len(users2024) != 0 {
		t.Errorf("Expected return %d of users, but got %d", 0, len(users2024))
	}
}
