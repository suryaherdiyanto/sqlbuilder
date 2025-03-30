package sqlbuilder

import (
	"context"
	"database/sql"
	"fmt"
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
	Dialect   string
	sql       *sql.DB
	hasWhere  bool
	statement *Statement
}

type Statement struct {
	SQL       string
	arguments []interface{}
}

type Builder interface {
	Table(table string, columns ...string) *SQLBuilder
	Where(statement string, vars ...interface{}) *SQLBuilder
	WhereFunc(statement string, b func(b Builder) *SQLBuilder) *SQLBuilder
	Join(table string, first string, operator string, second string)
	LeftJoin(table string, first string, operator string, second string)
	RightJoin(table string, first string, operator string, second string)
	OrderBy(column string, dir string) *SQLBuilder
	GroupBy(columns ...string) *SQLBuilder
	Limit(n int) *SQLBuilder
	Offset(n int) *SQLBuilder
}

func New(dialect string, db *sql.DB) *SQLBuilder {
	return &SQLBuilder{
		Dialect: dialect,
		sql:     db,
	}
}

func (s *SQLBuilder) NewSelect() *SQLBuilder {
	statement := &Statement{
		SQL: "SELECT ",
	}
	s.statement = statement
	return s
}

func (s *SQLBuilder) Table(table string, columns ...string) *SQLBuilder {
	s.statement.SQL += fmt.Sprintf("%s FROM %s", strings.Join(columns, ","), table)
	return s
}

func (s *SQLBuilder) GetSql() string {
	return s.statement.SQL
}

func (s *SQLBuilder) GetArguments() []interface{} {
	return s.statement.arguments
}

func (s *SQLBuilder) Where(statement string, vars ...interface{}) *SQLBuilder {
	s.statement.SQL += fmt.Sprintf(" WHERE %s", statement)
	s.statement.arguments = append(s.statement.arguments, vars...)
	return s
}

func (s *SQLBuilder) Join(table string, first string, operator string, second string) {
	s.statement.SQL += fmt.Sprintf(" INNER JOIN %s ON %s %s %s", table, first, operator, second)
}

func (s *SQLBuilder) LeftJoin(table string, first string, operator string, second string) {
	s.statement.SQL += fmt.Sprintf(" LEFT JOIN %s ON %s %s %s", table, first, operator, second)
}

func (s *SQLBuilder) RightJoin(table string, first string, operator string, second string) {
	s.statement.SQL += fmt.Sprintf(" RIGHT JOIN %s ON %s %s %s", table, first, operator, second)
}

func (s *SQLBuilder) WhereFunc(statement string, builder func(b Builder) *SQLBuilder) *SQLBuilder {
	newBuilder := builder(New(s.Dialect, s.sql).NewSelect())

	s.statement.SQL += fmt.Sprintf(" WHERE %s", statement)
	s.statement.SQL += fmt.Sprintf("(%s)", newBuilder.GetSql())
	s.statement.arguments = append(s.statement.arguments, newBuilder.GetArguments()...)
	return s
}

func (s *SQLBuilder) WhereExists(builder func(b Builder) *SQLBuilder) *SQLBuilder {
	newBuilder := builder(New(s.Dialect, s.sql).NewSelect())
	s.statement.SQL += fmt.Sprintf(" WHERE EXISTS (%s)", newBuilder.GetSql())
	s.statement.arguments = append(s.statement.arguments, newBuilder.GetArguments()...)

	return s
}

func (s *SQLBuilder) OrderBy(column string, dir string) *SQLBuilder {
	s.statement.SQL += fmt.Sprintf(" ORDER BY %s %s", column, dir)
	return s
}
func (s *SQLBuilder) GroupBy(columns ...string) *SQLBuilder {
	s.statement.SQL += fmt.Sprintf(" GROUP BY %s", strings.Join(columns, ", "))
	return s
}
func (s *SQLBuilder) Limit(n int) *SQLBuilder {
	s.statement.SQL += fmt.Sprintf(" LIMIT %d", n)
	return s
}
func (s *SQLBuilder) Offset(n int) *SQLBuilder {
	s.statement.SQL += fmt.Sprintf(" OFFSET %d", n)
	return s
}

func (s *SQLBuilder) Find(d interface{}, ctx context.Context) error {
	rows, err := s.runQuery(ctx)
	defer rows.Close()

	if err != nil {
		return err
	}

	if rows.Next() {
		return ScanStruct(d, rows)
	}

	return nil

}
func (b *SQLBuilder) Get(d interface{}, ctx context.Context) error {
	rows, err := b.runQuery(ctx)
	defer rows.Close()

	if err != nil {
		return err
	}

	if err = ScanAll(d, rows); err != nil {
		return err
	}

	return nil
}
func (b *SQLBuilder) Scan(d interface{}) error {
	rows, err := b.runQuery(context.Background())
	defer rows.Close()

	if err != nil {
		return err
	}

	if rows.Next() {
		return rows.Scan(d)
	}

	return nil
}
func (s *SQLBuilder) Exec(ctx context.Context) error {
	_, err := s.sql.ExecContext(ctx, s.statement.SQL, s.statement.arguments...)

	return err
}

func (s *SQLBuilder) runQuery(ctx context.Context) (*sql.Rows, error) {
	return s.sql.QueryContext(ctx, s.statement.SQL, s.statement.arguments...)
}
