package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type RequestPayloadValidator struct {
	validator *validator.Validate
}

func (rv *RequestPayloadValidator) Validate(obj interface{}) error {
	return rv.validator.Struct(obj)
}

func ValidateRequest[T any]() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req T
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		requestValidator := validator.New()
		if err := requestValidator.Struct(req); err != nil {
			validationErrors := err.(validator.ValidationErrors)
			c.JSON(http.StatusUnprocessableEntity, gin.H{"validation_error": validationErrors.Error()})
			c.Abort()
			return
		}

		c.Set("validated_data", req)
		c.Next()
	}
}
