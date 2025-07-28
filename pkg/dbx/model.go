package dbx

import (
	"github.com/qf0129/gox/pkg/timex"
	"github.com/rs/xid"
	"gorm.io/gorm"
)

type BaseModel struct {
	Id    int64          `gorm:"primaryKey"`
	Ctime *timex.Time    `gorm:"autoCreateTime;type:datetime"`
	Mtime *timex.Time    `gorm:"autoUpdateTime;type:datetime"`
	Dtime gorm.DeletedAt `gorm:"index;type:datetime" json:"-"`
}

type BaseUidModel struct {
	Id    string         `gorm:"primaryKey;type:varchar(50)"`
	Ctime *timex.Time    `gorm:"autoCreateTime;type:datetime"`
	Mtime *timex.Time    `gorm:"autoUpdateTime;type:datetime"`
	Dtime gorm.DeletedAt `gorm:"index;type:datetime" json:"-"`
}

func (m *BaseUidModel) BeforeCreate(tx *gorm.DB) error {
	if m.Id == "" {
		m.Id = xid.New().String()
	}
	return nil
}

type BaseUid struct {
	Uid string `gorm:"index;type:varchar(50)"`
}

func (m *BaseUid) BeforeCreate(tx *gorm.DB) error {
	if m.Uid == "" {
		m.Uid = xid.New().String()
	}
	return nil
}
