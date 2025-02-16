package sqlbuilder

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

const (
	Eq = "="
	Lte = "<="
	Lt = "<"
	Gte = ">="
	Gt= ">"
)

type SQLBuilder struct {
	dialect string
	sql *sql.DB
	Statement string
}

type Builder interface {
	Table(table string) Builder
	Select(cols []string) Builder
	Where(column string, op string, val interface{}) Builder
	WhereIn(column string, d interface{}) Builder
	WhereBetween(column string, start interface{}, end interface{}) Builder
	GetSql() string
}

func (b *SQLBuilder) Table(table string) Builder {
	b.Statement += " FROM " + table
	return b
}

func (b *SQLBuilder) Select(cols []string) Builder {
	b.Statement += "SELECT " + strings.Join(cols, ", ")
	return b
}

func (b *SQLBuilder) GetSql() string {
	return b.Statement
}

func (b *SQLBuilder) Where(column string, op string, val interface{}) Builder {
	clause := "WHERE"
	valueBind := "%v"

	if hasWhere(b.Statement) {
		clause = "AND"
	}

	t := reflect.ValueOf(val)
	if t.Kind() == reflect.String {
		valueBind = "'%v'"
	}

	b.Statement += fmt.Sprintf(" %s %s %s " + valueBind, clause, column, op, val)
	return b
}

func (b *SQLBuilder) WhereIn(column string, values interface{}) Builder {
	clause := "WHERE"
	var inValues string

	if hasWhere(b.Statement) {
		clause = "AND"
	}

	ref := reflect.ValueOf(values)
	if ref.Kind() != reflect.Slice {
		return b
	}
	lenValues := ref.Len()

	for i := 0; i < lenValues; i++ {
		vRef := ref.Index(i)
		val := ref.Index(i).Interface()

		con := fmt.Sprintf("%v", val)

		if vRef.Kind() == reflect.String {
			con = fmt.Sprintf("'%s'", val)
		}

		if i < (lenValues - 1) {
			con += ", "
		}

		inValues += con
	}


	b.Statement += fmt.Sprintf(" %s %s IN(%s)", clause, column, inValues)
	return b
}

func (b *SQLBuilder) WhereBetween(column string, start interface{}, end interface{}) Builder {
	clause := "WHERE"

	if hasWhere(b.Statement) {
		clause = "AND"
	}

	b.Statement += fmt.Sprintf(" %s %s BETWEEN %v AND %v", clause, column, start, end)
	return b
}

func NewSQLBuilder(dialect string, sql *sql.DB) *SQLBuilder {
	return &SQLBuilder{dialect: dialect, sql: sql}
}

func hasWhere(statement string) bool {
	return strings.Contains(statement, "WHERE")
}
