package ext

import "github.com/hi-sb/io-tail/auth"

//
type OpenSource struct {
	//
	// open source name
	Name string
	//
	ProfilePhotoUrl string
	//
	Describe string
	// create name
	CreateName string
}


// external_interface
// Interface definitions that allow external extensions.
type ExternalInterface interface {
	//Check whether a user has write access to a resource, that is,
	//whether messages can be sent to an identity
	CheckWritePermission(jwt *auth.JWT, name string) error
	// create open source
	//
	CreateOpenSource(openSource *OpenSource) (bool, error)
}
