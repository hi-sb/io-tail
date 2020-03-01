package model

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


// send request
type SendRequest struct {
	// send time
	SendTime int64
	// message body
	Body string
	// message type
	ContentType string
}