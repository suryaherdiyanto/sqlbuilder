package sqlbuilder

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

const (
	Eq   = "="
	Lte  = "<="
	Lt   = "<"
	Gte  = ">="
	Gt   = ">"
	Desc = "DESC"
	Asc  = "ASC"
)

const (
	whereOr  = "OR"
	whereAnd = "AND"
)

type SQLBuilder struct {
	dialect   string
	sql       *sql.DB
	hasWhere  bool
	arguments []interface{}
	Statement string
}

type Builder interface {
	Table(table string) *SQLBuilder
	Select(cols ...string) *SQLBuilder
	Where(column string, comp string, val interface{}) *SQLBuilder
	WhereIn(column string, d interface{}) *SQLBuilder
	WhereBetween(column string, start interface{}, end interface{}) *SQLBuilder
	OrWhere(column string, comp string, val interface{}) *SQLBuilder
	OrWhereIn(column string, d interface{}) *SQLBuilder
	OrWhereBetween(column string, start interface{}, end interface{}) *SQLBuilder
	OrderBy(column string, dir string) *SQLBuilder
	GroupBy(columns ...string) *SQLBuilder
	GetSql() string
}

func (b *SQLBuilder) Table(table string) *SQLBuilder {
	b.Statement += " FROM " + table
	return b
}

func (b *SQLBuilder) Select(cols ...string) *SQLBuilder {
	b.Statement += "SELECT " + strings.Join(cols, ", ")
	return b
}

func (b *SQLBuilder) GetSql() string {
	return b.Statement
}

func (b *SQLBuilder) buildWhere(column string, comp string, val interface{}) *SQLBuilder {
	b.Statement += fmt.Sprintf(" %s %s ?", column, comp)
	b.arguments = append(b.arguments, val)

	return b
}

func (b *SQLBuilder) buildWhereIn(column string, values interface{}) *SQLBuilder {
	var inValues string

	ref := reflect.ValueOf(values)
	if ref.Kind() != reflect.Slice {
		return b
	}
	lenValues := ref.Len()

	for i := 0; i < lenValues; i++ {
		val := ref.Index(i).Interface()

		con := "?"

		if i < (lenValues - 1) {
			con += ", "
		}

		inValues += con
		b.arguments = append(b.arguments, val)
	}

	b.Statement += fmt.Sprintf(" %s IN(%s)", column, inValues)
	return b
}

func (b *SQLBuilder) buildWhereBetween(column string, start interface{}, end interface{}) *SQLBuilder {
	startRef := reflect.ValueOf(start)
	endRef := reflect.ValueOf(end)

	if startRef.Kind() == reflect.String {
		start = fmt.Sprintf("'%s'", startRef.Interface())
	}

	if endRef.Kind() == reflect.String {
		end = fmt.Sprintf("'%s'", endRef.Interface())
	}

	b.Statement += fmt.Sprintf(" %s BETWEEN ? AND ?", column)
	b.arguments = append(b.arguments, start)
	b.arguments = append(b.arguments, end)
	return b
}

func (b *SQLBuilder) Where(column string, comp string, val interface{}) *SQLBuilder {
	b.setWhereOperator(whereAnd)
	return b.buildWhere(column, comp, val)
}

func (b *SQLBuilder) OrWhere(column string, comp string, val interface{}) *SQLBuilder {
	b.setWhereOperator(whereOr)
	return b.buildWhere(column, comp, val)
}

func (b *SQLBuilder) WhereIn(column string, values interface{}) *SQLBuilder {
	b.setWhereOperator(whereAnd)
	return b.buildWhereIn(column, values)
}
func (b *SQLBuilder) OrWhereIn(column string, values interface{}) *SQLBuilder {
	b.setWhereOperator(whereOr)
	return b.buildWhereIn(column, values)
}

func (b *SQLBuilder) WhereBetween(column string, start interface{}, end interface{}) *SQLBuilder {
	b.setWhereOperator(whereAnd)
	return b.buildWhereBetween(column, start, end)
}
func (b *SQLBuilder) OrWhereBetween(column string, start interface{}, end interface{}) *SQLBuilder {
	b.setWhereOperator(whereOr)
	return b.buildWhereBetween(column, start, end)
}
func (b *SQLBuilder) OrderBy(column string, dir string) *SQLBuilder {
	b.Statement += fmt.Sprintf(" ORDER BY %s %s", column, dir)
	return b
}
func (b *SQLBuilder) GroupBy(column ...string) *SQLBuilder {
	b.Statement += fmt.Sprintf(" GROUP BY %s", strings.Join(column, ", "))
	return b
}
func (b *SQLBuilder) Scan(d interface{}, ctx context.Context) error {
	rows, err := b.runQuery(ctx)
	defer rows.Close()

	if err != nil {
		return err
	}

	if rows.Next() {
		return ScanStruct(d, rows)
	}

	return nil

}
func (b *SQLBuilder) ScanAll(d interface{}, ctx context.Context) error {
	rows, err := b.runQuery(ctx)
	defer rows.Close()

	ScanAll(d, rows)
	return err
}
func (b *SQLBuilder) Exec(ctx context.Context) error {
	_, err := b.sql.ExecContext(ctx, b.Statement, b.arguments...)
	return err
}

func (b *SQLBuilder) setWhereOperator(op string) {
	if !b.HasWhere() {
		b.Statement += " WHERE"
	}

	if b.HasWhere() {
		b.Statement += fmt.Sprintf(" %s", op)
	}

	if strings.Contains(b.Statement, "WHERE") {
		b.hasWhere = true
	}
}

func (b *SQLBuilder) runQuery(ctx context.Context) (*sql.Rows, error) {
	return b.sql.QueryContext(ctx, b.Statement, b.arguments...)
}

func NewSQLBuilder(dialect string, sql *sql.DB) *SQLBuilder {
	return &SQLBuilder{dialect: dialect, sql: sql}
}

func (b *SQLBuilder) HasWhere() bool {
	return b.hasWhere
}
