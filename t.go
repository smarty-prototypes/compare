package equality

import "testing"

func T(t *testing.T) TT {
	return TT{T: t}
}

type TT struct{ *testing.T }

// Assert compares expected and actual and calls t.Error with
// a full report of any discrepancy between them.
func (this TT) Assert(expected, actual interface{}, options ...Option) bool {
	ok, report := Compare(expected, actual, options...)
	if !ok {
		this.Error(report)
	}
	return ok
}
