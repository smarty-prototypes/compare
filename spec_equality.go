package equality

import "reflect"

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
