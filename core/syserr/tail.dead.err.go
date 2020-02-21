package syserr

// tail dead  err
type TailDeadErr struct {
	BaseError
}

// create TailDeadErr
func NewTailDeadErr(message string) BaseErrorInterface {
	return &TailDeadErr{BaseError: BaseError{message: message, code: tailDeadErr}}
}
