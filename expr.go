package q2sql

import (
	"bytes"
	"io"
	"strings"

	"github.com/velmie/qparser"
)

// Sqlizer return an SQL with the list of arguments
type Sqlizer interface {
	ToSql() (string, []interface{}, error)
}

// Eq - Equality: field = value
type Eq struct {
	Field string
	Value interface{}
}

//nolint:stylecheck,gocritic // skip because this is an Interface implementation
func (eq *Eq) ToSql() (string, []interface{}, error) {
	return eq.Field + " = ?", []interface{}{eq.Value}, nil
}

// Neq - Non-equality: field != value
type Neq struct {
	Field string
	Value interface{}
}

//nolint:stylecheck,gocritic // skip because this is an Interface implementation
func (neq *Neq) ToSql() (string, []interface{}, error) {
	return neq.Field + " != ?", []interface{}{neq.Value}, nil
}

// Lt - Less than: field < value
type Lt struct {
	Field string
	Value interface{}
}

//nolint:stylecheck,gocritic // skip because this is an Interface implementation
func (lt *Lt) ToSql() (string, []interface{}, error) {
	return lt.Field + " < ?", []interface{}{lt.Value}, nil
}

// Le - Less than or equal to: field < value
type Le struct {
	Field string
	Value interface{}
}

//nolint:stylecheck,gocritic // skip because this is too small function; this is an Interface implementation
func (le *Le) ToSql() (string, []interface{}, error) {
	return le.Field + " <= ?", []interface{}{le.Value}, nil
}

// Gt - Greater than: field > value
type Gt struct {
	Field string
	Value interface{}
}

//nolint:stylecheck,gocritic // skip because this is too small function; this is an Interface implementation
func (gt *Gt) ToSql() (string, []interface{}, error) {
	return gt.Field + " > ?", []interface{}{gt.Value}, nil
}

// Ge - Greater than or equal to: field >= value
type Ge struct {
	Field string
	Value interface{}
}

//nolint:stylecheck,gocritic // skip because this is too small function; this is an Interface implementation
func (ge *Ge) ToSql() (string, []interface{}, error) {
	return ge.Field + " >= ?", []interface{}{ge.Value}, nil
}

// In - Equals one value from set: field IN (value, value2)
type In struct {
	Field  string
	Values []interface{}
}

//nolint:stylecheck,gocritic // skip because this is too small function; this is an Interface implementation
func (in *In) ToSql() (string, []interface{}, error) {
	return in.Field + " IN (?" + strings.Repeat(",?", len(in.Values)-1) + ")", in.Values, nil
}

// Like - contains text: field like value
type Like struct {
	Field string
	Value interface{}
}

//nolint:stylecheck,gocritic // skip because this is too small function; this is an Interface implementation
func (l *Like) ToSql() (string, []interface{}, error) {
	return l.Field + " LIKE ?", []interface{}{l.Value}, nil
}

// IsNull - Equal to null: field is null
type IsNull string

//nolint:stylecheck,gocritic // skip because this is too small function; this is an Interface implementation
func (i IsNull) ToSql() (string, []interface{}, error) {
	return string(i) + " IS NULL", nil, nil
}

// IsNotNull - Not equal to null: field is not null
type IsNotNull string

//nolint:stylecheck,gocritic // skip because this is too small function; this is an Interface implementation
func (i IsNotNull) ToSql() (string, []interface{}, error) {
	return string(i) + " IS NOT NULL", nil, nil
}

// Or connects multiple expressions with the "OR" statement
type Or []Sqlizer

//nolint:stylecheck // skip because this is an Interface implementation
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

// RawSqlWithArgs is a free form sql with possible arguments
type RawSqlWithArgs struct { //nolint:stylecheck // skip because fixing breaks client's code
	SQL  string
	Args []interface{}
}

//nolint:stylecheck,gocritic // skip because this is too small function; this is an Interface implementation
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

//nolint:stylecheck,gocritic // skip because this is too small function; this is an Interface implementation
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

//nolint:stylecheck,gocritic // skip because this is too small function; this is an Interface implementation
func (s OrderBy) ToSql() (string, []interface{}, error) {
	return s.String(), nil, nil
}

// RawSql is a raw SQL string without arguments
type RawSql string //nolint:stylecheck // skip because fixing breaks client's code

//nolint:stylecheck,gocritic // skip because this is too small function; this is an Interface implementation
func (s RawSql) ToSql() (string, []interface{}, error) {
	return string(s), nil, nil
}

func appendToSQL(parts []Sqlizer, w io.Writer, sep string, args []interface{}) ([]interface{}, error) {
	for i, p := range parts {
		partSQL, partArgs, err := p.ToSql()
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
