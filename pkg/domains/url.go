package domains

import (
	"gorm.io/gorm"
)

// URL is a struct that contains the URL
type URL struct {
	gorm.Model
	ShortURL string `gorm:"unique;column:short_url"`
	LongURL  string `gorm:"column:long_url"`
}

// TableName is a function that returns the table name
func (URL) TableName() string {
	return "url"
}
