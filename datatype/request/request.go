package request

type Request interface {
	Validation() bool
}
