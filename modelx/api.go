package modelx

type Api struct {
	BaseModel
	Group  string `gorm:"index;type:varchar(50);" json:"group"`
	Key    string `gorm:"index;type:varchar(200);" json:"key"`
	Method string `gorm:"index;type:varchar(20);" json:"method"`
	Path   string `gorm:"type:varchar(200);" json:"path"`
	// Func        string `gorm:"type:varchar(200);" json:"func"`
	Description string `gorm:"type:varchar(500);" json:"description"`
}
