package equality

import (
	"encoding/json"
	"fmt"
)

type Option func(*config)

var Options single

type single struct{}

func (single) CompareNumerics() Option {
	return func(this *config) {
		this.specs = append(this.specs, newNumericEqualitySpecification)
	}
}
func (single) CompareTimes() Option {
	return func(this *config) {
		this.specs = append(this.specs, newTimeEqualitySpecification)
	}
}
func (single) CompareDeep() Option {
	return func(this *config) {
		this.specs = append(this.specs, newDeepEqualitySpecification)
	}
}
func (single) FormatVerb(verb string) Option {
	return func(this *config) {
		this.format = func(a interface{}) string {
			return fmt.Sprintf(verb, a)
		}
	}
}
func (single) FormatJSON() Option {
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

// TODO: CompareEqual (==)
// TODO: CompareFloats (32 and/or 64, almost equal, within tolerance)
// TODO: IncludeDiffInFailureReport
// TODO: IncludeStackInFailureReport