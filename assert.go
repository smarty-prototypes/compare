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
func (this TT) Assert(expected, actual interface{}) bool {
	ok, report := Compare(expected, actual)
	if !ok {
		this.Error("\n" + report)
	}
	return ok
}

// Report compares expected and actual and returns
// a full report of any discrepancy between them.
func Report(expected, actual interface{}) string {
	_, report := Compare(expected, actual)
	return report
}

// Compare returns a comparison of expected and actual as well as
// a full report of any discrepancy between them.
func Compare(expected, actual interface{}) (ok bool, report string) {
	ok = Check(expected, actual)
	if !ok {
		return ok, newFormatter(expected, actual).String()
	}
	return ok, ""
}

// Check returns a comparison of expected and actual according
// to the specifications defined in this package.
func Check(expected, actual interface{}) bool {
	specs := []equalitySpecification{
		newNumericEqualitySpecification(expected, actual),
		newTimeEqualitySpecification(expected, actual),
		newDeepEqualitySpecification(expected, actual),
	}

	for _, spec := range specs {
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

func newNumericEqualitySpecification(expected, actual interface{}) equalitySpecification {
	return &numericEqualitySpecification{
		a:     expected,
		b:     actual,
		aType: reflect.TypeOf(expected),
		bType: reflect.TypeOf(actual),
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

func newFormatter(expected, actual interface{}) *formatter {
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

	// If the formatted values are of the same type, and appear equal, maybe we json serialize them?

	// TODO: %+v or %#v or just %v...
	// TODO: diff marker ^
	return fmt.Sprintf(""+
		"Expected: %v %v\n"+
		"Actual:   %v %v\n"+
		"Diff:        %s\n"+
		"Stack:     \n%s\n",
		expectedType, this.expected,
		actualType, this.actual,
		"",
		debug.Stack(),
	)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

/**********************************************************************/