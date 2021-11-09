package condition

import (
	"github.com/velmie/q2sql"
)

const nilSqlizer = q2sql.RawSql("")

// Eq - equal to
func Eq(field string, args ...interface{}) q2sql.Sqlizer {
	if len(args) == 0 {
		return nilSqlizer
	}
	return &q2sql.Eq{
		Field: field,
		Value: args[0],
	}
}

// Neq - not equal to
func Neq(field string, args ...interface{}) q2sql.Sqlizer {
	if len(args) == 0 {
		return nilSqlizer
	}
	return &q2sql.Neq{
		Field: field,
		Value: args[0],
	}
}

// Lt - less than
func Lt(field string, args ...interface{}) q2sql.Sqlizer {
	if len(args) == 0 {
		return nilSqlizer
	}
	return &q2sql.Lt{
		Field: field,
		Value: args[0],
	}
}

// Le - less than or equal to
func Le(field string, args ...interface{}) q2sql.Sqlizer {
	if len(args) == 0 {
		return nilSqlizer
	}
	return &q2sql.Le{
		Field: field,
		Value: args[0],
	}
}

// Gt - greater than
func Gt(field string, args ...interface{}) q2sql.Sqlizer {
	if len(args) == 0 {
		return nilSqlizer
	}
	return &q2sql.Gt{
		Field: field,
		Value: args[0],
	}
}

// Ge - greater than or equal to
func Ge(field string, args ...interface{}) q2sql.Sqlizer {
	if len(args) == 0 {
		return nilSqlizer
	}
	return &q2sql.Ge{
		Field: field,
		Value: args[0],
	}
}

// EndsWith - text ends with a substring
func EndsWith(field string, args ...interface{}) q2sql.Sqlizer {
	if len(args) == 0 {
		return nilSqlizer
	}
	val := args[0]
	switch v := val.(type) {
	case string:
		val = "%" + v
	case []byte:
		val = "%" + string(v)
	}
	return &q2sql.Like{
		Field: field,
		Value: val,
	}
}

// StartsWith - text starts with a substring
func StartsWith(field string, args ...interface{}) q2sql.Sqlizer {
	if len(args) == 0 {
		return nilSqlizer
	}
	val := args[0]
	switch v := val.(type) {
	case string:
		val = v + "%"
	case []byte:
		val = string(v) + "%"
	}
	return &q2sql.Like{
		Field: field,
		Value: val,
	}
}

// Contains - text contains a substring
func Contains(field string, args ...interface{}) q2sql.Sqlizer {
	if len(args) == 0 {
		return nilSqlizer
	}
	val := args[0]
	switch v := val.(type) {
	case string:
		val = "%" + v + "%"
	case []byte:
		val = "%" + string(v) + "%"
	}
	return &q2sql.Like{
		Field: field,
		Value: val,
	}
}

// Like - search for a specified pattern in a text
// where percent sign is a wildcard placeholder
func Like(field string, args ...interface{}) q2sql.Sqlizer {
	if len(args) == 0 {
		return nilSqlizer
	}
	return &q2sql.Like{
		Field: field,
		Value: args[0],
	}
}

// In - match any of the given arguments
func In(field string, args ...interface{}) q2sql.Sqlizer {
	if len(args) == 0 {
		return nilSqlizer
	}
	return &q2sql.In{
		Field:  field,
		Values: args,
	}
}

// IsNull - must be null
func IsNull(field string, _ ...interface{}) q2sql.Sqlizer {
	return q2sql.IsNull(field)
}

// IsNotNull - must not be null
func IsNotNull(field string, _ ...interface{}) q2sql.Sqlizer {
	return q2sql.IsNotNull(field)
}
