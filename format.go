package equality

import (
	"fmt"
	"reflect"
	"runtime/debug"
	"strings"
)

type formatter struct {
	expected reflect.Value
	actual   reflect.Value
	format   func(interface{}) string
}

func newFormatter(expected, actual interface{}, options ...Option) *formatter {
	config := new(config)
	config.apply(options...)
	config.applyDefaultFormatting(expected)

	return &formatter{
		expected: reflect.ValueOf(expected),
		actual:   reflect.ValueOf(actual),
		format:   config.format,
	}
}

func (this formatter) String() string {
	expectedType := fmt.Sprintf("<%v>", this.expected.Type())
	actualType := fmt.Sprintf("<%v>", this.actual.Type())
	longestTypeName := max(len(expectedType), len(actualType))
	expectedType += strings.Repeat(" ", longestTypeName-len(expectedType))
	actualType += strings.Repeat(" ", longestTypeName-len(actualType))
	expectedV := this.format(this.expected)
	actualV := this.format(this.actual)
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
