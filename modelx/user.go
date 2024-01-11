package modelx

type User struct {
	BaseModel
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
