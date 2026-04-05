package sqlbuilder

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/suryaherdiyanto/sqlbuilder/clause"
)

type SQLBuilder struct {
	Dialect         clause.SQLDialector
	sql             *sql.DB
	tx              *sql.Tx
	isTx            bool
	tempTable       string
	rawStatement    string
	statement       clause.Select
	insertStatement clause.Insert
	updateStatement clause.Update
	deleteStatement clause.Delete
	Grouping        clause.GroupBy
	Offseting       clause.Offset
	Limiting        clause.Limit
	Ordering        clause.Order
	Values          []any
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
	LockForShare() *SQLBuilder
	LockForUpdate() *SQLBuilder
	WhereDate(field string, operator clause.Operator, value any) *SQLBuilder
	WhereExists(builder func(b Builder) *SQLBuilder) *SQLBuilder
	WhereMonth(field string, operator clause.Operator, value any) *SQLBuilder
	WhereYear(field string, operator clause.Operator, value any) *SQLBuilder
	WhereDay(field string, operator clause.Operator, value any) *SQLBuilder
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

	defer transaction.Rollback()

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
	selectStatement := clause.Select{
		Table:   b.tempTable,
		Columns: columns,
	}

	stmt, _ := selectStatement.Parse(b.Dialect)
	b.rawStatement = stmt + b.rawStatement

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

	insertStatement := clause.Insert{
		Table: s.tempTable,
		Rows: []map[string]any{
			dataMap,
		},
	}

	stmt, insert := insertStatement.Parse(s.Dialect)
	s.rawStatement = stmt
	s.Values = append(s.Values, insert.Values...)

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
	insertStatement := clause.Insert{
		Table: s.tempTable,
		Rows:  data,
	}
	stmt, insert := insertStatement.Parse(s.Dialect)
	s.rawStatement = stmt
	s.Values = append(s.Values, insert.Values...)

	return s
}

func (s *SQLBuilder) Update(data map[string]any) (sql.Result, error) {
	updateStatement := clause.Update{
		Table: s.tempTable,
		Rows:  data,
	}

	stmt, update := updateStatement.Parse(s.Dialect)
	s.rawStatement = stmt + s.rawStatement

	updateValues := update.Values
	s.Values = append(updateValues, s.Values...)

	return s.Exec()
}

func (s *SQLBuilder) Delete() (sql.Result, error) {
	deleteStatement := clause.Delete{
		Table: s.tempTable,
	}

	stmt, _ := deleteStatement.Parse(s.Dialect)
	s.rawStatement = stmt + s.rawStatement

	return s.Exec()
}

func (s *SQLBuilder) Table(table string) *SQLBuilder {
	s.rawStatement = ""
	s.tempTable = ""
	s.Values = []any{}

	tableName := fmt.Sprintf("%s%s%s", s.Dialect.GetColumnQuoteLeft(), table, s.Dialect.GetColumnQuoteRight())
	s.tempTable = tableName

	return s
}

func (s *SQLBuilder) GetSql() string {
	return s.rawStatement
}

func (s *SQLBuilder) GetArguments() []any {
	return s.Values
}

func (s *SQLBuilder) Where(field string, Op clause.Operator, val any) *SQLBuilder {
	where := clause.Where{
		Field: field,
		Value: val,
		Op:    Op,
		Conj:  clause.ConjuctionAnd,
	}
	s.Values = append(s.Values, val)

	if strings.Contains(s.rawStatement, "WHERE") {
		s.rawStatement = s.rawStatement + " " + string(where.Conj) + " " + where.Parse(s.Dialect)
		return s
	}
	s.rawStatement = s.rawStatement + " WHERE " + where.Parse(s.Dialect)

	return s
}

func (s *SQLBuilder) WhereOr(field string, Op clause.Operator, val any) *SQLBuilder {
	where := clause.Where{
		Field: field,
		Value: val,
		Op:    Op,
		Conj:  clause.ConjuctionOr,
	}
	s.Values = append(s.Values, val)

	if strings.Contains(s.rawStatement, "WHERE") {
		s.rawStatement = s.rawStatement + " " + string(where.Conj) + " " + where.Parse(s.Dialect)
		return s
	}
	s.rawStatement = s.rawStatement + " WHERE " + where.Parse(s.Dialect)

	return s
}

func (s *SQLBuilder) WhereIn(field string, values []any) *SQLBuilder {
	wherein := clause.WhereIn{
		Field:  field,
		Values: values,
		Conj:   clause.ConjuctionAnd,
	}
	s.Values = append(s.Values, values...)

	if strings.Contains(s.rawStatement, "WHERE") {
		s.rawStatement = s.rawStatement + " " + string(wherein.Conj) + " " + wherein.Parse(s.Dialect)
		return s
	}
	s.rawStatement = s.rawStatement + " WHERE " + wherein.Parse(s.Dialect)

	return s
}

func (s *SQLBuilder) WhereNotIn(field string, values []any) *SQLBuilder {
	wherenotin := clause.WhereNotIn{
		Field:  field,
		Values: values,
	}
	s.Values = append(s.Values, values...)
	if strings.Contains(s.rawStatement, "WHERE") {
		s.rawStatement = s.rawStatement + " " + string(wherenotin.Conj) + " " + wherenotin.Parse(s.Dialect)
		return s
	}
	s.rawStatement = s.rawStatement + " WHERE " + wherenotin.Parse(s.Dialect)

	return s
}

func (s *SQLBuilder) WhereBetween(field string, start any, end any) *SQLBuilder {
	wherebetween := clause.WhereBetween{
		Field: field,
		Start: start,
		End:   end,
	}
	s.Values = append(s.Values, start, end)

	if strings.Contains(s.rawStatement, "WHERE") {
		s.rawStatement = s.rawStatement + " " + string(wherebetween.Conj) + " " + wherebetween.Parse(s.Dialect)
		return s
	}
	s.rawStatement = s.rawStatement + " WHERE " + wherebetween.Parse(s.Dialect)
	return s
}

func (s *SQLBuilder) WhereDate(field string, operator clause.Operator, value any) *SQLBuilder {
	wheredate := clause.WhereDate{
		Field: field,
		Op:    operator,
		Value: value,
		Conj:  clause.ConjuctionAnd,
	}
	s.Values = append(s.Values, value)

	if strings.Contains(s.rawStatement, "WHERE") {
		s.rawStatement = s.rawStatement + " " + string(wheredate.Conj) + " " + wheredate.Parse(s.Dialect)
		return s
	}
	s.rawStatement = s.rawStatement + " WHERE " + wheredate.Parse(s.Dialect)
	return s
}

func (s *SQLBuilder) WhereMonth(field string, operator clause.Operator, value any) *SQLBuilder {
	v := strconv.Itoa(value.(int))
	wheremonth := clause.WhereMonth{
		Field: field,
		Op:    operator,
		Value: v,
		Conj:  clause.ConjuctionAnd,
	}
	s.Values = append(s.Values, v)

	if strings.Contains(s.rawStatement, "WHERE") {
		s.rawStatement = s.rawStatement + " " + string(wheremonth.Conj) + " " + wheremonth.Parse(s.Dialect)
		return s
	}
	s.rawStatement = s.rawStatement + " WHERE " + wheremonth.Parse(s.Dialect)
	return s
}

func (s *SQLBuilder) WhereYear(field string, operator clause.Operator, value any) *SQLBuilder {
	v := strconv.Itoa(value.(int))
	whereyear := clause.WhereYear{
		Field: field,
		Op:    operator,
		Value: v,
		Conj:  clause.ConjuctionAnd,
	}
	s.Values = append(s.Values, v)
	if strings.Contains(s.rawStatement, "WHERE") {
		s.rawStatement = s.rawStatement + " " + string(whereyear.Conj) + " " + whereyear.Parse(s.Dialect)
		return s
	}

	s.rawStatement = s.rawStatement + " WHERE " + whereyear.Parse(s.Dialect)
	return s
}

func (s *SQLBuilder) WhereDay(field string, operator clause.Operator, value any) *SQLBuilder {
	v := strconv.Itoa(value.(int))
	whereday := clause.WhereDay{
		Field: field,
		Op:    operator,
		Value: v,
		Conj:  clause.ConjuctionAnd,
	}
	if strings.Contains(s.rawStatement, "WHERE") {
		s.rawStatement = s.rawStatement + " " + string(whereday.Conj) + " " + whereday.Parse(s.Dialect)
		return s
	}
	return s
}

func (s *SQLBuilder) LockForUpdate() *SQLBuilder {
	s.rawStatement = s.rawStatement + " " + clause.ForUpdate{IsLocking: true}.Parse()
	return s
}

func (s *SQLBuilder) LockForShare() *SQLBuilder {
	s.rawStatement = s.rawStatement + " " + clause.ForShare{IsLocking: true}.Parse()
	return s
}

func (s *SQLBuilder) Join(table string, first string, operator clause.Operator, second string) *SQLBuilder {
	join := clause.Join{
		Type:        clause.InnerJoin,
		SecondTable: table,
		On: clause.JoinON{
			LeftField:  first,
			Operator:   operator,
			RightField: second,
		},
	}
	s.rawStatement = s.rawStatement + " " + join.Parse(s.Dialect)
	return s
}

func (s *SQLBuilder) LeftJoin(table string, first string, operator clause.Operator, second string) *SQLBuilder {
	join := clause.Join{
		Type:        clause.LeftJoin,
		SecondTable: table,
		On: clause.JoinON{
			LeftField:  first,
			Operator:   operator,
			RightField: second,
		},
	}
	s.rawStatement = s.rawStatement + " " + join.Parse(s.Dialect)
	return s
}

func (s *SQLBuilder) RightJoin(table string, first string, operator clause.Operator, second string) *SQLBuilder {
	join := clause.Join{
		Type:        clause.RightJoin,
		SecondTable: table,
		On: clause.JoinON{
			LeftField:  first,
			Operator:   operator,
			RightField: second,
		},
	}
	s.rawStatement = s.rawStatement + " " + join.Parse(s.Dialect)
	return s
}

func (s *SQLBuilder) WhereFunc(field string, operator clause.Operator, builder func(b Builder) *SQLBuilder) *SQLBuilder {
	newBuilder := builder(New(s.Dialect, s.sql))
	where := clause.Where{
		Field: field,
		Op:    operator,
	}
	childStmt := newBuilder.GetSql()
	fieldQuoted := fmt.Sprintf("%s%s%s", s.Dialect.GetColumnQuoteLeft(), field, s.Dialect.GetColumnQuoteRight())

	if strings.Contains(s.rawStatement, "WHERE") {
		s.rawStatement = s.rawStatement + " " + string(where.Conj) + " " + fieldQuoted + " " + string(where.Op) + " " + "(" + childStmt + ")"
		return s
	}
	s.rawStatement = s.rawStatement + " WHERE " + fieldQuoted + " " + string(where.Op) + " " + "(" + childStmt + ")"
	s.Values = append(s.Values, newBuilder.Values...)
	return s
}

func (s *SQLBuilder) WhereExists(builder func(b Builder) *SQLBuilder) *SQLBuilder {
	newBuilder := builder(New(s.Dialect, s.sql))

	childStmt := newBuilder.GetSql()
	where := clause.Where{
		Op:   clause.OperatorExists,
		Conj: clause.ConjuctionAnd,
	}

	if strings.Contains(s.rawStatement, "WHERE") {
		s.rawStatement = s.rawStatement + " " + string(where.Conj) + " " + string(where.Op) + " " + "(" + childStmt + ")"
		return s
	}
	s.rawStatement = s.rawStatement + " WHERE " + string(where.Op) + " " + "(" + childStmt + ")"
	s.Values = append(s.Values, newBuilder.Values...)

	return s
}

func (s *SQLBuilder) OrderBy(column string, dir clause.OrderDirection) *SQLBuilder {
	order := clause.Order{
		OrderingFields: []clause.OrderField{
			{
				Field:     column,
				Direction: clause.OrderDirection(dir),
			},
		},
	}
	s.rawStatement = s.rawStatement + " " + order.Parse(s.Dialect)
	return s
}

func (s *SQLBuilder) GroupBy(columns ...string) *SQLBuilder {
	grouping := clause.GroupBy{
		Fields: columns,
	}
	s.rawStatement = s.rawStatement + " " + grouping.Parse(s.Dialect)
	return s
}
func (s *SQLBuilder) Limit(n int64) *SQLBuilder {
	limit := clause.Limit{
		Count: n,
	}
	s.rawStatement = s.rawStatement + " " + limit.Parse(s.Dialect)
	s.Values = append(s.Values, n)
	return s
}
func (s *SQLBuilder) Offset(n int64) *SQLBuilder {
	offset := clause.Offset{
		Count: n,
	}
	s.rawStatement = s.rawStatement + " " + offset.Parse(s.Dialect)
	return s
}

func (s *SQLBuilder) Raw(statement string, args ...any) *SQLBuilder {
	s.rawStatement = statement
	s.Values = append(s.Values, args...)
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
	statement := s.GetSql()
	arguments := s.GetArguments()

	if s.rawStatement != "" {
		s.rawStatement = ""
	}

	if s.isTx {
		return s.tx.Exec(statement, arguments...)
	}

	return s.sql.Exec(statement, arguments...)
}
func (s *SQLBuilder) ExecContext(ctx context.Context) (sql.Result, error) {
	statement := s.GetSql()
	arguments := s.GetArguments()

	if s.isTx {
		return s.tx.ExecContext(ctx, statement, arguments...)
	}

	return s.sql.ExecContext(ctx, statement, arguments...)
}
func (s *SQLBuilder) runQuery(ctx context.Context) (*sql.Rows, error) {
	sql := s.GetSql()
	arguments := s.GetArguments()

	if s.isTx {
		return s.tx.QueryContext(ctx, sql, arguments...)
	}

	return s.sql.QueryContext(ctx, sql, arguments...)
}
