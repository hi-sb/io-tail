package syserr

type ServiceErr struct {
	BaseError
}


// 业务错误
func NewServiceError(message string) BaseErrorInterface {
	return &ServiceErr{BaseError: BaseError{message: message, code: serviceErr}}
}