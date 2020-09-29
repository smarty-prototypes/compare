// Package equality facilitates comparisons of any two values,
// which are deemed equal herein if they meet any of the following
// specifications:
//
// 1. Numeric values are compared with the built-in equality
//    operator (`==`). Values of differing numeric reflect.Kind
//    are each converted to the type of the other and are
//    compared with `==` in both directions.
//    - https://golang.org/ref/spec#Comparison_operators
//    - https://golang.org/pkg/reflect/#Kind
// 2. Values both of type time.Time are compared using their Equal method.
//    - https://golang.org/pkg/time/#Time.Equal
// 3. All other values are compared using reflect.DeepEqual.
//    - https://golang.org/pkg/reflect/#DeepEqual
package equality

import (
	"testing"
)

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

// Report compares expected and actual and returns
// a full report of any discrepancy between them.
func Report(expected, actual interface{}, options ...Option) string {
	_, report := Compare(expected, actual, options...)
	return report
}

// Compare returns a comparison of expected and actual as well as
// a full report of any discrepancy between them.
func Compare(expected, actual interface{}, options ...Option) (ok bool, report string) {
	ok = Check(expected, actual, options...)
	if !ok {
		return ok, newFormatter(expected, actual, options...).String()
	}
	return ok, ""
}

// Check returns a comparison of expected and actual according
// to the specifications defined in this package.
func Check(expected, actual interface{}, options ...Option) bool {
	config := new(config)

	//config.apply(options...) // TODO: uncomment (tests)

	if len(config.specs) == 0 {
		config.apply(
			Options.CompareNumerics(),
			Options.CompareTimes(),
			Options.CompareDeep(),
		)
	}

	for _, factory := range config.specs {
		spec := factory(expected, actual)
		if !spec.IsSatisfied() {
			continue
		}
		if spec.AreEqual() {
			return true
		}
		break
	}
	return false
}
