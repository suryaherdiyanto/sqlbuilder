# SQLBuilder

A lightweight SQL builder for Go with pluggable dialects.

## Installation

```bash
go get github.com/suryaherdiyanto/sqlbuilder
```

## Usage Example

```go
package main

import (
    "fmt"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/suryaherdiyanto/sqlbuilder"
	"github.com/suryaherdiyanto/sqlbuilder/clause"
	"github.com/suryaherdiyanto/sqlbuilder/dialect"
)

func main() {
	db, _ := sql.Open("sqlite3", ":memory:")

	b := sqlbuilder.New(dialect.New("?", "`", "`"), db)

    type User struct {
        ID uint64
        name string
        username string
    }

    user := User{}
	err := b.
		Table("users").
		Select("id", "email").
		Where("id", "=", 1).
		Get(&user)

    fmt.Println(user)
}
```

For PostgreSQL, use `dialect.NewPostgres()` which emits placeholders like $1, $2, ... and double-quoted identifiers.

## Insert Example

```go
package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/suryaherdiyanto/sqlbuilder"
	"github.com/suryaherdiyanto/sqlbuilder/dialect"
)

func main() {
	db, _ := sql.Open("sqlite3", ":memory:")
	b := sqlbuilder.New(dialect.New("?", "`", "`"), db)

	id, err := b.Table("users").Insert(map[string]any{
		"username": "alice",
		"email":    "alice@example.com",
		"age":      29,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("last insert id:", id)

	// Insert multiple rows.
	_, err = b.Table("users").InsertMany([]map[string]any{
		{"username": "bob", "email": "bob@example.com", "age": 31},
		{"username": "carol", "email": "carol@example.com", "age": 27},
	}).Exec()
	if err != nil {
		panic(err)
	}
}
```

## Update Example

```go
package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/suryaherdiyanto/sqlbuilder"
	"github.com/suryaherdiyanto/sqlbuilder/clause"
	"github.com/suryaherdiyanto/sqlbuilder/dialect"
)

func main() {
	db, _ := sql.Open("sqlite3", ":memory:")
	b := sqlbuilder.New(dialect.New("?", "`", "`"), db)

	result, err := b.
		Table("users").
		Where("id", clause.OperatorEqual, 1).
		Update(map[string]any{
			"username": "alice_updated",
			"age":      30,
		})
	if err != nil {
		panic(err)
	}

	rowsAffected, _ := result.RowsAffected()
	fmt.Println("updated rows:", rowsAffected)
}
```

## Delete Example

```go
package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/suryaherdiyanto/sqlbuilder"
	"github.com/suryaherdiyanto/sqlbuilder/clause"
	"github.com/suryaherdiyanto/sqlbuilder/dialect"
)

func main() {
	db, _ := sql.Open("sqlite3", ":memory:")
	b := sqlbuilder.New(dialect.New("?", "`", "`"), db)

	result, err := b.
		Table("users").
		Where("username", clause.OperatorEqual, "alice_updated").
		Delete()
	if err != nil {
		panic(err)
	}

	rowsAffected, _ := result.RowsAffected()
	fmt.Println("deleted rows:", rowsAffected)
}
```
