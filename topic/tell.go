package topic

import (
	"fmt"
	"gitee.com/saltlamp/im-service/abstract"
	"gitee.com/saltlamp/im-service/body"
	"gitee.com/saltlamp/im-service/syserr"
	"github.com/hpcloud/tail"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	DataPath = "./data"
)

func SetDataPath(path string) {
	DataPath = path
}

// tell chan
type TellChan struct {
	// tell file
	// read string lin to chan
	Reader chan *body.Message
	//request connect close notify chan
	//
	Error chan error
}

// tell
type Tell struct {
	//tell config
	TellConfig abstract.TellConfig
}

// create default Tell
func NewDefaultTell(offset int64) *Tell {
	config := abstract.NewDefaultConfig()
	if offset != 0 {
		config.TailConfig.Location = &tail.SeekInfo{Offset: offset}
	}
	return &Tell{TellConfig: config}
}

// tell
func (tell *Tell) TellMessage(tellChan TellChan, request *http.Request) error {
	path, err := tell.TellConfig.FilePathAdapter.Handle(request.RequestURI)
	if err != nil {
		return err
	}
	return tell.TellBind(path, tellChan, request)
}

//
// bind request to tell file
func (tell *Tell) TellBind(path string, tellChan TellChan, request *http.Request) error {
	err := tell.checkAndInitFile(path)
	if err != nil {
		return err
	}
	tail, err := tail.TailFile(path, tell.TellConfig.TailConfig)
	if err != nil {
		return err
	}
	go tell.tellChan(tail, tellChan, request)
	return nil
}

// init file
func (tell *Tell) checkAndInitFile(path string) error {
	_, err := os.Lstat(path)
	// file not found
	// create
	if !os.IsNotExist(err) {
		return nil
	}
	i := strings.LastIndex(path, "/")
	dir := path[:i]
	_, err = os.Lstat(dir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
	}
	if err != nil {
		return err
	}
	file, err := os.Create(path)
	defer file.Close()
	return err
}

func (tell *Tell) tellChan(tail *tail.Tail, tellChan TellChan, request *http.Request) {
	for {
		select {
		case <-request.Context().Done():
			{
				_ = tail.Stop()
				if tellChan.Error != nil {
					tellChan.Error <- syserr.NewCloseErr("request done")
				}
				return
			}
		case <-tail.Dead():
			{
				_ = tail.Stop()
				if tellChan.Error != nil {
					tellChan.Error <- syserr.NewTailDeadErr("tail dead err")
				}
				return
			}
		case lin := <-tail.Lines:
			{
				message, err := tell.TellConfig.ReadAndWriteAdapter.Decoding(lin.Text)
				if err != nil {
					fmt.Println(err)
					continue
				}
				offset, _ := tail.Tell()
				message.Offset = offset
				// get tell offset err
				tellChan.Reader <- message
			}
		case <-time.After(5 * time.Second):
			{
				// time for done check
			}
		}
	}
}
