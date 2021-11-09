package q2sql

import (
	"reflect"
	"testing"
)

type delimitedArgsTest struct {
	in           string
	expectedName string
	expectedArgs []string
}

var delimitedArgsTestTests = []delimitedArgsTest{
	{
		in:           "notnull",
		expectedName: "notnull",
		expectedArgs: nil,
	},
	{
		in:           "Eq:1",
		expectedName: "Eq",
		expectedArgs: []string{"1"},
	},
	{
		in:           "oneOf:red,blue,green",
		expectedName: "oneOf",
		expectedArgs: []string{"red", "blue", "green"},
	},
}

func TestDelimitedArgsTest(t *testing.T) {
	parser := NewDelimitedArgsParser(':', ',')
	for _, test := range delimitedArgsTestTests {
		name, args, err := parser.ParseFilterExpression(test.in)
		if err != nil {
			t.Errorf("DelimitedArgsParser.ParseFilterExpression(%q) returned unexpected error %s", test.in, err)
			continue
		}
		if name != test.expectedName {
			t.Errorf(
				"DelimitedArgsParser.ParseFilterExpression(%q) returned name %q expected %q",
				test.in,
				name,
				test.expectedName,
			)
		}
		if !reflect.DeepEqual(args, test.expectedArgs) {
			t.Errorf(
				"DelimitedArgsParser.ParseFilterExpression(%q):\n\targs  %+v\n\twant %+v\n",
				test.in,
				args,
				test.expectedArgs,
			)
		}
	}
}
