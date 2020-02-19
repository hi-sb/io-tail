package syserr
// ContentType err
type ContentTypeErr struct {
	BaseError
}

// create content type err
func NewContentTypeErr(message string) BaseErrorInterface {
	return &ContentTypeErr{BaseError: BaseError{message: message, code: contentTypeErr}}
}
