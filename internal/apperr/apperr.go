package apperr

// Public is an error whose message is safe to show to API clients.
type Public struct {
	msg string
}

func New(msg string) *Public {
	return &Public{msg: msg}
}

func (e *Public) Error() string {
	return e.msg
}
