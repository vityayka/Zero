package errors

func Public(err error, msg string) error {
	return publicError{err, msg}
}

type publicError struct {
	err error
	msg string
}

func (pubErr publicError) Error() string {
	return pubErr.err.Error()
}

func (pubErr publicError) String() string {
	return pubErr.msg
}

func (pubErr publicError) Unwrap() error {
	return pubErr.err
}
