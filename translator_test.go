package q2sql

import (
	"reflect"
	"testing"
)

type mapTranslatorTest struct {
	in   []string
	out  []string
	dict map[string]string
	err  error
}

var mapTranslatorTests = []mapTranslatorTest{
	{
		in:   []string{},
		out:  []string{},
		dict: make(map[string]string),
	},
	{
		in:   []string{"test"},
		out:  nil,
		dict: make(map[string]string),
		err: &TranslationError{
			Entry:   "test",
			Message: "translation is not found",
		},
	},
	{
		in:  []string{"id", "name", "profileId"},
		out: []string{"users.id", "users.name", "users.profile_id"},
		dict: map[string]string{
			"id":        "users.id",
			"name":      "users.name",
			"profileId": "users.profile_id",
		},
	},
}

func TestMapTranslator(t *testing.T) {
	for _, test := range mapTranslatorTests {
		got, err := MapTranslator(test.dict)(test.in)
		if err != nil && test.err == nil {
			t.Errorf("map translator returned unexpected error %s", err)
			continue
		}
		if !reflect.DeepEqual(got, test.out) {
			t.Errorf("translation of %v:\n\tgot  %+v\n\twant %+v\n", test.in, got, test.out)
		}
		if err != nil && test.err.Error() != err.Error() {
			t.Errorf("expected error %s\ngot %s", test.err, err)
			continue
		}
	}
}
