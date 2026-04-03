package clause

import "github.com/suryaherdiyanto/sqlbuilder/dialect"

type Operator string
type JoinType string
type OrderDirection string
type Conjuction string

const (
	OperatorEqual             Operator = "="
	OperatorLessThan          Operator = "<"
	OperatorLessThanEqual     Operator = "<="
	OperatorGreaterThan       Operator = ">"
	OperatorGreatherThanEqual Operator = ">="
	OperatorNot               Operator = "!="
	OperatorNotQual           Operator = "<>"
	OperatorLike              Operator = "LIKE"
	OperatorILike             Operator = "ILIKE"
	OperatorNotLike           Operator = "NOT LIKE"
	OperatorExists            Operator = "EXISTS"
	OperatorNotExists         Operator = "NOT EXISTS"
)

const (
	LeftJoin  JoinType = "left join"
	RightJoin JoinType = "right join"
	InnerJoin JoinType = "inner join"
)

const (
	OrderDirectionASC  OrderDirection = "asc"
	OrderDirectionDESC OrderDirection = "desc"
)

const (
	ConjuctionAnd Conjuction = "AND"
	ConjuctionOr  Conjuction = "OR"
)

type SQLDialector interface {
	GetDelimiter() string
	GetColumnQuoteLeft() string
	GetColumnQuoteRight() string
	GetName() dialect.Dialect
}
