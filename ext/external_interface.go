package ext

import "github.com/hi-sb/io-tail/auth"

// source
type Source struct {
	//
	// open source name
	// AtNum : name@xxx.ooo:p
	Name string
	// rsa public key
	PublicKey string
	// Nickname
	Nickname string
	//
	ProfilePhotoUrl string
}

//
type OpenSource struct {
	Source
	//
	SourceType OpenSourceType
	//
	Describe string
	// value true
	IsOpenSource bool
	// create name
	CreateName string
}

// source type
type OpenSourceType int

const (
	//Subscription
	// Only creators are allowed to send messages
	SubscriptionOpenSourceType OpenSourceType = 1
	//Interactive
	//Allow everyone to send messages
	InteractiveOpenSourceType OpenSourceType = 2
)

// external_interface
// Interface definitions that allow external extensions.
type ExternalInterface interface {
	// We allow external management of our user systems.
	//Here we only need an interface to get the public key by user name. Therefore your username must be unique.
	GetUserPublicKey(name string) (string, error)
	//Check whether a user has write access to a resource, that is,
	//whether messages can be sent to an identity
	CheckWritePermission(jwt *auth.JWT, name string) error
	// create open source
	// If public key is null .
	// The content will not be encrypted when it is stored, and no password is required when it is join
	//
	CreateOpenSource(openSource *OpenSource) (bool, error)
	// get base data
	// get source ( source or open source ) return source data json string
	GetSourceBaseData(name string) (interface{}, error)
	// get private or public source rsa
	// public key
	GetRsaPublicKey( name string) (*string,error)
}
