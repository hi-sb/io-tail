package source


const (
	// user source
	privateSource = "private_source"
	// open source
	publicSource = "public_source"
)

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