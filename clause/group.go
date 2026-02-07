package clause

type GroupBy struct {
	Fields []string
}

func (g *GroupBy) Parse() string {
	stmt := ""
	for i, field := range g.Fields {
		stmt += field
		if i < len(g.Fields)-1 {
			stmt += ","
		}
	}

	return stmt
}
