package sqlbuilder

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

const (
	pgsqlPlaceholder  = "$%d"
	mysqlPlaceholder  = "?"
	sqlitePlaceholder = "?"
)

type SQLBuilder struct {
	Dialect   string
	sql       *sql.DB
	statement *Statement
}

type Statement struct {
	SQL          string
	Table        string
	Command      string
	Columns      []string
	Where        string
	Joins        []string
	arguments    []interface{}
	values       string
	subquery     string
	setStatement string
}

type Builder interface {
	Table(table string, columns ...string) *SQLBuilder
	Where(statement string, vars ...interface{}) *SQLBuilder
	WhereFunc(statement string, b func(b Builder) *SQLBuilder) *SQLBuilder
	Join(table string, first string, operator string, second string) *SQLBuilder
	LeftJoin(table string, first string, operator string, second string) *SQLBuilder
	RightJoin(table string, first string, operator string, second string) *SQLBuilder
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
		SQL:     "SELECT ",
		Command: "SELECT",
	}
	s.statement = statement
	return s
}

func (s *SQLBuilder) Insert(data interface{}) (sql.Result, error) {

	stmt, err := s.buildInsert(data)
	if err != nil {
		return nil, err
	}

	insertData, err := s.extractData(data)

	if err != nil {
		return nil, err
	}

	s.statement.Columns = stmt["columns"]
	s.statement.arguments = insertData
	s.statement.values = strings.Join(stmt["values"], ",")
	s.statement.Command = "INSERT"

	ctx := context.Background()
	return s.Exec(ctx)
}

func (s *SQLBuilder) Update(data interface{}) (sql.Result, error) {
	stmt, err := s.buildUpdate(data)
	if err != nil {
		return nil, err
	}

	args, err := s.extractData(data)
	if err != nil {
		return nil, err
	}

	s.statement.arguments = append(s.statement.arguments, args...)
	s.statement.setStatement = stmt
	s.statement.Command = "UPDATE"
	ctx := context.Background()

	return s.Exec(ctx)
}

func (s *SQLBuilder) Delete() *SQLBuilder {
	s.statement.Command = "DELETE"
	return nil
}

func (s *SQLBuilder) Table(table string, columns ...string) *SQLBuilder {
	vRef := reflect.ValueOf(s.statement)
	if vRef.IsZero() {
		s.statement = &Statement{}
	}

	s.statement.Table = table
	s.statement.Columns = columns

	return s
}

func (s *SQLBuilder) GetSql() string {
	switch s.statement.Command {
	case "SELECT":
		stmt := fmt.Sprintf("SELECT %s FROM %s", strings.Join(s.statement.Columns, ","), s.statement.Table)

		if len(s.statement.Joins) > 0 {
			for _, join := range s.statement.Joins {
				stmt += fmt.Sprintf("%s", join)
			}
		}

		if s.statement.Where != "" {
			stmt += fmt.Sprintf(" WHERE %s", s.statement.Where)

			if s.statement.subquery != "" {
				stmt += fmt.Sprintf(" (%s)", s.statement.subquery)
			}
		}

		s.statement.SQL = stmt
	case "INSERT":
		stmt := fmt.Sprintf("INSERT INTO %s(%s) VALUES(%s)", s.statement.Table, strings.Join(s.statement.Columns, ","), s.statement.values)
		s.statement.SQL = stmt
	case "UPDATE":
		stmt := fmt.Sprintf("UPDATE %s %s", s.statement.Table, s.statement.setStatement)
		if s.statement.Where != "" {
			stmt += fmt.Sprintf(" WHERE %s", s.statement.Where)
		}
		s.statement.SQL = stmt
	case "DELETE":
		stmt := fmt.Sprintf("DELETE FROM %s", s.statement.Table)
		if s.statement.Where != "" {
			stmt += fmt.Sprintf(" WHERE %s", s.statement.Where)
		}
		s.statement.SQL = stmt
	}

	return s.statement.SQL
}

func (s *SQLBuilder) GetArguments() []interface{} {
	return s.statement.arguments
}

func (s *SQLBuilder) Where(statement string, vars ...interface{}) *SQLBuilder {
	s.statement.SQL += fmt.Sprintf(" WHERE %s", statement)
	s.statement.arguments = append(s.statement.arguments, vars...)
	s.statement.Where = statement
	return s
}

func (s *SQLBuilder) Join(table string, first string, operator string, second string) *SQLBuilder {
	s.statement.Joins = append(s.statement.Joins, fmt.Sprintf(" INNER JOIN %s ON %s %s %s", table, first, operator, second))
	return s
}

func (s *SQLBuilder) LeftJoin(table string, first string, operator string, second string) *SQLBuilder {
	s.statement.Joins = append(s.statement.Joins, fmt.Sprintf(" LEFT JOIN %s ON %s %s %s", table, first, operator, second))
	return s
}

func (s *SQLBuilder) RightJoin(table string, first string, operator string, second string) *SQLBuilder {
	s.statement.Joins = append(s.statement.Joins, fmt.Sprintf(" RIGHT JOIN %s ON %s %s %s", table, first, operator, second))
	return s
}

func (s *SQLBuilder) WhereFunc(statement string, builder func(b Builder) *SQLBuilder) *SQLBuilder {
	newBuilder := builder(New(s.Dialect, s.sql).NewSelect())

	s.statement.Where = fmt.Sprintf("%s (%s)", statement, newBuilder.GetSql())
	s.statement.arguments = append(s.statement.arguments, newBuilder.GetArguments()...)
	return s
}

func (s *SQLBuilder) WhereExists(builder func(b Builder) *SQLBuilder) *SQLBuilder {
	newBuilder := builder(New(s.Dialect, s.sql).NewSelect())

	s.statement.Where = fmt.Sprintf("EXISTS (%s)", newBuilder.GetSql())
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
func (s *SQLBuilder) Exec(ctx context.Context) (sql.Result, error) {
	return s.sql.ExecContext(ctx, s.GetSql(), s.statement.arguments...)
}

func (s *SQLBuilder) runQuery(ctx context.Context) (*sql.Rows, error) {
	return s.sql.QueryContext(ctx, s.GetSql(), s.statement.arguments...)
}

func (s *SQLBuilder) extractData(data interface{}) ([]interface{}, error) {
	valRef := reflect.ValueOf(data)
	if valRef.Kind() == reflect.Ptr {
		valRef = valRef.Elem()
	}

	if valRef.Kind() != reflect.Struct {
		return nil, errors.New("data must be a struct")
	}

	var result []interface{}
	for i := 0; i < valRef.NumField(); i++ {
		field := valRef.Field(i)
		result = append(result, field.Interface())
	}

	return result, nil
}

func (s *SQLBuilder) buildInsert(data interface{}) (map[string][]string, error) {
	valRef := reflect.TypeOf(data)
	if valRef.Kind() == reflect.Ptr {
		valRef = valRef.Elem()
	}

	if valRef.Kind() != reflect.Struct {
		return make(map[string][]string), errors.New("data must be a struct")
	}

	var columns []string
	var values []string
	placeholder := func() string {
		switch s.Dialect {
		case "pgsql":
			return fmt.Sprintf(pgsqlPlaceholder, len(values)+1)
		case "mysql":
			return mysqlPlaceholder
		case "sqlite":
			return sqlitePlaceholder
		default:
			return mysqlPlaceholder
		}
	}

	for i := 0; i < valRef.NumField(); i++ {
		field := valRef.Field(i)
		columns = append(columns, field.Tag.Get("db"))
		values = append(values, placeholder())
	}

	return map[string][]string{
		"columns": columns,
		"values":  values,
	}, nil
}

func (s *SQLBuilder) buildUpdate(data interface{}) (string, error) {
	valRef := reflect.TypeOf(data)
	if valRef.Kind() == reflect.Ptr {
		valRef = valRef.Elem()
	}

	if valRef.Kind() != reflect.Struct {
		return "", errors.New("data must be a struct")
	}

	var set []string
	placeholder := func() string {
		switch s.Dialect {
		case "pgsql":
			return fmt.Sprintf(pgsqlPlaceholder, len(set)+1)
		case "mysql":
			return mysqlPlaceholder
		case "sqlite":
			return sqlitePlaceholder
		default:
			return mysqlPlaceholder
		}
	}

	for i := 0; i < valRef.NumField(); i++ {
		field := valRef.Field(i)
		set = append(set, fmt.Sprintf("%s = %s", field.Tag.Get("db"), placeholder()))
	}

	setStatement := fmt.Sprintf("%s", strings.Join(set, ", "))

	return fmt.Sprintf("SET %s", setStatement), nil
}
