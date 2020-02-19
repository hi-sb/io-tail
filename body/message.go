package body

import (
	"errors"
	"time"
)

const (
	// text message
	MessageTypeText string = "text/text"
	// voice message
	MessageTypeBase64Voice string = "voice/base64"
	// img message
	MessageTypeBase64Img string = "img/base64"
	// ack message
	MessageTypeAck string = "ack/text"
	// voice message
	MessageTypeUrlVoice string = "voice/url"
	// img message
	MessageTypeUrlImg string = "img/url"
	// sys uid
	sysId string = "00000000000000000000000000000000"
)

// message
type Message struct {
	// form user id
	FormId string
	// send time
	SendTime int64
	// message body
	Body string
	// offset
	Offset int64
	// message type
	ContentType string
}

func (message *Message) GetContentTypeNum() (int, error) {
	return message.GetMessageContentTypeNum(message.ContentType)
}

func (*Message) GetMessageContentTypeNum(contentType string) (int, error) {
	if contentType == "" {
		return -1, errors.New("content type is null")
	}
	switch contentType {
	case MessageTypeText:
		return 0, nil
	case MessageTypeBase64Voice:
		return 1, nil
	case MessageTypeBase64Img:
		return 2, nil
	case MessageTypeAck:
		return 3, nil
	case MessageTypeUrlVoice:
		return 4, nil
	case MessageTypeUrlImg:
		return 5, nil
	default:
		return -1, errors.New("unsupported content type")

	}
}

func NewErrMessage(err error) *Message {
	return &Message{
		FormId:      sysId,
		SendTime:    time.Now().UnixNano() / 1e6,
		Body:        err.Error(),
		Offset:      -1,
		ContentType: MessageTypeText,
	}
}

func (message *Message) CheckContentType(contentType string) error {
	_, err := message.GetMessageContentTypeNum(contentType)
	return err
}
