package dialect

import "fmt"

type PostgresDialect struct {
	placeholderIndex int
}

func (p *PostgresDialect) nextPlaceholder() string {
	p.placeholderIndex++
	return fmt.Sprintf("$%d", p.placeholderIndex)
}

func (p *PostgresDialect) GetDelimiter() string {
	return p.nextPlaceholder()
}

func (p *PostgresDialect) GetColumnQuoteLeft() string {
	return "\""
}

func (p *PostgresDialect) GetColumnQuoteRight() string {
	return "\""
}

// NewPostgres returns a PostgreSQL dialect instance.
func NewPostgres() *PostgresDialect {
	return &PostgresDialect{}
}
