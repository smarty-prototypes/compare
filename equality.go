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
	"fmt"
	"reflect"
	"runtime/debug"
	"strings"
	"testing"
	"time"
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
	config := new(Config)

	//config.apply(options...) // TODO: uncomment (tests)

	if len(config.specs) == 0 {
		Options.CompareNumerics()(config)
		Options.CompareTimes()(config)
		Options.CompareDeep()(config)
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

/**********************************************************************/

type equalitySpecification interface {
	IsSatisfied() bool
	AreEqual() bool
}

/**********************************************************************/

type deepEqualitySpecification struct {
	a, b interface{}

	aType, bType reflect.Type
}

func newDeepEqualitySpecification(a, b interface{}) equalitySpecification {
	return &deepEqualitySpecification{
		a: a,
		b: b,

		aType: reflect.TypeOf(a),
		bType: reflect.TypeOf(b),
	}
}
func (this *deepEqualitySpecification) IsSatisfied() bool {
	return this.aType == this.bType
}
func (this *deepEqualitySpecification) AreEqual() bool {
	return reflect.DeepEqual(this.a, this.b)
}

/**********************************************************************/

type timeEqualitySpecification struct {
	a time.Time
	b time.Time

	aOK bool
	bOK bool
}

func newTimeEqualitySpecification(a, b interface{}) equalitySpecification {
	this := &timeEqualitySpecification{}
	this.a, this.aOK = a.(time.Time)
	this.b, this.bOK = b.(time.Time)
	return this
}
func (this *timeEqualitySpecification) IsSatisfied() bool {
	return this.aOK && this.bOK
}
func (this *timeEqualitySpecification) AreEqual() bool {
	return this.a.Equal(this.b)
}

/**********************************************************************/

type numericEqualitySpecification struct {
	a, b interface{}

	aType, bType reflect.Type
}

func newNumericEqualitySpecification(a, b interface{}) equalitySpecification {
	return &numericEqualitySpecification{
		a:     a,
		b:     b,
		aType: reflect.TypeOf(a),
		bType: reflect.TypeOf(b),
	}
}
func (this *numericEqualitySpecification) IsSatisfied() bool {
	return isNumeric(this.aType.Kind()) && isNumeric(this.aType.Kind())
}

func (this *numericEqualitySpecification) AreEqual() bool {
	if this.a == this.b {
		return true
	}
	aValue := reflect.ValueOf(this.a)
	bValue := reflect.ValueOf(this.b)
	aAsB := aValue.Convert(this.bType).Interface()
	bAsA := bValue.Convert(this.aType).Interface()
	return this.a == bAsA && this.b == aAsB
}

func isNumeric(kind reflect.Kind) bool {
	return kind == reflect.Int ||
		kind == reflect.Int8 ||
		kind == reflect.Int16 ||
		kind == reflect.Int32 ||
		kind == reflect.Int64 ||
		kind == reflect.Uint ||
		kind == reflect.Uint8 ||
		kind == reflect.Uint16 ||
		kind == reflect.Uint32 ||
		kind == reflect.Uint64 ||
		kind == reflect.Float32 ||
		kind == reflect.Float64
}

/**********************************************************************/

type formatter struct {
	expected reflect.Value
	actual   reflect.Value
}

func newFormatter(expected, actual interface{}, options ...Option) *formatter {
	return &formatter{
		expected: reflect.ValueOf(expected),
		actual:   reflect.ValueOf(actual),
	}
}

func (this formatter) String() string {
	expectedType := fmt.Sprintf("<%v>", this.expected.Type())
	actualType := fmt.Sprintf("<%v>", this.actual.Type())
	longestTypeName := max(len(expectedType), len(actualType))
	expectedType += strings.Repeat(" ", longestTypeName-len(expectedType))
	actualType += strings.Repeat(" ", longestTypeName-len(actualType))
	// TODO: %+v or %#v or just %v... (maybe we try to rewrite pointers, interfaces, time.Times, or slices containing any of those?)
	// TODO: If the formatted values are of the same type, and appear equal, maybe we json serialize them?
	// TODO: perhaps we provide functional options that allow customization of the formatting?
	// - time.Time should use %v
	// - all numerics should be %v
	expectedV := fmt.Sprintf("%#v", this.expected)
	actualV := fmt.Sprintf("%#v", this.actual)
	valueDiff := this.diff(actualV, expectedV)
	typeDiff := this.diff(actualType, expectedType)

	return fmt.Sprintf("\n"+
		"Expected: %s %s\n"+
		"Actual:   %s %s\n"+
		"Diff:     %s %s\n"+
		"Stack:     \n%s\n",
		expectedType, expectedV,
		actualType, actualV,
		typeDiff, valueDiff,
		debug.Stack(),
	)
}

func (this formatter) diff(actualV string, expectedV string) string {
	result := new(strings.Builder)

	for x := 0; ; x++ {
		if x >= len(actualV) && x >= len(expectedV) {
			break
		}
		if x >= len(actualV) || x >= len(expectedV) || expectedV[x] != actualV[x] {
			result.WriteString("^")
		} else {
			result.WriteString(" ")
		}
	}
	return result.String()
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

/**********************************************************************/
