package sqlbuilder

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/suryaherdiyanto/sqlbuilder/clause"
)

type SQLBuilder struct {
	Dialect              clause.SQLDialector
	sql                  *sql.DB
	tx                   *sql.Tx
	isTx                 bool
	enableLogging        bool
	logger               queryLogger
	tempTable            string
	rawStatement         string
	whereClauseStatement string
	selectStatement      string
	joinClauseStatement  string
	lockClauseStatement  string
	tailClauseStatement  string
	Values               []any
}

type queryLogger interface {
	Printf(format string, v ...any)
}

type Option func(*SQLBuilder)

type Builder interface {
	Select(columns ...string) *SQLBuilder
	Table(table string) *SQLBuilder
	Where(field string, Op clause.Operator, val any) *SQLBuilder
	OrWhere(field string, Op clause.Operator, val any) *SQLBuilder
	WhereIn(field string, values []any) *SQLBuilder
	OrWhereIn(field string, values []any) *SQLBuilder
	WhereNotIn(field string, values []any) *SQLBuilder
	OrWhereNotIn(field string, values []any) *SQLBuilder
	WhereBetween(field string, start any, end any) *SQLBuilder
	OrWhereBetween(field string, start any, end any) *SQLBuilder
	WhereFunc(field string, operator clause.Operator, b func(b Builder) *SQLBuilder) *SQLBuilder
	OrWhereFunc(field string, operator clause.Operator, b func(b Builder) *SQLBuilder) *SQLBuilder
	LockForShare() *SQLBuilder
	LockForUpdate() *SQLBuilder
	WhereDate(field string, operator clause.Operator, value any) *SQLBuilder
	OrWhereDate(field string, operator clause.Operator, value any) *SQLBuilder
	WhereExists(builder func(b Builder) *SQLBuilder) *SQLBuilder
	OrWhereExists(builder func(b Builder) *SQLBuilder) *SQLBuilder
	WhereMonth(field string, operator clause.Operator, value any) *SQLBuilder
	OrWhereMonth(field string, operator clause.Operator, value any) *SQLBuilder
	WhereYear(field string, operator clause.Operator, value any) *SQLBuilder
	OrWhereYear(field string, operator clause.Operator, value any) *SQLBuilder
	WhereDay(field string, operator clause.Operator, value any) *SQLBuilder
	OrWhereDay(field string, operator clause.Operator, value any) *SQLBuilder
	Join(table string, first string, operator clause.Operator, second string) *SQLBuilder
	LeftJoin(table string, first string, operator clause.Operator, second string) *SQLBuilder
	RightJoin(table string, first string, operator clause.Operator, second string) *SQLBuilder
	OrderBy(column string, dir clause.OrderDirection) *SQLBuilder
	GroupBy(columns ...string) *SQLBuilder
	Limit(n int64) *SQLBuilder
	Offset(n int64) *SQLBuilder
}

func New(dialect clause.SQLDialector, db *sql.DB, opts ...Option) *SQLBuilder {
	builder := &SQLBuilder{
		Dialect:       dialect,
		sql:           db,
		enableLogging: true,
		logger:        log.Default(),
	}

	for _, opt := range opts {
		opt(builder)
	}

	return builder
}

func WithLogging(enabled bool) Option {
	return func(s *SQLBuilder) {
		s.enableLogging = enabled
	}
}

func WithLogger(logger queryLogger) Option {
	return func(s *SQLBuilder) {
		if logger != nil {
			s.logger = logger
		}
	}
}

func (s *SQLBuilder) Begin(tx func(s *SQLBuilder) error) error {
	transaction, err := s.sql.Begin()
	if err != nil {
		return err
	}

	defer transaction.Rollback()

	builder := &SQLBuilder{
		tx:            transaction,
		isTx:          true,
		Dialect:       s.Dialect,
		enableLogging: s.enableLogging,
		logger:        s.logger,
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
	b.selectStatement = stmt

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
	s.clearStatement()
	s.Values = []any{}

	tableName := fmt.Sprintf("%s%s%s", s.Dialect.GetColumnQuoteLeft(), table, s.Dialect.GetColumnQuoteRight())
	s.tempTable = tableName

	selectStatement := clause.Select{
		Table:   s.tempTable,
		Columns: []string{"*"},
	}
	stmt, _ := selectStatement.Parse(s.Dialect)
	s.selectStatement = stmt

	return s
}

func (s *SQLBuilder) GetSql() string {
	if s.rawStatement != "" {
		return s.rawStatement
	}

	statement := s.selectStatement
	if s.joinClauseStatement != "" {
		statement = statement + " " + s.joinClauseStatement
	}

	if s.whereClauseStatement != "" {
		statement = statement + " " + s.whereClauseStatement
	}

	if s.lockClauseStatement != "" {
		statement = statement + " " + s.lockClauseStatement
	}
	if s.tailClauseStatement != "" {
		statement = statement + " " + s.tailClauseStatement
	}

	statement = strings.TrimSpace(statement)
	return statement
}

func (s *SQLBuilder) GetArguments() []any {
	return s.Values
}

func (s *SQLBuilder) Where(field string, Op clause.Operator, val any) *SQLBuilder {
	return s.addWhere(field, Op, val, clause.ConjuctionAnd)
}

func (s *SQLBuilder) OrWhere(field string, Op clause.Operator, val any) *SQLBuilder {
	return s.addWhere(field, Op, val, clause.ConjuctionOr)
}

func (s *SQLBuilder) WhereIn(field string, values []any) *SQLBuilder {
	return s.addWhereIn(field, values, clause.ConjuctionAnd)
}

func (s *SQLBuilder) OrWhereIn(field string, values []any) *SQLBuilder {
	return s.addWhereIn(field, values, clause.ConjuctionOr)
}

func (s *SQLBuilder) WhereNotIn(field string, values []any) *SQLBuilder {
	return s.addWhereNotIn(field, values, clause.ConjuctionAnd)
}

func (s *SQLBuilder) OrWhereNotIn(field string, values []any) *SQLBuilder {
	return s.addWhereNotIn(field, values, clause.ConjuctionOr)
}

func (s *SQLBuilder) WhereBetween(field string, start any, end any) *SQLBuilder {
	return s.addWhereBetween(field, start, end, clause.ConjuctionAnd)
}

func (s *SQLBuilder) OrWhereBetween(field string, start any, end any) *SQLBuilder {
	return s.addWhereBetween(field, start, end, clause.ConjuctionOr)
}

func (s *SQLBuilder) WhereDate(field string, operator clause.Operator, value any) *SQLBuilder {
	return s.addWhereDate(field, operator, value, clause.ConjuctionAnd)
}

func (s *SQLBuilder) OrWhereDate(field string, operator clause.Operator, value any) *SQLBuilder {
	return s.addWhereDate(field, operator, value, clause.ConjuctionOr)
}

func (s *SQLBuilder) WhereMonth(field string, operator clause.Operator, value any) *SQLBuilder {
	return s.addWhereMonth(field, operator, value, clause.ConjuctionAnd)
}

func (s *SQLBuilder) OrWhereMonth(field string, operator clause.Operator, value any) *SQLBuilder {
	return s.addWhereMonth(field, operator, value, clause.ConjuctionOr)
}

func (s *SQLBuilder) WhereYear(field string, operator clause.Operator, value any) *SQLBuilder {
	return s.addWhereYear(field, operator, value, clause.ConjuctionAnd)
}

func (s *SQLBuilder) OrWhereYear(field string, operator clause.Operator, value any) *SQLBuilder {
	return s.addWhereYear(field, operator, value, clause.ConjuctionOr)
}

func (s *SQLBuilder) WhereDay(field string, operator clause.Operator, value any) *SQLBuilder {
	return s.addWhereDay(field, operator, value, clause.ConjuctionAnd)
}

func (s *SQLBuilder) OrWhereDay(field string, operator clause.Operator, value any) *SQLBuilder {
	return s.addWhereDay(field, operator, value, clause.ConjuctionOr)
}

func (s *SQLBuilder) LockForUpdate() *SQLBuilder {
	s.lockClauseStatement = clause.ForUpdate{IsLocking: true}.Parse()
	return s
}

func (s *SQLBuilder) LockForShare() *SQLBuilder {
	s.lockClauseStatement = clause.ForShare{IsLocking: true}.Parse()
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
	s.joinClauseStatement = s.concatJoinClause(s.joinClauseStatement, join)
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
	s.joinClauseStatement = s.concatJoinClause(s.joinClauseStatement, join)
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
	s.joinClauseStatement = s.concatJoinClause(s.joinClauseStatement, join)
	return s
}

func (s *SQLBuilder) WhereFunc(field string, operator clause.Operator, builder func(b Builder) *SQLBuilder) *SQLBuilder {
	return s.addWhereFunc(field, operator, clause.ConjuctionAnd, builder)
}

func (s *SQLBuilder) OrWhereFunc(field string, operator clause.Operator, builder func(b Builder) *SQLBuilder) *SQLBuilder {
	return s.addWhereFunc(field, operator, clause.ConjuctionOr, builder)
}

func (s *SQLBuilder) WhereExists(builder func(b Builder) *SQLBuilder) *SQLBuilder {
	return s.addWhereExists(clause.ConjuctionAnd, builder)
}

func (s *SQLBuilder) OrWhereExists(builder func(b Builder) *SQLBuilder) *SQLBuilder {
	return s.addWhereExists(clause.ConjuctionOr, builder)
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
	s.tailClauseStatement = s.tailClauseStatement + " " + order.Parse(s.Dialect)
	return s
}

func (s *SQLBuilder) GroupBy(columns ...string) *SQLBuilder {
	grouping := clause.GroupBy{
		Fields: columns,
	}
	s.tailClauseStatement = s.concatTailClause(s.tailClauseStatement, grouping)
	return s
}
func (s *SQLBuilder) Limit(n int64) *SQLBuilder {
	limit := clause.Limit{
		Count: n,
	}
	s.tailClauseStatement = s.concatTailClause(s.tailClauseStatement, limit)
	s.Values = append(s.Values, n)
	return s
}
func (s *SQLBuilder) Offset(n int64) *SQLBuilder {
	offset := clause.Offset{
		Count: n,
	}
	s.tailClauseStatement = s.concatTailClause(s.tailClauseStatement, offset)
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

func (b *SQLBuilder) Count() (int64, error) {
	var count int64
	selectStatement := clause.Select{
		Table:   b.tempTable,
		Columns: []string{"COUNT(*) AS count"},
	}
	stmt, _ := selectStatement.Parse(b.Dialect)
	b.selectStatement = stmt

	rows, err := b.runQuery(context.Background())
	if err != nil {
		return 0, err
	}

	if rows.Next() {
		if err = rows.Scan(&count); err != nil {
			return 0, err
		}
	}

	return count, nil
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

	startedAt := time.Now()
	defer s.logQuery(statement, arguments, startedAt)

	if s.isTx {
		return s.tx.Exec(statement, arguments...)
	}

	return s.sql.Exec(statement, arguments...)
}
func (s *SQLBuilder) ExecContext(ctx context.Context) (sql.Result, error) {
	statement := s.GetSql()
	arguments := s.GetArguments()

	startedAt := time.Now()
	defer s.logQuery(statement, arguments, startedAt)

	if s.isTx {
		return s.tx.ExecContext(ctx, statement, arguments...)
	}

	return s.sql.ExecContext(ctx, statement, arguments...)
}

func (s *SQLBuilder) clearStatement() {
	s.rawStatement = ""
	s.tempTable = ""
	s.selectStatement = ""
	s.joinClauseStatement = ""
	s.whereClauseStatement = ""
	s.lockClauseStatement = ""
	s.tailClauseStatement = ""
}

func (s *SQLBuilder) concatWhereClause(statement string, conj clause.Conjuction, where clause.WhereParser) string {
	stmt := statement
	if strings.Contains(stmt, "WHERE") {
		stmt = stmt + " " + string(conj) + " " + where.Parse(s.Dialect)
		return stmt
	}
	stmt = stmt + "WHERE " + where.Parse(s.Dialect)

	return stmt
}

func (s *SQLBuilder) concatWhereWithSubquery(statement string, w clause.Where, subquery string) string {
	stmt := statement
	field := w.GetField(s.Dialect)
	cl := string(w.Op) + " (" + subquery + ")"

	if strings.Contains(stmt, "WHERE") {
		stmt = stmt + " " + string(w.Conj) + " "
		log.Printf("stmt: %s", stmt)
		if (w.Op != clause.OperatorExists) && (w.Op != clause.OperatorNotExists) {
			cl = field + " " + cl
		}
		stmt = stmt + cl

		return stmt
	}

	if (w.Op != clause.OperatorExists) && (w.Op != clause.OperatorNotExists) {
		cl = field + " " + cl
	}
	stmt = stmt + "WHERE " + cl

	return stmt
}

func (s *SQLBuilder) concatJoinClause(statement string, join clause.JoinParser) string {
	stmt := statement
	if strings.Contains(stmt, "JOIN") {
		stmt = stmt + " " + join.Parse(s.Dialect)
		return stmt
	}
	stmt = join.Parse(s.Dialect)
	return stmt
}

func (s *SQLBuilder) concatTailClause(statement string, tail clause.TailParser) string {
	stmt := statement
	if strings.Contains(stmt, "GROUP BY") || strings.Contains(stmt, "ORDER BY") || strings.Contains(stmt, "LIMIT") || strings.Contains(stmt, "OFFSET") {
		stmt = stmt + " " + tail.Parse(s.Dialect)
		return stmt
	}
	stmt = tail.Parse(s.Dialect)
	return stmt
}

func (s *SQLBuilder) addWhere(field string, op clause.Operator, val any, conj clause.Conjuction) *SQLBuilder {
	where := clause.Where{
		Field: field,
		Value: val,
		Op:    op,
		Conj:  conj,
	}
	s.Values = append(s.Values, val)

	s.whereClauseStatement = s.concatWhereClause(s.whereClauseStatement, where.Conj, where)

	return s
}

func (s *SQLBuilder) addWhereIn(field string, values []any, conj clause.Conjuction) *SQLBuilder {
	wherein := clause.WhereIn{
		Field:  field,
		Values: values,
		Conj:   conj,
	}
	s.Values = append(s.Values, values...)

	s.whereClauseStatement = s.concatWhereClause(s.whereClauseStatement, wherein.Conj, wherein)

	return s
}

func (s *SQLBuilder) addWhereNotIn(field string, values []any, conj clause.Conjuction) *SQLBuilder {
	wherenotin := clause.WhereNotIn{
		Field:  field,
		Values: values,
		Conj:   conj,
	}
	s.Values = append(s.Values, values...)
	s.whereClauseStatement = s.concatWhereClause(s.whereClauseStatement, wherenotin.Conj, wherenotin)
	return s
}

func (s *SQLBuilder) addWhereBetween(field string, start any, end any, conj clause.Conjuction) *SQLBuilder {
	wherebetween := clause.WhereBetween{
		Field: field,
		Start: start,
		End:   end,
		Conj:  conj,
	}
	s.Values = append(s.Values, start, end)

	s.whereClauseStatement = s.concatWhereClause(s.whereClauseStatement, wherebetween.Conj, wherebetween)
	return s
}

func (s *SQLBuilder) addWhereDate(field string, operator clause.Operator, value any, conj clause.Conjuction) *SQLBuilder {
	wheredate := clause.WhereDate{
		Field: field,
		Op:    operator,
		Value: value,
		Conj:  conj,
	}
	s.Values = append(s.Values, value)

	s.whereClauseStatement = s.concatWhereClause(s.whereClauseStatement, wheredate.Conj, wheredate)

	return s
}

func (s *SQLBuilder) addWhereMonth(field string, operator clause.Operator, value any, conj clause.Conjuction) *SQLBuilder {
	v := strconv.Itoa(value.(int))
	wheremonth := clause.WhereMonth{
		Field: field,
		Op:    operator,
		Value: v,
		Conj:  conj,
	}
	s.Values = append(s.Values, v)

	s.whereClauseStatement = s.concatWhereClause(s.whereClauseStatement, wheremonth.Conj, wheremonth)

	return s
}

func (s *SQLBuilder) addWhereYear(field string, operator clause.Operator, value any, conj clause.Conjuction) *SQLBuilder {
	v := strconv.Itoa(value.(int))
	whereyear := clause.WhereYear{
		Field: field,
		Op:    operator,
		Value: v,
		Conj:  conj,
	}
	s.Values = append(s.Values, v)
	s.whereClauseStatement = s.concatWhereClause(s.whereClauseStatement, whereyear.Conj, whereyear)
	return s
}

func (s *SQLBuilder) addWhereDay(field string, operator clause.Operator, value any, conj clause.Conjuction) *SQLBuilder {
	v := strconv.Itoa(value.(int))
	whereday := clause.WhereDay{
		Field: field,
		Op:    operator,
		Value: v,
		Conj:  conj,
	}
	s.Values = append(s.Values, v)
	s.whereClauseStatement = s.concatWhereClause(s.whereClauseStatement, whereday.Conj, whereday)
	return s
}

func (s *SQLBuilder) addWhereFunc(field string, operator clause.Operator, conj clause.Conjuction, builder func(b Builder) *SQLBuilder) *SQLBuilder {
	newBuilder := builder(New(s.Dialect, s.sql))
	where := clause.Where{
		Field: field,
		Op:    operator,
		Conj:  conj,
	}
	childStmt := newBuilder.GetSql()

	s.Values = append(s.Values, newBuilder.Values...)
	s.whereClauseStatement = s.concatWhereWithSubquery(s.whereClauseStatement, where, childStmt)
	return s
}

func (s *SQLBuilder) addWhereExists(conj clause.Conjuction, builder func(b Builder) *SQLBuilder) *SQLBuilder {
	newBuilder := builder(New(s.Dialect, s.sql))

	childStmt := newBuilder.GetSql()
	where := clause.Where{
		Op:   clause.OperatorExists,
		Conj: conj,
	}

	s.Values = append(s.Values, newBuilder.Values...)
	s.whereClauseStatement = s.concatWhereWithSubquery(s.whereClauseStatement, where, childStmt)

	return s
}

func (s *SQLBuilder) runQuery(ctx context.Context) (*sql.Rows, error) {
	sql := s.GetSql()
	arguments := s.GetArguments()

	startedAt := time.Now()
	defer s.logQuery(sql, arguments, startedAt)

	if s.isTx {
		return s.tx.QueryContext(ctx, sql, arguments...)
	}

	return s.sql.QueryContext(ctx, sql, arguments...)
}

func (s *SQLBuilder) logQuery(statement string, arguments []any, startedAt time.Time) {
	if !s.enableLogging || s.logger == nil {
		return
	}

	_ = arguments
	s.logger.Printf("sqlbuilder: %s - %s", statement, time.Since(startedAt))
}
