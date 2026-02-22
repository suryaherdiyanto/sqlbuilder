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
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/suryaherdiyanto/sqlbuilder"
	"github.com/suryaherdiyanto/sqlbuilder/clause"
	"github.com/suryaherdiyanto/sqlbuilder/dialect"
)

func main() {
	db, _ := sql.Open("sqlite3", ":memory:")

	b := sqlbuilder.New(dialect.New("?", "`", "`"), db)

	query, _ := b.
		Table("users").
		Select("id", "email").
		Where("id", clause.OperatorEqual, 1).
		GetSql()

	args := b.GetArguments()
	_, _ = query, args
}
```

For PostgreSQL, use `dialect.NewPostgres()` which emits placeholders like $1, $2, ... and double-quoted identifiers.
