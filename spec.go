package equality

type equalitySpecification interface {
	IsSatisfied() bool
	AreEqual() bool
}
