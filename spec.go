package equality

type Specification interface {
	IsSatisfied() bool
	AreEqual() bool
}
