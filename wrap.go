package panik

type Wrapable interface {
	error
	Wrap(error)
}
