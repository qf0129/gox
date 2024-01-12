package tmplx

import (
	"github.com/qf0129/gox/modelx"
	"golang.org/x/crypto/bcrypt"
)

type Api struct {
	modelx.BaseModel
	Group  string `gorm:"index;type:varchar(50);" json:"group"`
	Key    string `gorm:"index;type:varchar(200);" json:"key"`
	Method string `gorm:"index;type:varchar(20);" json:"method"`
	Path   string `gorm:"type:varchar(200);" json:"path"`
	// Func        string `gorm:"type:varchar(200);" json:"func"`
	Description string `gorm:"type:varchar(500);" json:"description"`
}

type Role struct {
	modelx.BaseModel
	Name        string `gorm:"index;type:varchar(100);" json:"name"`
	Description string `gorm:"type:varchar(500);" json:"description"`
}

type User struct {
	modelx.BaseModel
	Username string `gorm:"type:varchar(200);index;not null" json:"username"`
	Nickname string `gorm:"type:varchar(200)" json:"nickname"`
	Realname string `gorm:"type:varchar(200)" json:"realname"`
	Password string `gorm:"type:varchar(500)" json:"-"`
	Phone    string `gorm:"type:varchar(50);index" json:"phone"`
	Mail     string `gorm:"type:varchar(500);index" json:"mail"`
	Source   string `gorm:"type:varchar(50)" json:"source"`
	Disabled bool   `gorm:"default:false;" json:"disabled"`
	Online   bool   `gorm:"default:false;" json:"online"`

	Roles []*Role `gorm:"many2many:user_role;" json:"roles"`
}

func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return false
	}
	return true
}
