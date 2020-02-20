package abstract

import (
	"github.com/hi-sb/io-tail/body"
	"github.com/hpcloud/tail"
)

// file path
// is file net source to sys path
type FilePathAdapter interface {
	// to sys path
	// return sys path
	Handle(uri string) (string, error)
}

//read and write
// Encoding and decoding
type ReadAndWriteAdapter interface {
	//Encoding
	Encoding(body *body.Message) (string, error)
	//Decoding
	Decoding(body string) (*body.Message, error)
}

type TellConfig struct {
	// path
	FilePathAdapter FilePathAdapter
	// tail file
	//config
	// is file offset or reopen .....
	TailConfig tail.Config
	// encode and decode
	ReadAndWriteAdapter ReadAndWriteAdapter
}
