package domains

import (
	"gorm.io/gorm"
)

type Url struct {
	gorm.Model
	ShortUrl string `gorm:"unique"`
	LongUrl  string
}

func (Url) TableName() string {
	return "url"
}
