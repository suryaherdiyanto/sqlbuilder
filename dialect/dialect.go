package dialect

type SQLDialect struct {
	Name             Dialect
	Delimiter        string
	ColumnQuoteLeft  string
	ColumnQuoteRight string
}

type Dialect string

const (
	MySQL      Dialect = "mysql"
	PostgreSQL Dialect = "postgresql"
	SQLite     Dialect = "sqlite"
)

func (s SQLDialect) GetDelimiter() string {
	return s.Delimiter
}

func (s SQLDialect) GetColumnQuoteLeft() string {
	return s.ColumnQuoteLeft
}

func (s SQLDialect) GetColumnQuoteRight() string {
	return s.ColumnQuoteRight
}

func (s SQLDialect) GetName() Dialect {
	return s.Name
}

func New(delimiter, columnQuoteLeft, columnQuoteRight string) *SQLDialect {
	return &SQLDialect{
		Name:             SQLite,
		Delimiter:        delimiter,
		ColumnQuoteLeft:  columnQuoteLeft,
		ColumnQuoteRight: columnQuoteRight,
	}
}

func NewMySQL() *SQLDialect {
	return &SQLDialect{
		Name:             MySQL,
		Delimiter:        "?",
		ColumnQuoteLeft:  "`",
		ColumnQuoteRight: "`",
	}
}
