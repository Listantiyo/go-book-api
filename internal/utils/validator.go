package utils

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New();
}

// ValidateStruct validates a struct based on the `validate` tags
func ValidateStruct(s interface{}) error {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	// Format validation errors become more readable
	var errors []string
	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range errs {
			errors = append(errors, formatValidationError(e))
		}
	}

	return fmt.Errorf("%s", strings.Join(errors, "; "))
}

// formatValidationError formats a single validation error
func formatValidationError(e validator.FieldError) string {
	field := e.Field()

	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", field, e.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", field, e.Param())
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", field, e.Param())
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", field, e.Param())
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}