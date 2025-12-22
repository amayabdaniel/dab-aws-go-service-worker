package middleware

import (
	"net/http"
	"strings"
	
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// ValidationError formats validation errors for API responses
func ValidationError(c *gin.Context, err error) {
	if ve, ok := err.(validator.ValidationErrors); ok {
		errors := make(map[string]string)
		for _, e := range ve {
			field := strings.ToLower(e.Field())
			switch e.Tag() {
			case "required":
				errors[field] = field + " is required"
			case "min":
				errors[field] = field + " must be at least " + e.Param() + " characters"
			case "max":
				errors[field] = field + " must be at most " + e.Param() + " characters"
			default:
				errors[field] = field + " is invalid"
			}
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "validation failed",
			"fields": errors,
		})
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
}

// ValidateStruct validates any struct with validation tags
func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}