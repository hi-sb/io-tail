package syserr

// jwt auth err
type AuthError struct {
	BaseError
}

// create newAuthError
// token check err
func NewTokenAuthError(message string) BaseErrorInterface {
	return &AuthError{BaseError: BaseError{message: message, code: authTokenErr}}
}

// create newAuthError
// unbound user
func NewUnboundAuthError(message string) BaseErrorInterface {
	return &AuthError{BaseError: BaseError{message: message, code: authUnboundUserErr}}
}

// create newAuthError
// id is null
func NewIdIsNullAuthError(message string) BaseErrorInterface {
	return &AuthError{BaseError: BaseError{message: message, code: authIdIsNullErr}}
}

// create newAuthError
// public key err
func NewPublicKeyAuthError(message string) BaseErrorInterface {
	return &AuthError{BaseError: BaseError{message: message, code: authPublicKeyErr}}
}

// create newAuthError
// check randomKey err
func NewCheckRandomKeyError(message string) BaseErrorInterface {
	return &AuthError{BaseError: BaseError{message: message, code: checkRandomKeyError}}
}

// create newAuthError
// federation auth err
func NewFederationAuthErr(message string) BaseErrorInterface {
	return &AuthError{BaseError: BaseError{message: message, code: federationAuthErr}}
}

//create newAuthError
// permission check err
func NewPermissionErr(message string) BaseErrorInterface {
	return &AuthError{BaseError: BaseError{message: message, code: federationAuthErr}}
}
