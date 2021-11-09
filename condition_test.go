package q2sql

import "testing"

type conditionMapTest struct {
	cmap ConditionMap
	name string
	err  error
}

var conditionMapTests = []conditionMapTest{
	{
		cmap: ConditionMap{
			"someName": nil,
		},
		name: "someName",
		err:  nil,
	},
	{
		cmap: ConditionMap{},
		name: "undefined",
		err:  ErrUndefinedCondition,
	},
}

func TestConditionMap(t *testing.T) {
	for _, test := range conditionMapTests {
		_, err := test.cmap.CreateCondition(test.name)
		if err != test.err {
			t.Errorf("ConditionMap.CreateCondition(%q) returned unexpected error %s", test.name, err)
			continue
		}

	}
}
