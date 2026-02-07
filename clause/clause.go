package clause

type Operator string
type JoinType string
type OrderDirection string
type Conjuction string

const (
	OperatorEqual             Operator = "="
	OperatorLessThan                   = "<"
	OperatorLessThanEqual              = "<="
	OperatorGreaterThan                = ">"
	OperatorGreatherThanEqual          = ">="
	OperatorNot                        = "!="
	OperatorLike                       = "LIKE"
	OperatorNotLike                    = "NOT LIKE"
	OperatorExists                     = "EXISTS"
	OperatorNotExists                  = "NOT EXISTS"
)

const (
	LeftJoin  JoinType = "left join"
	RightJoin          = "right join"
	InnerJoin          = "inner join"
)

const (
	OrderDirectionASC  OrderDirection = "asc"
	OrderDirectionDESC                = "desc"
)

const (
	ConjuctionAnd Conjuction = "AND"
	ConjuctionOr             = "OR"
)

type SQLDialector interface {
	ParseWhere() string
	ParseJoin() string
	ParseWhereBetween() string
	ParseWhereNotBetween() string
	ParseWhereIn() string
	ParseWhereNotIn() string
	ParseOrder() string
	ParseGroup() string

	NewWhere(field string, op Operator, value any, conj Conjuction)
	NewJoin(joinType JoinType, firstTable string, secondTable string, on JoinON)
	NewOrder(orderFields []OrderField)
	NewGroup(fields []string)
	NewWhereBetween(field string, start, end any, conj Conjuction)
	NewWhereNotBetween(field string, start, end any, conj Conjuction)
	NewWhereIn(field string, values []any, conj Conjuction)
	NewWhereNotIn(field string, values []any, conj Conjuction)
}
