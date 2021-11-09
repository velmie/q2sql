package q2sql

import "fmt"

// Error defines string error
type Error string

// Error returns error message
func (e Error) Error() string {
	return string(e)
}

// FilterError
type FilterError struct {
	Filter  string
	Field   string
	Message string
}

func (a FilterError) Error() string {
	return a.Message
}

// TranslationError specifies what entry is failed to translate
type TranslationError struct {
	Entry   string
	Message string
}

func (t *TranslationError) Error() string {
	return fmt.Sprintf("failed to translate format of the %q entry because %s", t.Entry, t.Message)
}
