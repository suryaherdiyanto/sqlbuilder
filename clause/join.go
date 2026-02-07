package clause

type JoinON struct {
	Operator   Operator
	LeftValue  any
	RightValue any
}

type Join struct {
	Type        JoinType
	FirstTable  string
	SecondTable string
	On          JoinON
}

func (j Join) Parse(dialect SQLDialector) string {
	return dialect.ParseJoin(j)
}
