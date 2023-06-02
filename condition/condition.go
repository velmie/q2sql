package condition

import (
	"github.com/velmie/q2sql"
)

// Eq - equal to
func Eq(field string, args ...interface{}) (s q2sql.Sqlizer, err error) {
	if len(args) == 0 {
		return &q2sql.Eq{Field: field}, nil
	}
	return &q2sql.Eq{Field: field, Value: args[0]}, nil
}

// Neq - not equal to
func Neq(field string, args ...interface{}) (q2sql.Sqlizer, error) {
	if len(args) == 0 {
		return &q2sql.Neq{Field: field}, nil
	}
	return &q2sql.Neq{Field: field, Value: args[0]}, nil
}

// Lt - less than
func Lt(field string, args ...interface{}) (q2sql.Sqlizer, error) {
	if len(args) == 0 {
		return &q2sql.Lt{Field: field}, nil
	}
	return &q2sql.Lt{Field: field, Value: args[0]}, nil
}

// Le - less than or equal to
func Le(field string, args ...interface{}) (q2sql.Sqlizer, error) {
	if len(args) == 0 {
		return &q2sql.Le{Field: field}, nil
	}
	return &q2sql.Le{Field: field, Value: args[0]}, nil
}

// Gt - greater than
func Gt(field string, args ...interface{}) (q2sql.Sqlizer, error) {
	if len(args) == 0 {
		return &q2sql.Gt{Field: field}, nil
	}
	return &q2sql.Gt{Field: field, Value: args[0]}, nil
}

// Ge - greater than or equal to
func Ge(field string, args ...interface{}) (q2sql.Sqlizer, error) {
	if len(args) == 0 {
		return &q2sql.Ge{Field: field}, nil
	}
	return &q2sql.Ge{Field: field, Value: args[0]}, nil
}

// EndsWith - text ends with a substring
func EndsWith(field string, args ...interface{}) (q2sql.Sqlizer, error) {
	if len(args) == 0 {
		return &q2sql.Like{Field: field}, nil
	}
	val := args[0]
	switch v := val.(type) {
	case string:
		val = "%" + v
	case []byte:
		val = "%" + string(v)
	}
	return &q2sql.Like{Field: field, Value: val}, nil
}

// StartsWith - text starts with a substring
func StartsWith(field string, args ...interface{}) (q2sql.Sqlizer, error) {
	if len(args) == 0 {
		return &q2sql.Like{Field: field}, nil
	}
	val := args[0]
	switch v := val.(type) {
	case string:
		val = v + "%"
	case []byte:
		val = string(v) + "%"
	}
	return &q2sql.Like{Field: field, Value: val}, nil
}

// Contains - text contains a substring
func Contains(field string, args ...interface{}) (q2sql.Sqlizer, error) {
	if len(args) == 0 {
		return &q2sql.Like{Field: field}, nil
	}
	val := args[0]
	switch v := val.(type) {
	case string:
		val = "%" + v + "%"
	case []byte:
		val = "%" + string(v) + "%"
	}
	return &q2sql.Like{Field: field, Value: val}, nil
}

// Like - search for a specified pattern in a text
// where percent sign is a wildcard placeholder
func Like(field string, args ...interface{}) (q2sql.Sqlizer, error) {
	if len(args) == 0 {
		return &q2sql.Like{Field: field}, nil
	}
	return &q2sql.Like{Field: field, Value: args[0]}, nil
}

// In - match any of the given arguments
func In(field string, args ...interface{}) (q2sql.Sqlizer, error) {
	if len(args) == 0 {
		return &q2sql.In{Field: field, Values: []interface{}{"NULL"}}, nil
	}
	return &q2sql.In{Field: field, Values: args}, nil
}

// IsNull - must be null
func IsNull(field string, _ ...interface{}) (q2sql.Sqlizer, error) {
	return q2sql.IsNull(field), nil
}

// IsNotNull - must not be null
func IsNotNull(field string, _ ...interface{}) (q2sql.Sqlizer, error) {
	return q2sql.IsNotNull(field), nil
}
