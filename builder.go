package sqlbuilder

import (
	"context"
	"database/sql"
	"errors"
	"reflect"
	"slices"

	"github.com/suryaherdiyanto/sqlbuilder/clause"
)

const (
	pgsqlPlaceholder  = "$%d"
	mysqlPlaceholder  = "?"
	sqlitePlaceholder = "?"
)

type SQLBuilder struct {
	Dialect         clause.SQLDialector
	sql             *sql.DB
	tx              *sql.Tx
	isTx            bool
	tempTable       string
	statement       clause.Select
	insertStatement clause.Insert
	updateStatement clause.Update
	deleteStatement clause.Delete
	Grouping        clause.GroupBy
	Offseting       clause.Offset
	Limiting        clause.Limit
	Ordering        clause.Order
	clause.WhereStatements
}

type Builder interface {
	Select(columns ...string) *SQLBuilder
	Table(table string) *SQLBuilder
	Where(field string, Op clause.Operator, val any) *SQLBuilder
	WhereOr(field string, Op clause.Operator, val any) *SQLBuilder
	WhereIn(field string, values []any) *SQLBuilder
	WhereNotIn(field string, values []any) *SQLBuilder
	WhereBetween(field string, start any, end any) *SQLBuilder
	WhereFunc(field string, operator clause.Operator, b func(b Builder) *SQLBuilder) *SQLBuilder
	Join(table string, first string, operator clause.Operator, second string) *SQLBuilder
	LeftJoin(table string, first string, operator clause.Operator, second string) *SQLBuilder
	RightJoin(table string, first string, operator clause.Operator, second string) *SQLBuilder
	OrderBy(column string, dir clause.OrderDirection) *SQLBuilder
	GroupBy(columns ...string) *SQLBuilder
	Limit(n int64) *SQLBuilder
	Offset(n int64) *SQLBuilder
}

func New(dialect clause.SQLDialector, db *sql.DB) *SQLBuilder {
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
		tx:      transaction,
		isTx:    true,
		Dialect: s.Dialect,
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
	b.statement = clause.Select{
		Table:   b.tempTable,
		Columns: columns,
	}
	b.WhereStatements = clause.WhereStatements{}
	b.Values = []any{}

	return b
}

func (s *SQLBuilder) Insert(data any) (int64, error) {
	dataMap := map[string]any{}
	dataType := reflect.TypeOf(data)

	if dataType.Kind() == reflect.Struct {
		err := toMap(data, &dataMap)

		if err != nil {
			return 0, err
		}
	} else {
		dataMap = data.(map[string]any)
	}

	s.insertStatement = clause.Insert{
		Table: s.tempTable,
		Rows: []map[string]any{
			dataMap,
		},
	}

	res, err := s.Exec()
	if err != nil {
		return 0, err
	}

	lastInsertId, err := res.LastInsertId()

	if err != nil {
		return 0, nil
	}

	return lastInsertId, nil
}

func (s *SQLBuilder) InsertMany(data []map[string]any) *SQLBuilder {
	s.insertStatement = clause.Insert{
		Table: s.tempTable,
		Rows:  data,
	}

	return s
}

func (s *SQLBuilder) Update(data map[string]any) *SQLBuilder {
	s.updateStatement = clause.Update{
		Table: s.tempTable,
		Rows:  data,
	}
	s.WhereStatements = clause.WhereStatements{}
	s.Values = []any{}

	return s
}

func (s *SQLBuilder) Delete() *SQLBuilder {
	s.deleteStatement = clause.Delete{
		Table: s.tempTable,
	}
	s.WhereStatements = clause.WhereStatements{}
	s.Values = []any{}

	return s
}

func (s *SQLBuilder) Table(table string) *SQLBuilder {
	s.tempTable = table

	return s
}

func (s *SQLBuilder) GetSql() (string, error) {

	if s.statement.Table != "" {
		stmt, statementObj := s.statement.Parse(s.Dialect)
		stmt += s.WhereStatements.Parse(s.Dialect)
		stmt += s.Dialect.ParseGroup(s.Grouping)
		stmt += s.Dialect.ParseLimit(s.Limiting)
		stmt += s.Dialect.ParseOffset(s.Offseting)
		stmt += s.Dialect.ParseOrder(s.Ordering)

		s.Values = append(s.Values, statementObj.Values...)
		s.Values = append(s.Values, s.WhereStatements.Values...)
		return stmt, nil
	}

	if s.insertStatement.Table != "" {
		return s.insertStatement.Parse(s.Dialect), nil
	}

	if s.updateStatement.Table != "" {
		stmt, statementObj := s.updateStatement.Parse(s.Dialect)
		stmt += s.WhereStatements.Parse(s.Dialect)

		s.Values = append(s.Values, statementObj.Values...)
		s.Values = append(s.Values, s.WhereStatements.Values...)
		return stmt, nil
	}

	if s.deleteStatement.Table != "" {
		stmt, statementObj := s.deleteStatement.Parse(s.Dialect)
		stmt += s.WhereStatements.Parse(s.Dialect)

		s.Values = append(s.Values, statementObj.Values...)
		s.Values = append(s.Values, s.WhereStatements.Values...)
		return stmt, nil
	}

	return "", errors.New("no valid statement found")
}

func (s *SQLBuilder) GetArguments() []any {

	values := []any{}
	if !reflect.DeepEqual(s.statement, clause.Select{}) {
		values = append(values, s.Values...)
	}

	if !reflect.DeepEqual(s.insertStatement, clause.Insert{}) {
		for _, row := range s.insertStatement.Rows {
			keys := make([]string, 0, len(row))
			for k := range row {
				keys = append(keys, k)
			}
			slices.Sort(keys)

			for k := range keys {
				if val, ok := row[keys[k]]; ok {
					values = append(values, val)
				}
			}
		}
	}

	if !reflect.DeepEqual(s.updateStatement, clause.Update{}) {
		values = append(values, s.Values...)
	}

	if !reflect.DeepEqual(s.deleteStatement, clause.Delete{}) {
		values = append(values, s.Values...)
	}

	return values
}

func (s *SQLBuilder) Where(field string, Op clause.Operator, val any) *SQLBuilder {
	where := clause.Where{
		Field: field,
		Value: val,
		Op:    Op,
		Conj:  clause.ConjuctionAnd,
	}
	s.WhereStatements.Where = append(s.WhereStatements.Where, where)
	return s
}

func (s *SQLBuilder) WhereOr(field string, Op clause.Operator, val any) *SQLBuilder {
	where := clause.Where{
		Field: field,
		Value: val,
		Op:    Op,
		Conj:  clause.ConjuctionOr,
	}
	s.WhereStatements.Where = append(s.WhereStatements.Where, where)
	return s
}

func (s *SQLBuilder) WhereIn(field string, values []any) *SQLBuilder {
	wherein := clause.WhereIn{
		Field:  field,
		Values: values,
	}
	s.WhereStatements.WhereIn = append(s.WhereStatements.WhereIn, wherein)

	return s
}

func (s *SQLBuilder) WhereNotIn(field string, values []any) *SQLBuilder {
	wherenotin := clause.WhereNotIn{
		Field:  field,
		Values: values,
	}
	s.WhereStatements.WhereNotIn = append(s.WhereStatements.WhereNotIn, wherenotin)
	return s
}

func (s *SQLBuilder) WhereBetween(field string, start any, end any) *SQLBuilder {
	wherebetween := clause.WhereBetween{
		Field: field,
		Start: start,
		End:   end,
	}
	s.WhereStatements.WhereBetween = append(s.WhereStatements.WhereBetween, wherebetween)
	return s
}

func (s *SQLBuilder) Join(table string, first string, operator clause.Operator, second string) *SQLBuilder {
	s.statement.Joins = append(s.statement.Joins, clause.Join{
		Type:        clause.InnerJoin,
		SecondTable: table,
		FirstTable:  s.statement.Table,
		On: clause.JoinON{
			LeftField:  first,
			Operator:   operator,
			RightField: second,
		},
	})
	return s
}

func (s *SQLBuilder) LeftJoin(table string, first string, operator clause.Operator, second string) *SQLBuilder {
	s.statement.Joins = append(s.statement.Joins, clause.Join{
		Type:        clause.LeftJoin,
		SecondTable: table,
		FirstTable:  s.statement.Table,
		On: clause.JoinON{
			LeftField:  first,
			Operator:   operator,
			RightField: second,
		},
	})
	return s
}

func (s *SQLBuilder) RightJoin(table string, first string, operator clause.Operator, second string) *SQLBuilder {
	s.statement.Joins = append(s.statement.Joins, clause.Join{
		Type:        clause.RightJoin,
		SecondTable: table,
		FirstTable:  s.statement.Table,
		On: clause.JoinON{
			LeftField:  first,
			Operator:   operator,
			RightField: second,
		},
	})
	return s
}

func (s *SQLBuilder) WhereFunc(field string, operator clause.Operator, builder func(b Builder) *SQLBuilder) *SQLBuilder {
	newBuilder := builder(New(s.Dialect, s.sql))
	where := clause.Where{
		Field: field,
		Op:    operator,
		SubStatement: clause.SubStatement{
			Select:          newBuilder.statement,
			WhereStatements: newBuilder.WhereStatements,
		},
	}
	s.WhereStatements.Where = append(s.WhereStatements.Where, where)
	return s
}

func (s *SQLBuilder) WhereExists(builder func(b Builder) *SQLBuilder) *SQLBuilder {
	newBuilder := builder(New(s.Dialect, s.sql))

	where := clause.Where{
		Op: clause.OperatorExists,
		SubStatement: clause.SubStatement{
			Select:          newBuilder.statement,
			WhereStatements: newBuilder.WhereStatements,
		},
	}
	s.WhereStatements.Where = append(s.WhereStatements.Where, where)

	return s
}

func (s *SQLBuilder) OrderBy(column string, dir clause.OrderDirection) *SQLBuilder {
	if len(s.Ordering.OrderingFields) > 0 {
		s.Ordering.OrderingFields = append(s.Ordering.OrderingFields, clause.OrderField{
			Field:     column,
			Direction: clause.OrderDirection(dir),
		})
		return s
	}

	s.Ordering = clause.Order{
		OrderingFields: []clause.OrderField{
			{
				Field:     column,
				Direction: clause.OrderDirection(dir),
			},
		},
	}
	return s
}
func (s *SQLBuilder) GroupBy(columns ...string) *SQLBuilder {
	if len(s.Grouping.Fields) > 0 {
		s.Grouping.Fields = append(s.Grouping.Fields, columns...)
	} else {
		s.Grouping = clause.GroupBy{
			Fields: columns,
		}
	}
	return s
}
func (s *SQLBuilder) Limit(n int64) *SQLBuilder {
	s.Limiting.Count = n
	s.statement.Values = append(s.statement.Values, n)
	return s
}
func (s *SQLBuilder) Offset(n int64) *SQLBuilder {
	s.Offseting.Count = n
	s.statement.Values = append(s.statement.Values, n)
	return s
}

func (b *SQLBuilder) Get(d any) error {
	ctx := context.Background()
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
func (b *SQLBuilder) GetContext(d any, ctx context.Context) error {
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

	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		return rows.Scan(d)
	}

	return nil
}
func (s *SQLBuilder) Exec() (sql.Result, error) {
	statement, err := s.GetSql()
	arguments := s.GetArguments()
	if err != nil {
		return nil, err
	}

	if s.isTx {
		return s.tx.Exec(statement, arguments...)
	}

	return s.sql.Exec(statement, arguments...)
}
func (s *SQLBuilder) ExecContext(ctx context.Context) (sql.Result, error) {
	statement, err := s.GetSql()
	arguments := s.GetArguments()
	if err != nil {
		return nil, err
	}

	if s.isTx {
		return s.tx.ExecContext(ctx, statement, arguments...)
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
		return s.tx.QueryContext(ctx, sql, arguments...)
	}

	return s.sql.QueryContext(ctx, sql, arguments...)
}
