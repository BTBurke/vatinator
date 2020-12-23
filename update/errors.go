package update

// NonFatalError is returned when an error doesn't affect the running binary during the update process
type NonFatalError struct {
	Err error
}

func (e NonFatalError) Error() string {
	return e.Err.Error()
}

// FatalError is returned when an error during update affects the running binary and a new version should be downloaded
type FatalError struct {
	Err error
}

func (e FatalError) Error() string {
	return e.Err.Error()
}

var _ error = NonFatalError{}
var _ error = FatalError{}
