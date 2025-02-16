package models

import "fmt"

// CreateURLShortenModel is a struct that contains the long URL
type CreateURLShortenModel struct {
	LongURL string `json:"long_url" binding:"required"`
}

// Validate method to implement the RequestBody interface
func (m *CreateURLShortenModel) Validate() error {
	if m.LongURL == "" {
		return fmt.Errorf("LongUrl cannot be empty")
	}
	return nil
}

// RetrieveURLShortenModel is a struct that contains the short URL
type RetrieveURLShortenModel struct {
	ShortURL string `json:"shortened_url"`
}

// Validate method to implement the RequestBody interface
func (m *RetrieveURLShortenModel) Validate() error {
	if m.ShortURL == "" {
		return fmt.Errorf("ShortUrl cannot be empty")
	}
	return nil
}
