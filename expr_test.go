package q2sql

import (
	"fmt"
	"reflect"
	"testing"
)

type expressionTest struct {
	in        Sqlizer
	out       string
	args      []interface{}
	expectErr bool
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
		in: &In{
			Field:  fieldName,
			Values: []interface{}{},
		},
		expectErr: true,
	},
	{
		in: &NotIn{
			Field:  fieldName,
			Values: []interface{}{"4", "5", "6"},
		},
		out:  fmt.Sprintf("%s NOT IN (?,?,?)", fieldName),
		args: []interface{}{"4", "5", "6"},
	},
	{
		in: &NotIn{
			Field:  fieldName,
			Values: []interface{}{},
		},
		expectErr: true,
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
		in:   &Not{Expr: &Eq{Field: fieldName, Value: "not_eq"}},
		out:  fmt.Sprintf("NOT (%s = ?)", fieldName),
		args: []interface{}{"not_eq"},
	},
}

func TestExpressions(t *testing.T) {
	for i, tt := range expressionTests {
		meta := fmt.Sprintf("test %d (%s)", i, reflect.TypeOf(tt.in))
		expr, args, err := tt.in.ToSQL()
		if (err != nil) != tt.expectErr {
			t.Errorf("%s: unexpected error status: got %v, want %v", meta, err != nil, tt.expectErr)
			continue
		}
		if tt.expectErr {
			continue
		}
		if expr != tt.out {
			t.Errorf("%s: expected expression %q, got %q", meta, tt.out, expr)
			continue
		}
		if !reflect.DeepEqual(args, tt.args) {
			t.Errorf("%s: expected args %+v, got %+v", meta, tt.args, args)
		}
	}
}
