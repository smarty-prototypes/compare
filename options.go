package equality

type Config struct {
	specs []func(a, b interface{}) equalitySpecification
}

func (this *Config) apply(options ...Option) {
	for _, option := range options {
		option(this)
	}
}

type Option func(*Config)

var Options single

type single struct {}

func (single) CompareNumerics() Option {
	return func(this *Config) {
		this.specs = append(this.specs, newNumericEqualitySpecification)
	}
}
func (single) CompareTimes() Option {
	return func(this *Config) {
		this.specs = append(this.specs, newTimeEqualitySpecification)
	}
}
func (single) CompareDeep() Option {
	return func(this *Config) {
		this.specs = append(this.specs, newDeepEqualitySpecification)
	}
}
// TODO: CompareEqual (==)
// TODO: CompareFloats (32 and/or 64, almost equal, within tolerance)
// TODO: FormatItem("%+v")
