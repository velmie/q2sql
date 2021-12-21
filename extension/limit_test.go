package extension

import (
	"context"
	"strconv"

	"github.com/velmie/qparser"

	"github.com/velmie/q2sql"

	"testing"
)

func TestDefaultLimit(t *testing.T) {
	const (
		limit1 = 125
		limit2 = 999
	)

	b := new(q2sql.SelectBuilder)
	ext := DefaultLimit(limit1)
	err := ext(context.Background(), new(qparser.Query), b)
	if err != nil {
		t.Fatalf("unexpected error %s", err)
	}
	if strconv.FormatInt(limit1, 10) != b.LimitPart {
		t.Fatalf("want limit1 part %d, got %s", limit1, b.LimitPart)
	}

	ext = DefaultLimit(limit2)
	err = ext(context.Background(), new(qparser.Query), b)
	if err != nil {
		t.Fatalf("unexpected error %s", err)
	}
	if strconv.FormatInt(limit1, 10) != b.LimitPart {
		t.Fatal("limit part must not be changed when set")
	}
}
