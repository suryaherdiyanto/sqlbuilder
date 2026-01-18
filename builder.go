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
	Dialect         string
	sql             *sql.DB
	tx              *sql.Tx
	isTx            bool
	tempTable       string
	statement       SelectStatement
	insertStatement InsertStatement
	updateStatement UpdateStatement
	deleteStatement DeleteStatement
}

type Builder interface {
	Select(columns ...string) *SQLBuilder
	Table(table string) *SQLBuilder
	Where(field string, Op Operator, val any) *SQLBuilder
	WhereOr(field string, Op Operator, val any) *SQLBuilder
	WhereIn(field string, values []any) *SQLBuilder
	WhereNotIn(field string, values []any) *SQLBuilder
	WhereBetween(field string, start any, end any) *SQLBuilder
	WhereFunc(field string, operator Operator, b func(b Builder) *SQLBuilder) *SQLBuilder
	Join(table string, first string, operator Operator, second string) *SQLBuilder
	LeftJoin(table string, first string, operator Operator, second string) *SQLBuilder
	RightJoin(table string, first string, operator Operator, second string) *SQLBuilder
	OrderBy(column string, dir OrderDirection) *SQLBuilder
	GroupBy(columns ...string) *SQLBuilder
	Limit(n int64) *SQLBuilder
	Offset(n int64) *SQLBuilder
}

func New(dialect string, db *sql.DB) *SQLBuilder {
	return &SQLBuilder{
		Dialect: dialect,
		sql:     db,
	}
}

func (s *SQLBuilder) Begin(tx func(s *SQLBuilder) error) error {
	transaction, err := s.sql.Begin()
	if err != nil {
		return err
	}

	builder := &SQLBuilder{
		tx:   transaction,
		isTx: true,
	}

	err = tx(builder)

	if err != nil {
		if rollbackErr := transaction.Rollback(); rollbackErr != nil {
			return rollbackErr
		}

		return err
	}

	if err = transaction.Commit(); err != nil {
		return err
	}

	return nil
}

func (b *SQLBuilder) Select(columns ...string) *SQLBuilder {
	b.statement = SelectStatement{
		Table:   b.tempTable,
		Columns: columns,
	}

	return b
}

func (s *SQLBuilder) Insert(data []map[string]any) *SQLBuilder {
	s.insertStatement = InsertStatement{
		Table: s.tempTable,
		Rows:  data,
	}

	return s
}

func (s *SQLBuilder) Update(data map[string]any) *SQLBuilder {
	s.updateStatement = UpdateStatement{
		Table: s.tempTable,
	}

	return s
}

func (s *SQLBuilder) Delete() *SQLBuilder {
	s.deleteStatement = DeleteStatement{
		Table: s.tempTable,
	}

	return s
}

func (s *SQLBuilder) Table(table string) *SQLBuilder {
	s.tempTable = table

	if !reflect.DeepEqual(s.statement, SelectStatement{}) {
		s.statement.Table = table
		return s
	}

	if !reflect.DeepEqual(s.insertStatement, InsertStatement{}) {
		s.insertStatement.Table = table
		return s
	}

	if !reflect.DeepEqual(s.updateStatement, UpdateStatement{}) {
		s.updateStatement.Table = table
		return s
	}

	if !reflect.DeepEqual(s.deleteStatement, DeleteStatement{}) {
		s.deleteStatement.Table = table
		return s
	}

	return s
}

func (s *SQLBuilder) GetSql() (string, error) {

	if !reflect.DeepEqual(s.statement, SelectStatement{}) {
		return s.statement.Parse(), nil
	}

	if !reflect.DeepEqual(s.insertStatement, InsertStatement{}) {
		return s.insertStatement.Parse(), nil
	}

	if !reflect.DeepEqual(s.updateStatement, UpdateStatement{}) {
		return s.updateStatement.Parse(), nil
	}

	if !reflect.DeepEqual(s.deleteStatement, DeleteStatement{}) {
		return s.deleteStatement.Parse(), nil
	}

	return "", errors.New("no valid statement found")
}

func (s *SQLBuilder) GetArguments() []any {

	values := []any{}
	if !reflect.DeepEqual(s.statement, SelectStatement{}) {
		values = append(values, s.statement.GetArguments()...)
	}

	if !reflect.DeepEqual(s.insertStatement, InsertStatement{}) {
		for _, row := range s.insertStatement.Rows {
			for _, val := range row {
				values = append(values, val)
			}
		}
	}

	if !reflect.DeepEqual(s.updateStatement, UpdateStatement{}) {
		values = append(values, s.updateStatement.GetArguments()...)
	}

	if !reflect.DeepEqual(s.deleteStatement, DeleteStatement{}) {
		values = append(values, s.deleteStatement.GetArguments()...)
	}

	return values
}

func (s *SQLBuilder) Where(field string, Op Operator, val any) *SQLBuilder {
	s.statement.WhereStatements.Where = append(s.statement.WhereStatements.Where, Where{
		Field: field,
		Value: val,
		Op:    Op,
		Conj:  ConjuctionAnd,
	})
	return s
}

func (s *SQLBuilder) WhereOr(field string, Op Operator, val any) *SQLBuilder {
	s.statement.WhereStatements.Where = append(s.statement.WhereStatements.Where, Where{
		Field: field,
		Value: val,
		Op:    Op,
		Conj:  ConjuctionOr,
	})
	return s
}

func (s *SQLBuilder) WhereIn(field string, values []any) *SQLBuilder {
	s.statement.WhereStatements.WhereIn = append(s.statement.WhereStatements.WhereIn, WhereIn{
		Field:  field,
		Values: values,
	})

	return s
}

func (s *SQLBuilder) WhereNotIn(field string, values []any) *SQLBuilder {
	s.statement.WhereStatements.WhereNotIn = append(s.statement.WhereStatements.WhereNotIn, WhereNotIn{
		Field:  field,
		Values: values,
	})

	return s
}

func (s *SQLBuilder) WhereBetween(field string, start any, end any) *SQLBuilder {
	s.statement.WhereStatements.WhereBetween = append(s.statement.WhereStatements.WhereBetween, WhereBetween{
		Field: field,
		Start: start,
		End:   end,
	})

	return s
}

func (s *SQLBuilder) Join(table string, first string, operator Operator, second string) *SQLBuilder {
	s.statement.JoinStatements = append(s.statement.JoinStatements, Join{
		Type:        InnerJoin,
		SecondTable: table,
		FirstTable:  s.statement.Table,
		On: JoinON{
			LeftValue:  first,
			Operator:   operator,
			RightValue: second,
		},
	})
	return s
}

func (s *SQLBuilder) LeftJoin(table string, first string, operator Operator, second string) *SQLBuilder {
	s.statement.JoinStatements = append(s.statement.JoinStatements, Join{
		Type:        LeftJoin,
		SecondTable: table,
		FirstTable:  s.statement.Table,
		On: JoinON{
			LeftValue:  first,
			Operator:   operator,
			RightValue: second,
		},
	})
	return s
}

func (s *SQLBuilder) RightJoin(table string, first string, operator Operator, second string) *SQLBuilder {
	s.statement.JoinStatements = append(s.statement.JoinStatements, Join{
		Type:        RightJoin,
		SecondTable: table,
		FirstTable:  s.statement.Table,
		On: JoinON{
			LeftValue:  first,
			Operator:   operator,
			RightValue: second,
		},
	})
	return s
}

func (s *SQLBuilder) WhereFunc(field string, operator Operator, builder func(b Builder) *SQLBuilder) *SQLBuilder {
	s.statement.WhereStatements.Where = append(s.statement.WhereStatements.Where, Where{
		Field:        field,
		Op:           operator,
		SubStatement: builder(New(s.Dialect, s.sql)).statement,
	})
	return s
}

func (s *SQLBuilder) WhereExists(builder func(b Builder) *SQLBuilder) *SQLBuilder {
	newBuilder := builder(New(s.Dialect, s.sql))

	s.statement.HasExistsClause = true
	s.statement.WhereStatements.Where = append(s.statement.WhereStatements.Where, Where{
		Op:           OperatorExists,
		SubStatement: newBuilder.statement,
	})

	return s
}

func (s *SQLBuilder) OrderBy(column string, dir OrderDirection) *SQLBuilder {
	if len(s.statement.Ordering.OrderingFields) > 0 {
		s.statement.Ordering.OrderingFields = append(s.statement.Ordering.OrderingFields, OrderField{
			Field:     column,
			Direction: OrderDirection(dir),
		})
		return s
	}

	s.statement.Ordering = Order{
		OrderingFields: []OrderField{
			{
				Field:     column,
				Direction: OrderDirection(dir),
			},
		},
	}
	return s
}
func (s *SQLBuilder) GroupBy(columns ...string) *SQLBuilder {
	if len(s.statement.GroupByStatement.Fields) > 0 {
		s.statement.GroupByStatement.Fields = append(s.statement.GroupByStatement.Fields, columns...)
	} else {
		s.statement.GroupByStatement = GroupBy{
			Fields: columns,
		}
	}
	return s
}
func (s *SQLBuilder) Limit(n int64) *SQLBuilder {
	s.statement.Limit = n
	return s
}
func (s *SQLBuilder) Offset(n int64) *SQLBuilder {
	s.statement.Offset = n
	return s
}

func (b *SQLBuilder) Get(d any, ctx context.Context) error {
	rows, err := b.runQuery(ctx)
	if err != nil {
		return err
	}

	defer rows.Close()

	ref := reflect.TypeOf(d)
	if ref.Kind() == reflect.Ptr {
		ref = ref.Elem()
	}

	if ref.Kind() == reflect.Struct {
		if rows.Next() {
			if err = ScanStruct(d, rows); err != nil {
				return err
			}
		}
	}

	if ref.Kind() == reflect.Slice {
		if err = ScanAll(d, rows); err != nil {
			return err
		}
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
	statement, err := s.GetSql()
	arguments := s.GetArguments()
	if err != nil {
		return nil, err
	}

	if s.isTx {
		// return s.tx.ExecContext(ctx, statement, s.statement.arguments...)
	}

	return s.sql.ExecContext(ctx, statement, arguments...)
}

func (s *SQLBuilder) runQuery(ctx context.Context) (*sql.Rows, error) {
	sql, err := s.GetSql()
	arguments := s.GetArguments()
	if err != nil {
		return nil, err
	}

	if s.isTx {
		// return s.tx.QueryContext(ctx, sql, s.statement.arguments...)
	}

	return s.sql.QueryContext(ctx, sql, arguments...)
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
