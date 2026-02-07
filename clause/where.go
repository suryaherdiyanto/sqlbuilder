package clause

type WhereGroup struct {
	Conj   Conjuction
	Wheres []Where
}

type Where struct {
	Field  string
	Op     Operator
	Value  any
	Conj   Conjuction
	Groups []WhereGroup
}

func (w Where) Parse(dialect SQLDialector) string {
	return dialect.ParseWhere(w)
}
