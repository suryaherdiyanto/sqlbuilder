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
	ParseWhere(w Where) string
	ParseJoin(j Join) string
	ParseWhereBetween(wb WhereBetween) string
	ParseWhereNotBetween(wb WhereNotBetween) string
	ParseWhereIn(wi WhereIn) string
	ParseWhereNotIn(wi WhereNotIn) string
	ParseOrder(o Order) string
	ParseGroup(gb GroupBy) string
	ParseLimit(l Limit) string
	ParseOffset(o Offset) string
	ParseInsert(in Insert) string
}
