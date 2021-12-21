package extension

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/velmie/qparser"

	"github.com/velmie/q2sql"
)

type limitOffsetPaginationTest struct {
	p           q2sql.Extension
	q           *qparser.Query
	expectedErr bool
	limit       int64
	offset      int64
}

var limitOffsetPaginationTests = []limitOffsetPaginationTest{
	{
		p:           LimitOffsetPagination(Unlimited, Unlimited),
		q:           &qparser.Query{Page: &qparser.Page{Offset: "10"}},
		expectedErr: true,
	},
	{
		p:           LimitOffsetPagination(Unlimited, Unlimited),
		q:           &qparser.Query{Page: &qparser.Page{Limit: "invalid"}},
		expectedErr: true,
	},
	{
		p:           LimitOffsetPagination(Unlimited, Unlimited),
		q:           &qparser.Query{Page: &qparser.Page{Limit: "1", Offset: "invalid"}},
		expectedErr: true,
	},
	{
		p:           LimitOffsetPagination(Unlimited, Unlimited),
		q:           &qparser.Query{Page: &qparser.Page{Limit: "10"}},
		expectedErr: true,
		limit:       10,
	},
	{
		p:           LimitOffsetPagination(Unlimited, Unlimited),
		q:           &qparser.Query{Page: &qparser.Page{Limit: "500", Offset: "10"}},
		expectedErr: false,
		limit:       500,
		offset:      10,
	},
	{
		p:           LimitOffsetPagination(Unlimited, 5),
		q:           &qparser.Query{Page: &qparser.Page{Limit: "500", Offset: "10"}},
		expectedErr: true,
	},
	{
		p:           LimitOffsetPagination(5, Unlimited),
		q:           &qparser.Query{Page: &qparser.Page{Limit: "500", Offset: "10"}},
		expectedErr: true,
	},
}

func TestLimitOffsetPagination(t *testing.T) {
	ctx := context.Background()
	for i, tt := range limitOffsetPaginationTests {
		meta := fmt.Sprintf("test %d", i)
		b := new(q2sql.SelectBuilder)
		err := tt.p(ctx, tt.q, b)
		if !tt.expectedErr && err != nil {
			t.Errorf("%s, unexpected error %s", meta, err)
			continue
		} else if err != nil {
			continue
		}
		if tt.limit > 0 && strconv.FormatInt(tt.limit, 10) != b.LimitPart {
			t.Errorf("%s, unexpected limit part, want %d, got %s", meta, tt.limit, b.LimitPart)
		}
		if tt.offset > 0 && strconv.FormatInt(tt.offset, 10) != b.OffsetPart {
			t.Errorf("%s, unexpected offset part, want %d, got %s", meta, tt.offset, b.OffsetPart)
		}
	}
}

type limitNumberPaginationTest struct {
	p           q2sql.Extension
	q           *qparser.Query
	expectedErr bool
	limit       int64
	offset      int64
}

var limitNumberPaginationTests = []limitNumberPaginationTest{
	{
		p:           LimitNumberPagination(Unlimited),
		q:           &qparser.Query{Page: &qparser.Page{Number: "10"}},
		expectedErr: false,
	},
	{
		p:           LimitNumberPagination(Unlimited),
		q:           &qparser.Query{Page: &qparser.Page{Limit: "invalid"}},
		expectedErr: true,
	},
	{
		p:           LimitNumberPagination(Unlimited),
		q:           &qparser.Query{Page: &qparser.Page{Limit: "1", Number: "invalid"}},
		expectedErr: true,
	},
	{
		p:           LimitNumberPagination(Unlimited),
		q:           &qparser.Query{Page: &qparser.Page{Limit: "10", Number: "1"}},
		expectedErr: false,
		limit:       10,
	},
	{
		p:           LimitNumberPagination(Unlimited),
		q:           &qparser.Query{Page: &qparser.Page{Limit: "50", Number: "10"}},
		expectedErr: false,
		limit:       50,
		offset:      450,
	},
	{
		p:           LimitNumberPagination(5),
		q:           &qparser.Query{Page: &qparser.Page{Limit: "500"}},
		expectedErr: true,
	},
}

func TestLimitNumberPagination(t *testing.T) {
	ctx := context.Background()
	for i, tt := range limitNumberPaginationTests {
		meta := fmt.Sprintf("test %d", i)
		b := new(q2sql.SelectBuilder)
		err := tt.p(ctx, tt.q, b)
		if !tt.expectedErr && err != nil {
			t.Errorf("%s, unexpected error %s", meta, err)
			continue
		} else if err != nil {
			continue
		}
		if tt.limit > 0 && strconv.FormatInt(tt.limit, 10) != b.LimitPart {
			t.Errorf("%s, unexpected limit part, want %d, got %s", meta, tt.limit, b.LimitPart)
		}
		if tt.offset > 0 && strconv.FormatInt(tt.offset, 10) != b.OffsetPart {
			t.Errorf("%s, unexpected offset part, want %d, got %s", meta, tt.offset, b.OffsetPart)
		}
	}
}
