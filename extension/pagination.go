package extension

import (
	"context"
	"fmt"
	"strconv"

	"github.com/velmie/qparser"

	"github.com/velmie/q2sql"
)

const Unlimited = int64(-1)

type LimitOffsetPaginationParams struct {
	MaxLimit            int64
	MaxOffset           int64
	LimitParameterName  string
	OffsetParameterName string
}

// LimitOffsetPagination is the extension
// that sets limit and offset based on the corresponding fields of the given query.Page
//
//nolint:gocognit // skip because it is covered by tests
func LimitOffsetPagination(maxLimit, maxOffset int64) q2sql.Extension {
	return func(_ context.Context, query *qparser.Query, builder *q2sql.SelectBuilder) error {
		if page := query.Page; page != nil {
			if page.Limit != "" {
				limit, err := strconv.ParseUint(page.Limit, 10, 32)
				if err != nil {
					return fmt.Errorf("page limit must be unsigned integer, got %q", page.Limit)
				}
				if maxLimit != Unlimited && int64(limit) > maxLimit {
					return fmt.Errorf("page limit cannot be greater than %d", maxLimit)
				}
				builder.Limit(limit)
			}
			if page.Offset != "" && page.Limit == "" {
				return fmt.Errorf("offset cannot be used without specifying limit")
			}
			if page.Offset != "" {
				offset, err := strconv.ParseUint(page.Offset, 10, 32)
				if err != nil {
					return fmt.Errorf("page offset must be unsigned integer, got %q", page.Offset)
				}
				if maxOffset != Unlimited && int64(offset) > maxOffset {
					return fmt.Errorf("page offset cannot be greater than %d", maxOffset)
				}
				builder.Offset(offset)
			}
		}
		return nil
	}
}

// LimitNumberPagination is the extension
// that sets limit and offset based on the corresponding fields of the given query.Page
// where offset is a result of the expression limit * (number - 1)
func LimitNumberPagination(maxLimit int64) q2sql.Extension {
	return Pagination(
		maxLimit,
		func(query *qparser.Query) (limit uint64, number uint64, err error) {
			page := query.Page
			if page == nil {
				return 0, 0, nil
			}
			if page.Limit == "" {
				return 0, 0, nil
			}
			limit, err = strconv.ParseUint(page.Limit, 10, 64)
			if err != nil {
				return 0, 0, fmt.Errorf("page limit must be unsigned integer, got %q", page.Limit)
			}
			number, err = strconv.ParseUint(page.Number, 10, 64)
			if err != nil {
				return 0, 0, fmt.Errorf("page number must be unsigned integer, got %q", page.Number)
			}
			return limit, number, nil
		})
}

// Pagination is the extension
// that sets limit and offset based on the size and number values returned by the "sizeAndNumberGetter"
// argument. Offset is a result of the expression size * (number - 1)
func Pagination(
	maxLimit int64,
	sizeAndNumberGetter func(*qparser.Query) (size uint64, number uint64, err error),
) q2sql.Extension {
	return func(_ context.Context, query *qparser.Query, builder *q2sql.SelectBuilder) error {
		size, number, err := sizeAndNumberGetter(query)
		if err != nil {
			return err
		}

		if size == 0 {
			return nil
		}

		if maxLimit != Unlimited && int64(size) > maxLimit {
			return fmt.Errorf("page limit cannot be greater than %d", maxLimit)
		}
		builder.Limit(size)

		if number > 1 {
			builder.Offset(size * (number - 1))
		}
		return nil
	}
}
