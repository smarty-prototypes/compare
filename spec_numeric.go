package equality

import "reflect"

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