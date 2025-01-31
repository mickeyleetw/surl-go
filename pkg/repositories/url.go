package repositories

import (
	utils "shorten_url/pkg/core"
	database "shorten_url/pkg/core/database"
	domains "shorten_url/pkg/domains"

	"gorm.io/gorm"
)

func GetOrCreateShortUrl(longUrl string) (string, error) {
	var shortUrl string
	err := database.WithTransaction(func(tx *gorm.DB) error {
		var url domains.Url
		// Check if the long URL already exists
		err := tx.Where(&domains.Url{LongUrl: longUrl}).First(&url).Error
		if err == nil {
			shortUrl = url.ShortUrl
			return nil
		} else if err != gorm.ErrRecordNotFound {
			return err
		}

		// Generate a new short URL
		shortUrl = utils.GenerateShortURL(longUrl)
		// Check if the generated short URL already exists
		err = tx.Where(&domains.Url{ShortUrl: shortUrl}).First(&url).Error
		if err == nil {
			return nil
		} else if err != gorm.ErrRecordNotFound {
			return err
		}

		// Create a new URL record
		err = tx.Create(&domains.Url{LongUrl: longUrl, ShortUrl: shortUrl}).Error
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	return shortUrl, nil
}

func GetUrlByShortUrl(shortUrl string) (string, error) {
	var url domains.Url
	var longUrl string = ""
	err := database.WithTransaction(func(tx *gorm.DB) error {
		err := tx.Where(&domains.Url{ShortUrl: shortUrl}).First(&url).Error
		if err != nil {
			return err
		}
		longUrl = url.LongUrl
		return nil
	})
	if err != nil {
		return longUrl, err
	}

	return longUrl, nil
}
