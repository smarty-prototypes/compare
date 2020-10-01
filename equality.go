// Package equality facilitates comparisons of any two values.
package equality

import (
	"encoding/json"
	"fmt"
	"reflect"
	"runtime/debug"
	"strings"
	"testing"
	"time"
)

type Comparer interface{
	Compare(expected, actual interface{}) (result Comparison)
}

type comparer struct {
	options []Option
}

func New(options ...Option) Comparer {
	return comparer{options: options}
}

// Compare returns a comparison of expected and actual as well as
// a full report of any discrepancy between them.
func (this comparer) Compare(expected, actual interface{}) (result Comparison) {
	config := new(config)
	config.apply(this.options...)
	config.applyDefaultEqualitySpecs()

	result.ok = check(expected, actual, config.specs...)
	result.report = newFormatter(expected, actual, config).Format(result.ok)
	config.reportT(result)
	return result
}

func check(expected, actual interface{}, specs ...specFunc) bool {
	for _, factory := range specs {
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

type Comparison struct {
	ok     bool
	report string
}

func (this Comparison) OK() bool {
	return this.ok
}
func (this Comparison) Report() string {
	return this.report
}

type Option func(*config)

var Options options

type options struct{}

func (options) TestingT(t *testing.T) Option {
	return func(this *config) { this.t = t }
}
func (options) CompareNumerics() Option {
	return func(this *config) { this.appendSpec(newNumericEqualitySpecification) }
}
func (options) CompareTimes() Option {
	return func(this *config) { this.appendSpec(newTimeEqualitySpecification) }
}
func (options) CompareDeep() Option {
	return func(this *config) { this.appendSpec(newDeepEqualitySpecification) }
}
func (options) CompareEqual() Option {
	return func(this *config) { this.appendSpec(newEqualitySpecification) }
}
func (options) FormatVerb(verb string) Option {
	return func(this *config) {
		this.format = func(a interface{}) string {
			return fmt.Sprintf(verb, a)
		}
	}
}
func (options) FormatJSON() Option {
	return func(this *config) {
		this.format = func(a interface{}) string {
			serialized, err := json.Marshal(a)
			if err != nil {
				return err.Error()
			}
			return string(serialized)
		}
	}
}

type specFunc func(a, b interface{}) Specification

type config struct {
	t      *testing.T
	specs  []specFunc
	format func(interface{}) string
}

func (this *config) appendSpec(f specFunc) {
	this.specs = append(this.specs, f)
}

func (this *config) apply(options ...Option) {
	for _, option := range options {
		option(this)
	}
}
func (this *config) applyDefaultEqualitySpecs() {
	if len(this.specs) > 0 {
		return
	}
	this.apply(
		Options.CompareNumerics(),
		Options.CompareTimes(),
		Options.CompareDeep(),
	)
}
func (this *config) applyDefaultFormatting(expected interface{}) {
	if this.format != nil {
		return
	}

	switch {
	case isNumeric(expected):
		this.apply(Options.FormatVerb("%v"))
	case isTime(expected):
		this.apply(Options.FormatVerb("%v"))
	default:
		this.apply(Options.FormatVerb("%#v"))
	}
}

func (this *config) reportT(result Comparison) {
	if !result.OK() && this.t != nil {
		this.t.Error(result.Report())
	}
}

type Specification interface {
	IsSatisfied() bool
	AreEqual() bool
}

// deepEqualitySpecification compares any two values
// using reflect.DeepEqual.
//
// https://golang.org/pkg/reflect/#DeepEqual
//
type deepEqualitySpecification struct {
	a, b interface{}

	aType, bType reflect.Type
}

func newDeepEqualitySpecification(a, b interface{}) Specification {
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

// equalitySpecification compares any two values using
// the built-in equality operator (`==`).
//
// https://golang.org/ref/spec#Comparison_operators
//
type equalitySpecification struct {
	a, b interface{}
}

func newEqualitySpecification(a, b interface{}) Specification {
	return &equalitySpecification{a: a, b: b}
}

func (this *equalitySpecification) IsSatisfied() bool {
	return reflect.TypeOf(this.a) == reflect.TypeOf(this.b)
}

func (this *equalitySpecification) AreEqual() bool {
	return this.a == this.b
}

// numericEqualitySpecification compares numeric values using
// the built-in equality operator (`==`). Values of differing
// numeric reflect.Kind are each converted to the type of the
// other and are compared with `==` in both directions.
//
// https://golang.org/ref/spec#Comparison_operators
// https://golang.org/pkg/reflect/#Kind
//
type numericEqualitySpecification struct {
	a, b interface{}
}

func newNumericEqualitySpecification(a, b interface{}) Specification {
	return &numericEqualitySpecification{a: a, b: b}
}
func (this *numericEqualitySpecification) IsSatisfied() bool {
	return isNumeric(this.a) && isNumeric(this.b)
}

func (this *numericEqualitySpecification) AreEqual() bool {
	if this.a == this.b {
		return true
	}
	aValue := reflect.ValueOf(this.a)
	bValue := reflect.ValueOf(this.b)
	aAsB := aValue.Convert(bValue.Type()).Interface()
	bAsA := bValue.Convert(aValue.Type()).Interface()
	return this.a == bAsA && this.b == aAsB
}

func isNumeric(v interface{}) bool {
	kind := reflect.TypeOf(v).Kind()
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

// timeEqualitySpecification compares values both of type
// time.Time using their Equal method.
//
// https://golang.org/pkg/time/#Time.Equal
//
type timeEqualitySpecification struct {
	a time.Time
	b time.Time

	aOK bool
	bOK bool
}

func newTimeEqualitySpecification(a, b interface{}) Specification {
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

func isTime(v interface{}) bool {
	_, ok := v.(time.Time)
	return ok
}

type formatter struct {
	expected reflect.Value
	actual   reflect.Value
	config   *config
}

func newFormatter(expected, actual interface{}, config *config) *formatter {
	config.applyDefaultFormatting(expected)

	return &formatter{
		config:   config,
		expected: reflect.ValueOf(expected),
		actual:   reflect.ValueOf(actual),
	}
}

func (this formatter) Format(equal bool) string {
	if equal {
		return this.formatEqual()
	} else {
		return this.formatUnequal()
	}
}

func (this formatter) formatEqual() string {
	return fmt.Sprintf("%s == %s",
		this.config.format(this.actual),
		this.config.format(this.expected),
	)
}
func (this formatter) formatUnequal() string {
	expectedType := fmt.Sprintf("<%v>", this.expected.Type())
	actualType := fmt.Sprintf("<%v>", this.actual.Type())
	longestTypeName := max(len(expectedType), len(actualType))
	expectedType += strings.Repeat(" ", longestTypeName-len(expectedType))
	actualType += strings.Repeat(" ", longestTypeName-len(actualType))
	expectedV := this.config.format(this.expected)
	actualV := this.config.format(this.actual)
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
		debug.Stack(), // TODO: filter out lines from this file, go runtime, and go testing lib
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
