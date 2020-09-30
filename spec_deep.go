package equality

import "reflect"

// deepEqualitySpecification compares any two values
// using reflect.DeepEqual.
//
// https://golang.org/pkg/reflect/#DeepEqual
//
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

