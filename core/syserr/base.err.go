package syserr

// sys base error
type BaseError struct {
	// error message
	message string
	// error code
	code int
}

// base error interface
type BaseErrorInterface interface {
	error
	Code() int
}

func (err *BaseError) Error() string {
	return err.message
}

func (err *BaseError) Code() int {
	return err.code
}
