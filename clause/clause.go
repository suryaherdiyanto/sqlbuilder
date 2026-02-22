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
	GetDelimiter() string
	GetColumnQuoteLeft() string
	GetColumnQuoteRight() string
}
