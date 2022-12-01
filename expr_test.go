package q2sql

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/velmie/qparser"
)

type expressionTest struct {
	in   Sqlizer
	out  string
	args []interface{}
}

const fieldName = "field_x"

var expressionTests = []expressionTest{
	{
		in: &Eq{
			Field: fieldName,
			Value: "eq",
		},
		out:  fmt.Sprintf("%s = ?", fieldName),
		args: []interface{}{"eq"},
	},
	{
		in: &Neq{
			Field: fieldName,
			Value: "neq",
		},
		out:  fmt.Sprintf("%s != ?", fieldName),
		args: []interface{}{"neq"},
	},
	{
		in: &Lt{
			Field: fieldName,
			Value: "lt",
		},
		out:  fmt.Sprintf("%s < ?", fieldName),
		args: []interface{}{"lt"},
	},
	{
		in: &Le{
			Field: fieldName,
			Value: "le",
		},
		out:  fmt.Sprintf("%s <= ?", fieldName),
		args: []interface{}{"le"},
	},
	{
		in: &Gt{
			Field: fieldName,
			Value: "gt",
		},
		out:  fmt.Sprintf("%s > ?", fieldName),
		args: []interface{}{"gt"},
	},
	{
		in: &Ge{
			Field: fieldName,
			Value: "ge",
		},
		out:  fmt.Sprintf("%s >= ?", fieldName),
		args: []interface{}{"ge"},
	},
	{
		in: &Like{
			Field: fieldName,
			Value: "%like%",
		},
		out:  fmt.Sprintf("%s LIKE ?", fieldName),
		args: []interface{}{"%like%"},
	},
	{
		in: &In{
			Field:  fieldName,
			Values: []interface{}{"1", "2", "3"},
		},
		out:  fmt.Sprintf("%s IN (?,?,?)", fieldName),
		args: []interface{}{"1", "2", "3"},
	},
	{
		in:   IsNull(fieldName),
		out:  fmt.Sprintf("%s IS NULL", fieldName),
		args: nil,
	},
	{
		in:   IsNotNull(fieldName),
		out:  fmt.Sprintf("%s IS NOT NULL", fieldName),
		args: nil,
	},
	{
		in: OrderBy{
			{FieldName: "fieldA", Order: qparser.OrderDesc},
			{FieldName: "fieldB", Order: qparser.OrderAsc},
		},
		out:  "fieldA DESC, fieldB ASC",
		args: nil,
	},
	{
		in:   Columns{"col1", "col2", "col3"},
		out:  "col1, col2, col3",
		args: nil,
	},
	{
		in:   RawSql("SELECT 1"),
		out:  "SELECT 1",
		args: nil,
	},
	{
		in:   &RawSqlWithArgs{"id = ? OR id = ?", []interface{}{1, 2}},
		out:  "id = ? OR id = ?",
		args: []interface{}{1, 2},
	},
	{
		in:   Or{&RawSqlWithArgs{"field = ?", []interface{}{"value"}}, RawSql("field = 42")},
		out:  "(field = ? OR field = 42)",
		args: []interface{}{"value"},
	},
	{
		in:   Or{},
		out:  "",
		args: nil,
	},
	{
		in:   Or{RawSql("field = value")},
		out:  "field = value",
		args: []interface{}{},
	},
}

func TestExpressions(t *testing.T) {
	for i, tt := range expressionTests {
		meta := fmt.Sprintf("test %d (%s)", i, reflect.TypeOf(tt.in))
		expr, args, err := tt.in.ToSql()
		if err != nil {
			t.Errorf("%s: expr unexpected error: %s", meta, err)
			continue
		}
		if expr != tt.out {
			t.Errorf("%s: expected expression %q, got %q", meta, tt.out, expr)
			continue
		}
		if !reflect.DeepEqual(args, tt.args) {
			t.Errorf("%s: expected args  %+v, got %+v", meta, tt.args, args)
			continue
		}
	}
}
