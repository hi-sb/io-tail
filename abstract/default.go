package abstract

import (
	"encoding/json"
	"errors"
	"gitee.com/saltlamp/im-service/body"
	"gitee.com/saltlamp/im-service/config"
	"gitee.com/saltlamp/im-service/utils"
	"github.com/hpcloud/tail"
)


func NewDefaultConfigByTailConfig(config tail.Config) TellConfig {
	filePathAdapter := NewDefaultFilePathAdapter()
	readAndWriteHandle := NewDefaultReadAndWriteAdapter()
	return TellConfig{
		TailConfig:          config,
		FilePathAdapter:     filePathAdapter,
		ReadAndWriteAdapter: readAndWriteHandle,
	}
}

func NewDefaultConfig() TellConfig {
	tailConfig := tail.Config{Follow: true, ReOpen: true, Poll: true}
	return NewDefaultConfigByTailConfig(tailConfig)
}

// default file path adapter
type DefaultFilePathAdapter struct {
	FilePathAdapter
}

//new
func NewDefaultFilePathAdapter() FilePathAdapter {
	return &DefaultFilePathAdapter{}
}

// source path  to sys path
func (defaultFilePathAdapter *DefaultFilePathAdapter) Handle(uri string) (string, error) {
	if uri == "" || uri == "/" {
		return "", errors.New("path is null")
	}
	return defaultFilePathAdapter.hashPath(uri)
}

func (*DefaultFilePathAdapter) hashPath(uri string) (string, error) {
	md5 := utils.Md5V2(uri)
	return config.DataPath + "/" + md5[0:1] + "/" + md5[2:3] + "/" + md5[4:5] + "/" + md5[6:7] + "/" + md5, nil
}

// default ReadAndWriteAdapter
type DefaultReadAndWriteAdapter struct {
	ReadAndWriteAdapter
}

func NewDefaultReadAndWriteAdapter() ReadAndWriteAdapter {
	return &DefaultReadAndWriteAdapter{}
}

//
func (*DefaultReadAndWriteAdapter) Encoding(body *body.Message) (string, error) {
	byte, err := json.Marshal(body)
	if err != nil {
		return "", err
	}
	return string(byte) + "\n", err
}

//
func (*DefaultReadAndWriteAdapter) Decoding(requestBody string) (*body.Message, error) {
	message := new(body.Message)
	err := json.Unmarshal([]byte(requestBody), message)
	if err != nil {
		return nil, err
	}
	return message, nil
}
