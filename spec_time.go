package equality

import "time"

// timeEqualitySpecification comparse values both of type
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
