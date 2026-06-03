package dbx

import (
	"strings"

	"github.com/google/uuid"
	"github.com/qf0129/gox/pkg/timex"
	"github.com/rs/xid"
	"gorm.io/gorm"
)

type BaseModel struct {
	Id        int64          `gorm:"primaryKey" json:"-"`
	CreatedAt *timex.Time    `gorm:"autoCreateTime;type:datetime"`
	UpdatedAt *timex.Time    `gorm:"autoUpdateTime;type:datetime"`
	DeletedAt gorm.DeletedAt `gorm:"index;type:datetime" json:"-"`
}

type BaseUidModel struct {
	Id        string         `gorm:"primaryKey;type:varchar(50)"`
	CreatedAt *timex.Time    `gorm:"autoCreateTime;type:datetime"`
	UpdatedAt *timex.Time    `gorm:"autoUpdateTime;type:datetime"`
	DeletedAt gorm.DeletedAt `gorm:"index;type:datetime" json:"-"`
}

func (m *BaseUidModel) BeforeCreate(tx *gorm.DB) error {
	if m.Id == "" {
		m.Id = strings.ReplaceAll(uuid.New().String(), "-", "")
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
