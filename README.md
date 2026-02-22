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
