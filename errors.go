package arngin

import "strings"

// Errors is an alias for slice of errors which satisfies the error interface.
type Errors []error

func (errs Errors) Error() string {
	switch len(errs) {
	case 0:
		return "<nil>"
	case 1:
		return errs[0].Error()
	}

	ss := make([]string, len(errs))
	for i, err := range errs {
		ss[i] = err.Error()
	}

	return "multiple errors:\n\t" + strings.Join(ss, "\n\t")
}
