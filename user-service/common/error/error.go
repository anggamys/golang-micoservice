package error

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ValidationResponse struct {
	Field   string `json:"field.omitempty"`
	Message string `json:"message,omitempty"`
}

var ErrValidator = map[string]string{}

func ErrValidationResponse(err error) (validationResponse []ValidationResponse) {
	var fieldErrors validator.ValidationErrors
	if errors.As(err, &fieldErrors) {
		for _, err := range fieldErrors {
			switch err.Tag() {
			case "required":
				validationResponse = append(validationResponse, ValidationResponse{
					Field:   err.Field(),
					Message: fmt.Sprintf("%s is required", err.Field()),
				})

			case "email":
				validationResponse = append(validationResponse, ValidationResponse{
					Field:   err.Field(),
					Message: fmt.Sprintf("%s is not a valid email", err.Field()),
				})

			default:
				ErrValidator, ok := ErrValidator[err.Tag()]

				if ok {
					count := strings.Count(ErrValidator, "%s")
					if count == 1 {
						validationResponse = append(validationResponse, ValidationResponse{
							Field:   err.Field(),
							Message: fmt.Sprintf(ErrValidator, err.Field()),
						})
					} else {
						validationResponse = append(validationResponse, ValidationResponse{
							Field:   err.Field(),
							Message: fmt.Sprintf(ErrValidator, err.Field(), err.Param()),
						})
					}
				} else {
					validationResponse = append(validationResponse, ValidationResponse{
						Field:   err.Field(),
						Message: fmt.Sprintf("Something wrong on %s: %s", err.Field(), err.Tag()),
					})
				}
			}
		}
	}

	return validationResponse
}
