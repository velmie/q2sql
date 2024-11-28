package q2sql

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/velmie/qparser"
)

// Sqlizer return an SQL with the list of arguments
type Sqlizer interface {
	ToSQL() (string, []interface{}, error)
}

// Eq - Equality: field = value
type Eq struct {
	Field string
	Value interface{}
}

func (eq *Eq) ToSQL() (string, []interface{}, error) {
	return eq.Field + " = ?", []interface{}{eq.Value}, nil
}

// Neq - Non-equality: field != value
type Neq struct {
	Field string
	Value interface{}
}

func (neq *Neq) ToSQL() (string, []interface{}, error) {
	return neq.Field + " != ?", []interface{}{neq.Value}, nil
}

// Lt - Less than: field < value
type Lt struct {
	Field string
	Value interface{}
}

func (lt *Lt) ToSQL() (string, []interface{}, error) {
	return lt.Field + " < ?", []interface{}{lt.Value}, nil
}

// Le - Less than or equal to: field < value
type Le struct {
	Field string
	Value interface{}
}

func (le *Le) ToSQL() (string, []interface{}, error) {
	return le.Field + " <= ?", []interface{}{le.Value}, nil
}

// Gt - Greater than: field > value
type Gt struct {
	Field string
	Value interface{}
}

func (gt *Gt) ToSQL() (string, []interface{}, error) {
	return gt.Field + " > ?", []interface{}{gt.Value}, nil
}

// Ge - Greater than or equal to: field >= value
type Ge struct {
	Field string
	Value interface{}
}

func (ge *Ge) ToSQL() (string, []interface{}, error) {
	return ge.Field + " >= ?", []interface{}{ge.Value}, nil
}

// In - Equals one value from set: field IN (value, value2)
type In struct {
	Field  string
	Values []interface{}
}

func (in *In) ToSQL() (string, []interface{}, error) {
	if len(in.Values) == 0 {
		return "", nil, fmt.Errorf("'In' condition requires at least one value for field %s", in.Field)
	}
	return in.Field + " IN (?" + strings.Repeat(",?", len(in.Values)-1) + ")", in.Values, nil
}

// NotIn - Not in a set of values: field NOT IN (value, value2)
type NotIn struct {
	Field  string
	Values []interface{}
}

func (n *NotIn) ToSQL() (string, []interface{}, error) {
	if len(n.Values) == 0 {
		return "", nil, fmt.Errorf("'NotIn' condition requires at least one value for field %s", n.Field)
	}
	return n.Field + " NOT IN (?" + strings.Repeat(",?", len(n.Values)-1) + ")", n.Values, nil
}

// Like - contains text: field like value
type Like struct {
	Field string
	Value interface{}
}

func (l *Like) ToSQL() (string, []interface{}, error) {
	return l.Field + " LIKE ?", []interface{}{l.Value}, nil
}

// IsNull - Equal to null: field is null
type IsNull string

func (i IsNull) ToSQL() (string, []interface{}, error) {
	return string(i) + " IS NULL", nil, nil
}

// IsNotNull - Not equal to null: field is not null
type IsNotNull string

func (i IsNotNull) ToSQL() (string, []interface{}, error) {
	return string(i) + " IS NOT NULL", nil, nil
}

// Or connects multiple expressions with the "OR" statement
type Or []Sqlizer

func (or Or) ToSQL() (sql string, args []interface{}, err error) {
	if len(or) == 0 {
		return "", nil, fmt.Errorf("'Or' requires at least one condition")
	}
	group := len(or) > 1
	expr := new(bytes.Buffer)
	args = make([]interface{}, 0)
	if group {
		expr.WriteByte('(')
	}
	args, err = appendToSQL(or, expr, " OR ", args)
	if err != nil {
		return
	}
	if group {
		expr.WriteByte(')')
	}
	sql = expr.String()
	return
}

// And connects multiple expressions with the "AND" statement
type And []Sqlizer

func (and And) ToSQL() (sql string, args []interface{}, err error) {
	if len(and) == 0 {
		return "", nil, fmt.Errorf("'And' requires at least one condition")
	}

	group := len(and) > 1
	expr := new(bytes.Buffer)
	args = make([]interface{}, 0)

	if group {
		expr.WriteByte('(')
	}

	args, err = appendToSQL(and, expr, " AND ", args)
	if err != nil {
		return "", nil, err
	}

	if group {
		expr.WriteByte(')')
	}

	return expr.String(), args, nil
}

// RawSQLWithArgs is a free form sql with possible arguments
type RawSQLWithArgs struct {
	SQL  string
	Args []interface{}
}

func (r *RawSQLWithArgs) ToSQL() (string, []interface{}, error) {
	return r.SQL, r.Args, nil
}

// Columns is a helper that simplifies columns list creation
type Columns []string

func (s Columns) String() string {
	columns := ""
	if len(s) > 0 {
		columns = strings.Join(s, ", ")
	}
	return columns
}

func (s Columns) ToSQL() (string, []interface{}, error) {
	return s.String(), nil, nil
}

// OrderBy is a helper that simplifies creation
// of the "ORDER BY" SQL statement
type OrderBy []qparser.Sort

func (s OrderBy) String() string {
	if len(s) == 0 {
		return ""
	}

	var sb strings.Builder
	for i := 0; i < len(s); i++ {
		if i != 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(s[i].FieldName)
		sb.WriteString(" ")
		sb.WriteString(s[i].Order.String())
	}

	return sb.String()
}

func (s OrderBy) ToSQL() (string, []interface{}, error) {
	return s.String(), nil, nil
}

// Not - Negates a single expression: NOT (expression)
type Not struct {
	Expr Sqlizer
}

func (n *Not) ToSQL() (string, []interface{}, error) {
	sql, args, err := n.Expr.ToSQL()
	if err != nil {
		return "", nil, err
	}
	return "NOT (" + sql + ")", args, nil
}

// RawSQL is a raw SQL string without arguments
type RawSQL string

func (s RawSQL) ToSQL() (string, []interface{}, error) {
	return string(s), nil, nil
}

func appendToSQL(parts []Sqlizer, w io.Writer, sep string, args []interface{}) ([]interface{}, error) {
	for i, p := range parts {
		partSQL, partArgs, err := p.ToSQL()
		if err != nil {
			return nil, err
		}
		if partSQL == "" {
			continue
		}
		if i > 0 {
			if _, inErr := io.WriteString(w, sep); inErr != nil {
				return nil, inErr
			}
		}
		_, err = io.WriteString(w, partSQL)
		if err != nil {
			return nil, err
		}
		args = append(args, partArgs...)
	}
	return args, nil
}
