package equality

type config struct {
	specs  []func(a, b interface{}) Specification
	format func(interface{}) string
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
