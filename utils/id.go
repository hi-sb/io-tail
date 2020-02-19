package utils

import (
	"github.com/satori/go.uuid"
	"strings"
)

func GetID() string {
	return strings.ReplaceAll(uuid.NewV4().String(), "-", "")
}
