package sqlbuilder

type WhereParser interface {
	Parse() string
}

type WhereInParser interface {
	Parse() string
}

type WhereBetweenParser interface {
	Parse() string
}

type WhereNotBetweenParser interface {
	Parse() string
}

type JoinParser interface {
	Parse() string
}

type StatementParser interface {
	Parse() string
}

type OrderParser interface {
	Parse() string
}

type ClauseParser interface {
	ParseWheres() string
	ParseJoins() string
	ParseWhereBetweens() string
	ParseWhereNotBetweens() string
	ParseWhereIn() string
	ParseWhereNotIn() string
	ParseOrdering() string
}
