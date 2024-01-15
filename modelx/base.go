package modelx

import (
	"github.com/rs/xid"
	"gorm.io/gorm"
)

type BaseModel struct {
	Id    string          `gorm:"primaryKey;type:varchar(50);" json:"id"`
	Ctime *Time           `gorm:"autoCreateTime;comment:'CreatedTime'" json:"ctime"`
	Utime *Time           `gorm:"autoUpdateTime;comment:'UpdatedTime'" json:"utime"`
	Dtime *gorm.DeletedAt `gorm:"index;comment:'DeletedTime'" json:"-" `
}

type BaseModelSimple struct {
	Id    string `gorm:"primaryKey;type:varchar(50);" json:"id"`
	Ctime *Time  `gorm:"autoCreateTime;comment:'CreatedTime'" json:"ctime"`
}

func (m *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	if m.Id == "" {
		m.Id = xid.New().String()
	}
	return
}

func (m *BaseModelSimple) BeforeCreate(tx *gorm.DB) (err error) {
	if m.Id == "" {
		m.Id = xid.New().String()
	}
	return
}
