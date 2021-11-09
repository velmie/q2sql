package q2sql

import "strings"

const (
	mainDelim = ':'
	argsDelim = ','
)

// DefaultFilterExpressionParser is the default FilterExpressionParser
var DefaultFilterExpressionParser = NewDelimitedArgsParser(mainDelim, argsDelim)

// FilterExpressionParser parses given expression
//
// it returns filter name and its arguments
// error is returned in case if expression cannot be parsed
type FilterExpressionParser interface {
	ParseFilterExpression(expr string) (name string, args []string, err error)
}

// DelimitedArgsParser uses delimiters to split filter name and arguments
type DelimitedArgsParser struct {
	mainDelim byte
	argsDelim string
}

// NewDelimitedArgsParser is a DelimitedArgsParser constructor
func NewDelimitedArgsParser(
	mainDelim byte,
	argsDelim byte,
) *DelimitedArgsParser {
	return &DelimitedArgsParser{
		mainDelim,
		string([]byte{argsDelim}),
	}
}

// ParseFilterExpression retrieves filter name and arguments
//
// for example for this expression "any:1,2,3" colon delimiter might be used
// in order to specify filter name "any" and the comma delimiter might be used
// in order to specify a list of arguments 1 2 and 3
// so for the input expression "any:1,2,3" it returns ("any", []string{"1","2","3"}, nil)
func (d *DelimitedArgsParser) ParseFilterExpression(expr string) (name string, args []string, err error) {
	i := strings.IndexByte(expr, d.mainDelim)
	if i == -1 {
		return expr, nil, nil
	}
	name, expr = expr[:i], expr[i+1:]
	if expr == "" {
		return name, nil, nil
	}
	return name, strings.Split(expr, d.argsDelim), nil
}
