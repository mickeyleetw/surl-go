package models

import "fmt"

type CreateUrlShortenModel struct {
	LongUrl string `json:"long_url" binding:"required"`
}

// Validate method to implement the RequestBody interface
func (m *CreateUrlShortenModel) Validate() error {
	if m.LongUrl == "" {
		return fmt.Errorf("LongUrl cannot be empty")
	}
	return nil
}

type RetrieveUrlShortenModel struct {
	ShortUrl string `json:"shortened_url"`
}

// Validate method to implement the RequestBody interface
func (m *RetrieveUrlShortenModel) Validate() error {
	if m.ShortUrl == "" {
		return fmt.Errorf("ShortUrl cannot be empty")
	}
	return nil
}
