package equality

import "testing"

func T(t *testing.T) TT {
	return TT{T: t}
}

// TT embeds *testing.T to provide an Assert method.
type TT struct{ *testing.T }

// Assert compares expected and actual and calls t.Error with
// a full report of any discrepancy between them.
func (this TT) Assert(expected, actual interface{}, options ...Option) bool {
	this.T.Helper()
	ok, report := Compare(expected, actual, options...)
	if !ok {
		this.T.Error(report)
	}
	return ok
}
