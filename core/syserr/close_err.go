package syserr

// close dead  err
type CloseErr struct {
	BaseError
}

// create newAuthError
func NewCloseErr(message string) BaseErrorInterface {
	return &CloseErr{BaseError: BaseError{message: message, code: closeErr}}
}
