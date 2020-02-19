package syserr

// close dead  err
type BadRequestErr struct {
	BaseError
}

// create new bad request Error
func NewBadRequestErr(message string) BaseErrorInterface {
	return &BadRequestErr{BaseError: BaseError{message: message, code: badRequest}}
}

// SourceNotFound(
type SourceNotFoundErr struct {
	BaseError
}

// create new not found err
func NewSourceNotFound( message string) BaseErrorInterface {
	return &SourceNotFoundErr{BaseError: BaseError{message: message, code: sourceNotFoundErr}}
}

// Parameter error
type ParameterError struct {
	BaseError
}

// create Parameter error
func NewParameterError( message string) BaseErrorInterface {
	return &ParameterError{BaseError: BaseError{message: message, code: parameterError}}
}
