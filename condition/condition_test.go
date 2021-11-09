package condition

import (
	"reflect"
	"testing"

	"github.com/velmie/q2sql"
)

type conditionTest struct {
	name      string
	condition q2sql.Condition
	field     string
	args      []interface{}
	expected  q2sql.Sqlizer
	test      func(t *testing.T, sqlizer q2sql.Sqlizer, field string, args []interface{})
}

var conditionTests = []conditionTest{
	{
		name:      "Eq",
		condition: Eq,
		field:     "field",
		args:      []interface{}{"eq_value"},
		expected:  &q2sql.Eq{},
		test: func(t *testing.T, sqlizer q2sql.Sqlizer, field string, args []interface{}) {
			eq := sqlizer.(*q2sql.Eq)
			if eq.Field != field {
				t.Errorf("expected Eq.Field to be %q, got %q", field, eq.Field)
			}
			if eq.Value != args[0] {
				t.Errorf("expected Eq.Value to be %+v, got %+v", args[0], eq.Value)
			}
		},
	},
	{
		name:      "Neq",
		condition: Neq,
		field:     "field",
		args:      []interface{}{"neq_value"},
		expected:  &q2sql.Neq{},
		test: func(t *testing.T, sqlizer q2sql.Sqlizer, field string, args []interface{}) {
			neq := sqlizer.(*q2sql.Neq)
			if neq.Field != field {
				t.Errorf("expected Neq.Field to be %q, got %q", field, neq.Field)
			}
			if neq.Value != args[0] {
				t.Errorf("expected Neq.Value to be %+v, got %+v", args[0], neq.Value)
			}
		},
	},
	{
		name:      "Lt",
		condition: Lt,
		field:     "field",
		args:      []interface{}{"lt_value"},
		expected:  &q2sql.Lt{},
		test: func(t *testing.T, sqlizer q2sql.Sqlizer, field string, args []interface{}) {
			lt := sqlizer.(*q2sql.Lt)
			if lt.Field != field {
				t.Errorf("expected Lt.Field to be %q, got %q", field, lt.Field)
			}
			if lt.Value != args[0] {
				t.Errorf("expected Lt.Value to be %+v, got %+v", args[0], lt.Value)
			}
		},
	},
	{
		name:      "Le",
		condition: Le,
		field:     "field",
		args:      []interface{}{"le_value"},
		expected:  &q2sql.Le{},
		test: func(t *testing.T, sqlizer q2sql.Sqlizer, field string, args []interface{}) {
			le := sqlizer.(*q2sql.Le)
			if le.Field != field {
				t.Errorf("expected Le.Field to be %q, got %q", field, le.Field)
			}
			if le.Value != args[0] {
				t.Errorf("expected Le.Value to be %+v, got %+v", args[0], le.Value)
			}
		},
	},
	{
		name:      "Gt",
		condition: Gt,
		field:     "field",
		args:      []interface{}{"gt_value"},
		expected:  &q2sql.Gt{},
		test: func(t *testing.T, sqlizer q2sql.Sqlizer, field string, args []interface{}) {
			gt := sqlizer.(*q2sql.Gt)
			if gt.Field != field {
				t.Errorf("expected Gt.Field to be %q, got %q", field, gt.Field)
			}
			if gt.Value != args[0] {
				t.Errorf("expected Gt.Value to be %+v, got %+v", args[0], gt.Value)
			}
		},
	},
	{
		name:      "Ge",
		condition: Ge,
		field:     "field",
		args:      []interface{}{"ge_value"},
		expected:  &q2sql.Ge{},
		test: func(t *testing.T, sqlizer q2sql.Sqlizer, field string, args []interface{}) {
			ge := sqlizer.(*q2sql.Ge)
			if ge.Field != field {
				t.Errorf("expected Ge.Field to be %q, got %q", field, ge.Field)
			}
			if ge.Value != args[0] {
				t.Errorf("expected Ge.Value to be %+v, got %+v", args[0], ge.Value)
			}
		},
	},
	{
		name:      "Like",
		condition: Like,
		field:     "field",
		args:      []interface{}{"like_value"},
		expected:  &q2sql.Like{},
		test: func(t *testing.T, sqlizer q2sql.Sqlizer, field string, args []interface{}) {
			like := sqlizer.(*q2sql.Like)
			if like.Field != field {
				t.Errorf("expected Like.Field to be %q, got %q", field, like.Field)
			}
			if like.Value != args[0] {
				t.Errorf("expected Like.Value to be %+v, got %+v", args[0], like.Value)
			}
		},
	},
	{
		name:      "EndsWith(bytes)",
		condition: EndsWith,
		field:     "field",
		args:      []interface{}{[]byte("endsWith_value")},
		expected:  &q2sql.Like{},
		test: func(t *testing.T, sqlizer q2sql.Sqlizer, field string, args []interface{}) {
			endsWith := sqlizer.(*q2sql.Like)
			if endsWith.Field != field {
				t.Errorf("expected EndsWith.Field to be %q, got %q", field, endsWith.Field)
			}
			arg := args[0].([]byte)
			expectedVal := "%" + string(arg)
			val := endsWith.Value.(string)
			if val != expectedVal {
				t.Errorf("expected EndsWith.Value to be %+v, got %+v", expectedVal, endsWith.Value)
			}
		},
	},
	{
		name:      "EndsWith(string)",
		condition: EndsWith,
		field:     "field",
		args:      []interface{}{"endsWith_value"},
		expected:  &q2sql.Like{},
		test: func(t *testing.T, sqlizer q2sql.Sqlizer, field string, args []interface{}) {
			endsWith := sqlizer.(*q2sql.Like)
			if endsWith.Field != field {
				t.Errorf("expected EndsWith.Field to be %q, got %q", field, endsWith.Field)
			}
			arg := args[0].(string)
			expectedVal := "%" + arg
			val := endsWith.Value.(string)
			if val != expectedVal {
				t.Errorf("expected EndsWith.Value to be %+v, got %+v", expectedVal, endsWith.Value)
			}
		},
	},
	{
		name:      "StartsWith(bytes)",
		condition: StartsWith,
		field:     "field",
		args:      []interface{}{[]byte("startsWith_value")},
		expected:  &q2sql.Like{},
		test: func(t *testing.T, sqlizer q2sql.Sqlizer, field string, args []interface{}) {
			startsWith := sqlizer.(*q2sql.Like)
			if startsWith.Field != field {
				t.Errorf("expected StartsWith.Field to be %q, got %q", field, startsWith.Field)
			}
			arg := args[0].([]byte)
			expectedVal := string(arg) + "%"
			val := startsWith.Value.(string)
			if val != expectedVal {
				t.Errorf("expected StartsWith.Value to be %+v, got %+v", expectedVal, startsWith.Value)
			}
		},
	},
	{
		name:      "StartsWith(string)",
		condition: StartsWith,
		field:     "field",
		args:      []interface{}{"startsWith_value"},
		expected:  &q2sql.Like{},
		test: func(t *testing.T, sqlizer q2sql.Sqlizer, field string, args []interface{}) {
			startsWith := sqlizer.(*q2sql.Like)
			if startsWith.Field != field {
				t.Errorf("expected StartsWith.Field to be %q, got %q", field, startsWith.Field)
			}
			arg := args[0].(string)
			expectedVal := arg + "%"
			val := startsWith.Value.(string)
			if val != expectedVal {
				t.Errorf("expected StartsWith.Value to be %+v, got %+v", expectedVal, startsWith.Value)
			}
		},
	},
	{
		name:      "Contains(bytes)",
		condition: Contains,
		field:     "field",
		args:      []interface{}{[]byte("contains_value")},
		expected:  &q2sql.Like{},
		test: func(t *testing.T, sqlizer q2sql.Sqlizer, field string, args []interface{}) {
			contains := sqlizer.(*q2sql.Like)
			if contains.Field != field {
				t.Errorf("expected Contains.Field to be %q, got %q", field, contains.Field)
			}
			arg := args[0].([]byte)
			expectedVal := "%" + string(arg) + "%"
			val := contains.Value.(string)
			if val != expectedVal {
				t.Errorf("expected Contains.Value to be %+v, got %+v", expectedVal, contains.Value)
			}
		},
	},
	{
		name:      "Contains(string)",
		condition: Contains,
		field:     "field",
		args:      []interface{}{"contains_value"},
		expected:  &q2sql.Like{},
		test: func(t *testing.T, sqlizer q2sql.Sqlizer, field string, args []interface{}) {
			contains := sqlizer.(*q2sql.Like)
			if contains.Field != field {
				t.Errorf("expected Contains.Field to be %q, got %q", field, contains.Field)
			}
			arg := args[0].(string)
			expectedVal := "%" + arg + "%"
			val := contains.Value.(string)
			if val != expectedVal {
				t.Errorf("expected Contains.Value to be %+v, got %+v", expectedVal, contains.Value)
			}
		},
	},
	{
		name:      "In",
		condition: In,
		field:     "field",
		args:      []interface{}{"in_value"},
		expected:  &q2sql.In{},
		test: func(t *testing.T, sqlizer q2sql.Sqlizer, field string, args []interface{}) {
			in := sqlizer.(*q2sql.In)
			if in.Field != field {
				t.Errorf("expected In.Field to be %q, got %q", field, in.Field)
			}
			if !reflect.DeepEqual(in.Values, args) {
				t.Errorf("expected In.Value to be %+v, got %+v", args, in.Values)
			}
		},
	},
	{
		name:      "IsNull",
		condition: IsNull,
		field:     "field",
		expected:  q2sql.IsNull(""),
		test: func(t *testing.T, sqlizer q2sql.Sqlizer, field string, args []interface{}) {
			isNull := sqlizer.(q2sql.IsNull)
			if string(isNull) != field {
				t.Errorf("expected IsNull to be wrapped field %q, got %q", field, isNull)
			}
		},
	},
	{
		name:      "IsNotNull",
		condition: IsNotNull,
		field:     "field",
		expected:  q2sql.IsNotNull(""),
		test: func(t *testing.T, sqlizer q2sql.Sqlizer, field string, args []interface{}) {
			isNotNull := sqlizer.(q2sql.IsNotNull)
			if string(isNotNull) != field {
				t.Errorf("expected IsNotNull to be wrapped field %q, got %q", field, isNotNull)
			}
		},
	},
}

func TestConditions(t *testing.T) {
	for _, tt := range conditionTests {
		sqlizer := tt.condition(tt.field, tt.args...)
		expectedType, gotType := reflect.TypeOf(tt.expected), reflect.TypeOf(sqlizer)
		if expectedType != gotType {
			t.Errorf("%s condition returned unexpected type: %s", tt.name, gotType)
			continue
		}
		if tt.test != nil {
			tt.test(t, sqlizer, tt.field, tt.args)
		}
	}
}

var nilSqlizerConditions = q2sql.ConditionMap{
	"Eq":         Eq,
	"In":         In,
	"Neq":        Neq,
	"Lt":         Lt,
	"Le":         Le,
	"Gt":         Gt,
	"Ge":         Ge,
	"StartsWith": StartsWith,
	"EndsWith":   EndsWith,
	"Contains":   Contains,
	"Like":       Like,
}

func TestReturnNilSqlizerIfNoArgs(t *testing.T) {
	const field = "field"
	nilSqlizerType := reflect.TypeOf(nilSqlizer)
	for name, c := range nilSqlizerConditions {
		sqlizerType := reflect.TypeOf(c(field))
		if sqlizerType != nilSqlizerType {
			t.Errorf("%s condition returned unexpected type: %s", name, sqlizerType)
			continue
		}
	}
}
