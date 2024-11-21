package utils

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func ValidationMessages(err error) []string {
	var errorMessages []string
	for _, e := range err.(validator.ValidationErrors) {
		var errorMsg string
		switch e.Tag() {
		case "required":
			errorMsg = fmt.Sprintf("El campo %s es obligatorio", e.Field())
		case "email":
			errorMsg = "el formato no es un email válido"
		case "min":
			errorMsg = fmt.Sprintf("El campo %s debe tener al menos %s caracteres", e.Field(), e.Param())
		default:
			errorMsg = fmt.Sprintf("El campo %s no es válido", e.Field())
		}
		errorMessages = append(errorMessages, errorMsg)
	}

	return errorMessages
}
