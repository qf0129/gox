package modelx

type Role struct {
	BaseModel
	Name        string `gorm:"index;type:varchar(100);" json:"name"`
	Description string `gorm:"type:varchar(500);" json:"description"`
}
