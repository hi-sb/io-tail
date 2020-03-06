package body

import (
	"errors"
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
	// sys notify message
	MessageTypeNotify string = "notify/text"
	// Add friends
	MessageTypeAddFriends string = "add-friends/json"
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
		return -1, errors.New("ContentType 不能为空")
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
	case MessageTypeNotify:
		return 6, nil
	case MessageTypeAddFriends:
		return 7, nil
	default:
		return -1, errors.New("非法的 ContentType")
	}
}

func (message *Message) CheckContentType(contentType string) error {
	_, err := message.GetMessageContentTypeNum(contentType)
	return err
}
