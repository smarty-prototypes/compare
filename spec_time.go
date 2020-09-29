package equality

import "time"

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
