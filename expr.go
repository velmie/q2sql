package q2sql

import (
	"bytes"
	"io"
	"strings"

	"github.com/velmie/qparser"
)

const alwaysTrue = RawSql("1 = 1")

// Sqlizer return an SQL with the list of arguments
type Sqlizer interface {
	ToSql() (string, []interface{}, error)
}

// Eq - Equality: field = value
type Eq struct {
	Field string
	Value interface{}
}

func (eq *Eq) ToSql() (string, []interface{}, error) {
	return eq.Field + " = ?", []interface{}{eq.Value}, nil
}

// Neq - Non-equality: field != value
type Neq struct {
	Field string
	Value interface{}
}

func (neq *Neq) ToSql() (string, []interface{}, error) {
	return neq.Field + " != ?", []interface{}{neq.Value}, nil
}

// Lt - Less than: field < value
type Lt struct {
	Field string
	Value interface{}
}

func (lt *Lt) ToSql() (string, []interface{}, error) {
	return lt.Field + " < ?", []interface{}{lt.Value}, nil
}

// Le - Less than or equal to: field < value
type Le struct {
	Field string
	Value interface{}
}

func (le *Le) ToSql() (string, []interface{}, error) {
	return le.Field + " <= ?", []interface{}{le.Value}, nil
}

// Gt - Greater than: field > value
type Gt struct {
	Field string
	Value interface{}
}

func (gt *Gt) ToSql() (string, []interface{}, error) {
	return gt.Field + " > ?", []interface{}{gt.Value}, nil
}

// Ge - Greater than or equal to: field >= value
type Ge struct {
	Field string
	Value interface{}
}

func (ge *Ge) ToSql() (string, []interface{}, error) {
	return ge.Field + " >= ?", []interface{}{ge.Value}, nil
}

// In - Equals one value from set: field IN (value, value2)
type In struct {
	Field  string
	Values []interface{}
}

func (in *In) ToSql() (string, []interface{}, error) {
	return in.Field + " IN (?" + strings.Repeat(",?", len(in.Values)-1) + ")", in.Values, nil
}

// Like - contains text: field like value
type Like struct {
	Field string
	Value interface{}
}

func (l *Like) ToSql() (string, []interface{}, error) {
	return l.Field + " LIKE ?", []interface{}{l.Value}, nil
}

// IsNull - Equal to null: field is null
type IsNull string

func (i IsNull) ToSql() (string, []interface{}, error) {
	return string(i) + " IS NULL", nil, nil
}

// IsNotNull - Not equal to null: field is not null
type IsNotNull string

func (i IsNotNull) ToSql() (string, []interface{}, error) {
	return string(i) + " IS NOT NULL", nil, nil
}

// Or connects multiple expressions with the "OR" statement
type Or []Sqlizer

func (or Or) ToSql() (sql string, args []interface{}, err error) {
	if len(or) == 0 {
		return
	}
	group := len(or) > 1
	expr := new(bytes.Buffer)
	args = make([]interface{}, 0)
	if group {
		expr.WriteByte('(')
	}
	args, err = appendToSql(or, expr, " OR ", args)
	if err != nil {
		return
	}
	if group {
		expr.WriteByte(')')
	}
	sql = expr.String()
	return
}

// RawSqlWithArgs is a free form sql with possible arguments
type RawSqlWithArgs struct {
	SQL  string
	Args []interface{}
}

func (r *RawSqlWithArgs) ToSql() (string, []interface{}, error) {
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

func (s Columns) ToSql() (string, []interface{}, error) {
	return s.String(), nil, nil
}

// OrderBy is a helper that simplifies creation
// of the "ORDER BY" SQL statement
type OrderBy []qparser.Sort

func (s OrderBy) String() string {
	if len(s) == 0 {
		return ""
	}
	var order string
	for i := 0; i < len(s); i++ {
		if i != 0 {
			order = order + ", "
		}
		order += s[i].FieldName + " " + s[i].Order.String()
	}
	return order
}

func (s OrderBy) ToSql() (string, []interface{}, error) {
	return s.String(), nil, nil
}

// RawSql is a raw SQL string without arguments
type RawSql string

func (s RawSql) ToSql() (string, []interface{}, error) {
	return string(s), nil, nil
}

func appendToSql(parts []Sqlizer, w io.Writer, sep string, args []interface{}) ([]interface{}, error) {
	for i, p := range parts {
		partSql, partArgs, err := p.ToSql()
		if err != nil {
			return nil, err
		}
		if partSql == "" {
			continue
		}
		if i > 0 {
			_, err := io.WriteString(w, sep)
			if err != nil {
				return nil, err
			}
		}
		_, err = io.WriteString(w, partSql)
		if err != nil {
			return nil, err
		}
		args = append(args, partArgs...)
	}
	return args, nil
}
