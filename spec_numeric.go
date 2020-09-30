package equality

import "reflect"

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
