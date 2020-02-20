package db
import (
	uuid "github.com/satori/go.uuid"
	"reflect"
	"strings"
	"time"
)

var BaseModelType = reflect.TypeOf(BaseModel{})

// base Model
type BaseModel struct {
	ID        string `gorm:"type:varchar(32);primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type BaseModelInterface interface {
	GetID() string
}

func (baseModel *BaseModel) GetID() string {
	return baseModel.ID
}

func (baseModel *BaseModel) Bind() {
	if baseModel.CreatedAt.Unix() <= 0 {
		baseModel.CreatedAt = time.Now()
	}
	if baseModel.UpdatedAt.Unix() <= 0 {
		baseModel.UpdatedAt = time.Now()
	}
	if baseModel.ID == "" {
		baseModel.ID = strings.ReplaceAll(uuid.NewV4().String(), "-", "")
	}
}