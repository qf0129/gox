package dbx

import (
	"github.com/qf0129/gox/pkg/timex"
	"github.com/rs/xid"
	"gorm.io/gorm"
)

type BaseModel struct {
	Id    int64          `gorm:"primaryKey" json:"-"`
	Ctime *timex.Time    `gorm:"autoCreateTime;type:datetime(3)"`
	Utime *timex.Time    `gorm:"autoUpdateTime;type:datetime(3)"`
	Dtime gorm.DeletedAt `gorm:"index;type:datetime(3)" json:"-"`
}

type BaseUidModel struct {
	Id    string         `gorm:"primaryKey;type:varchar(64)"`
	Ctime *timex.Time    `gorm:"autoCreateTime;type:datetime(3)"`
	Utime *timex.Time    `gorm:"autoUpdateTime;type:datetime(3)"`
	Dtime gorm.DeletedAt `gorm:"index;type:datetime(3)" json:"-"`
}

func (m *BaseUidModel) BeforeCreate(tx *gorm.DB) error {
	if m.Id == "" {
		m.Id = xid.New().String()
	}
	return nil
}
