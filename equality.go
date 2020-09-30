// Package equality facilitates comparisons of any two values.
package equality

// Report compares expected and actual and returns
// a full report of any discrepancy between them.
func Report(expected, actual interface{}, options ...Option) string {
	_, report := Compare(expected, actual, options...)
	return report
}

// Compare returns a comparison of expected and actual as well as
// a full report of any discrepancy between them.
func Compare(expected, actual interface{}, options ...Option) (ok bool, report string) {
	ok = Check(expected, actual, options...)
	if !ok {
		report = newFormatter(expected, actual, options...).String()
	}
	return ok, report
}

// Check returns a comparison of expected and actual according
// to the specifications defined in this package.
func Check(expected, actual interface{}, options ...Option) bool {
	config := new(config)
	config.apply(options...)
	config.applyDefaultEqualitySpecs()

	for _, factory := range config.specs {
		spec := factory(expected, actual)
		if !spec.IsSatisfied() {
			continue
		}
		if spec.AreEqual() {
			return true
		}
		break
	}
	return false
}
