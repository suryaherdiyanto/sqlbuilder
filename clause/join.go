package clause

type JoinON struct {
	Operator   Operator
	LeftField  string
	RightField string
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
