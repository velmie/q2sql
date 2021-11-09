package condition

import "github.com/velmie/q2sql"

// These constants are recommended names for the filters
const (
	NameEq         = "eq"
	NameIn         = "in"
	NameNeq        = "neq"
	NameIsNull     = "null"
	NameIsNotNull  = "notnull"
	NameLt         = "lt"
	NameLe         = "le"
	NameGt         = "gt"
	NameGe         = "ge"
	NameStartsWith = "startswith"
	NameEndsWith   = "endswith"
	NameContains   = "contains"
	NameLike       = "like"
)

var DefaultConditionMap = q2sql.ConditionMap{
	NameEq:         Eq,
	NameIn:         In,
	NameNeq:        Neq,
	NameIsNull:     IsNull,
	NameIsNotNull:  IsNotNull,
	NameLt:         Lt,
	NameLe:         Le,
	NameGt:         Gt,
	NameGe:         Ge,
	NameStartsWith: StartsWith,
	NameEndsWith:   EndsWith,
	NameContains:   Contains,
	NameLike:       Like,
}
