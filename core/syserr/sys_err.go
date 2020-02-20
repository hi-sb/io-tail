package syserr

// sys err
type SysErr struct {
	BaseError
}

// create new sys err
func NewSysErr(message string) BaseErrorInterface {
	return &SysErr{BaseError: BaseError{message: message, code: sysErr}}
}

