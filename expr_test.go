package q2sql

import (
	"reflect"
	"testing"

	"github.com/velmie/qparser"
)

const fieldName = "field_x"

type expressionTest struct {
	Name      string
	in        Sqlizer
	out       string
	args      []interface{}
	expectErr bool
}

var expressionTests = []expressionTest{
	{
		Name: "Eq operator",
		in: &Eq{
			Field: fieldName,
			Value: "eq",
		},
		out:  "field_x = ?",
		args: []interface{}{"eq"},
	},
	{
		Name: "Neq operator",
		in: &Neq{
			Field: fieldName,
			Value: "neq",
		},
		out:  "field_x != ?",
		args: []interface{}{"neq"},
	},
	{
		Name: "Lt operator",
		in: &Lt{
			Field: fieldName,
			Value: "lt",
		},
		out:  "field_x < ?",
		args: []interface{}{"lt"},
	},
	{
		Name: "Le operator",
		in: &Le{
			Field: fieldName,
			Value: "le",
		},
		out:  "field_x <= ?",
		args: []interface{}{"le"},
	},
	{
		Name: "Gt operator",
		in: &Gt{
			Field: fieldName,
			Value: "gt",
		},
		out:  "field_x > ?",
		args: []interface{}{"gt"},
	},
	{
		Name: "Ge operator",
		in: &Ge{
			Field: fieldName,
			Value: "ge",
		},
		out:  "field_x >= ?",
		args: []interface{}{"ge"},
	},
	{
		Name: "Like operator",
		in: &Like{
			Field: fieldName,
			Value: "%like%",
		},
		out:  "field_x LIKE ?",
		args: []interface{}{"%like%"},
	},
	{
		Name: "In operator with values",
		in: &In{
			Field:  fieldName,
			Values: []interface{}{"1", "2", "3"},
		},
		out:  "field_x IN (?,?,?)",
		args: []interface{}{"1", "2", "3"},
	},
	{
		Name:      "In operator with empty values",
		in:        &In{Field: fieldName, Values: []interface{}{}},
		expectErr: true,
	},
	{
		Name: "NotIn operator with values",
		in: &NotIn{
			Field:  fieldName,
			Values: []interface{}{"4", "5", "6"},
		},
		out:  "field_x NOT IN (?,?,?)",
		args: []interface{}{"4", "5", "6"},
	},
	{
		Name:      "NotIn operator with empty values",
		in:        &NotIn{Field: fieldName, Values: []interface{}{}},
		expectErr: true,
	},
	{
		Name: "IsNull operator",
		in:   IsNull(fieldName),
		out:  "field_x IS NULL",
	},
	{
		Name: "IsNotNull operator",
		in:   IsNotNull(fieldName),
		out:  "field_x IS NOT NULL",
	},
	{
		Name: "Not operator",
		in:   &Not{Expr: &Eq{Field: fieldName, Value: "not_eq"}},
		out:  "NOT (field_x = ?)",
		args: []interface{}{"not_eq"},
	},
	{
		Name: "Or operator with single condition",
		in: Or{
			&Eq{Field: "field1", Value: "value1"},
		},
		out:  "field1 = ?",
		args: []interface{}{"value1"},
	},
	{
		Name: "Or operator with multiple conditions",
		in: Or{
			&Eq{Field: "field1", Value: "value1"},
			&Gt{Field: "field2", Value: 10},
		},
		out:  "(field1 = ? OR field2 > ?)",
		args: []interface{}{"value1", 10},
	},
	{
		Name: "Nested Or operators",
		in: Or{
			Or{&Eq{Field: "field1", Value: "value1"}},
			Or{&Gt{Field: "field2", Value: 10}},
		},
		out:  "(field1 = ? OR field2 > ?)",
		args: []interface{}{"value1", 10},
	},
	{
		Name: "And operator with single condition",
		in: And{
			&Neq{Field: "field3", Value: "value2"},
		},
		out:  "field3 != ?",
		args: []interface{}{"value2"},
	},
	{
		Name: "And operator with multiple conditions",
		in: And{
			&Neq{Field: "field3", Value: "value2"},
			&Lt{Field: "field4", Value: 20},
		},
		out:  "(field3 != ? AND field4 < ?)",
		args: []interface{}{"value2", 20},
	},
	{
		Name: "Nested And operators",
		in: And{
			&And{
				&Eq{Field: "field1", Value: "value1"},
				&Gt{Field: "field2", Value: 10},
			},
			&And{
				&Lt{Field: "field3", Value: 20},
				&Neq{Field: "field4", Value: "value2"},
			},
		},
		out:  "((field1 = ? AND field2 > ?) AND (field3 < ? AND field4 != ?))",
		args: []interface{}{"value1", 10, 20, "value2"},
	},
	{
		Name: "RawSQL without args",
		in:   RawSQL("field5 = 'fixed_value'"),
		out:  "field5 = 'fixed_value'",
	},
	{
		Name: "RawSQL with args",
		in: &RawSQLWithArgs{
			SQL:  "field6 = ?",
			Args: []interface{}{"dynamic_value"},
		},
		out:  "field6 = ?",
		args: []interface{}{"dynamic_value"},
	},
	{
		Name: "Columns",
		in:   Columns{"field7", "field8", "field9"},
		out:  "field7, field8, field9",
	},
	{
		Name: "OrderBy",
		in: OrderBy{
			qparser.Sort{FieldName: "field10", Order: qparser.OrderAsc},
			qparser.Sort{FieldName: "field11", Order: qparser.OrderDesc},
		},
		out: "field10 ASC, field11 DESC",
	},
}

func TestExpressions(t *testing.T) {
	for _, tt := range expressionTests {
		tt := tt // capture range variable
		t.Run(tt.Name, func(t *testing.T) {
			t.Parallel()
			expr, args, err := tt.in.ToSQL()
			if (err != nil) != tt.expectErr {
				t.Errorf("unexpected error status: got %v, want %v", err != nil, tt.expectErr)
				return
			}
			if tt.expectErr {
				return
			}
			if expr != tt.out {
				t.Errorf("expected expression %q, got %q", tt.out, expr)
			}
			if !reflect.DeepEqual(args, tt.args) {
				t.Errorf("expected args %+v, got %+v", tt.args, args)
			}
		})
	}
}
