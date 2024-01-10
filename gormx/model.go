package gormx

import (
	"time"

	"github.com/rs/xid"
	"gorm.io/gorm"
)

type BaseModel struct {
	Id    string          `gorm:"primaryKey;type:varchar(50);" json:"id"`
	Ctime *time.Time      `gorm:"autoCreateTime;comment:'CreatedTime'" json:"ctime"`
	Utime *time.Time      `gorm:"autoUpdateTime;comment:'UpdatedTime'" json:"utime"`
	Dtime *gorm.DeletedAt `gorm:"index;comment:'DeletedTime'" json:"-" `
}

func (m *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	if m.Id == "" {
		m.Id = xid.New().String()
	}
	return
}
