package dialect

type SQLDialect struct {
	Delimiter        string
	ColumnQuoteLeft  string
	ColumnQuoteRight string
}

func (s SQLDialect) GetDelimiter() string {
	return s.Delimiter
}

func (s SQLDialect) GetColumnQuoteLeft() string {
	return s.ColumnQuoteLeft
}

func (s SQLDialect) GetColumnQuoteRight() string {
	return s.ColumnQuoteRight
}

func New(delimiter, columnQuoteLeft, columnQuoteRight string) *SQLDialect {
	return &SQLDialect{
		Delimiter:        delimiter,
		ColumnQuoteLeft:  columnQuoteLeft,
		ColumnQuoteRight: columnQuoteRight,
	}
}
