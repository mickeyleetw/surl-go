package repositories

import (
	utils "shorten_url/pkg/core"
	database "shorten_url/pkg/core/database"
	domains "shorten_url/pkg/domains"

	"gorm.io/gorm"
)

// GetOrCreateShortURL is a function that gets or creates a short URL
func GetOrCreateShortURL(longURL string) (string, error) {
	var shortURL string
	err := database.WithTransaction(func(tx *gorm.DB) error {
		var url domains.URL
		// Check if the long URL already exists
		result := tx.Where(&domains.URL{LongURL: longURL}).Find(&url)
		if result.RowsAffected > 0 {
			shortURL = url.ShortURL
			return nil
		} else if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
			return result.Error
		}

		// Generate a new short URL
		shortURL = utils.GenerateShortURL(longURL)
		// Check if the generated short URL already exists
		result = tx.Where(&domains.URL{ShortURL: shortURL}).Find(&url)
		if result.RowsAffected > 0 {
			return nil
		} else if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
			return result.Error
		}

		// Create a new URL record
		result = tx.Create(&domains.URL{LongURL: longURL, ShortURL: shortURL})
		if result.Error != nil {
			return result.Error
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	return shortURL, nil
}

// GetLongURLByShortURL is a function that gets a long URL from a short URL
func GetLongURLByShortURL(shortURL string) (string, error) {
	var url domains.URL
	var longURL string = ""
	err := database.WithTransaction(func(tx *gorm.DB) error {
		err := tx.Where(&domains.URL{ShortURL: shortURL}).First(&url).Error
		if err != nil {
			return err
		}
		longURL = url.LongURL
		return nil
	})
	if err != nil {
		return longURL, err
	}

	return longURL, nil
}
