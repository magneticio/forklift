package models

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ValidateDTO func(obj interface{}) error

func NewValidateDTO() ValidateDTO {
	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		jsonTag := fld.Tag.Get("json")
		if jsonTag == "" {
			return fld.Name
		}
		name := strings.SplitN(jsonTag, ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return func(obj interface{}) error {
		if kindOfData(obj) == reflect.Struct {
			if err := validate.Struct(obj); err != nil {
				validationErrors := err.(validator.ValidationErrors)
				var sb strings.Builder
				prefix := ""
				for _, validationError := range validationErrors {
					sb.WriteString(prefix)
					prefix = "\n"
					sb.WriteString(getSingleValidationErrorMessage(validationError))

				}
				return fmt.Errorf(sb.String())

			}
		}
		return nil
	}
}

func kindOfData(data interface{}) reflect.Kind {

	value := reflect.ValueOf(data)
	valueType := value.Kind()

	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	return valueType
}

func getSingleValidationErrorMessage(err validator.FieldError) string {
	trimmedNamespace := ""
	namespace := err.Namespace()
	if idx := strings.IndexByte(namespace, '.'); idx >= 0 {
		trimmedNamespace = namespace[idx+1:]
	}

	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s field is required", trimmedNamespace)
	case "max":
		switch err.Kind() {
		case reflect.String:
			return fmt.Sprintf("%s field must be at most %s character(s)", trimmedNamespace, err.Param())
		case reflect.Array, reflect.Slice, reflect.Map:
			return fmt.Sprintf("%s field must contain at most %s element(s)", trimmedNamespace, err.Param())
		default:
			return fmt.Sprintf("%s field must be at most %s", trimmedNamespace, err.Param())
		}
	case "min":
		switch err.Kind() {
		case reflect.String:
			return fmt.Sprintf("%s field must be at least %s character(s)", trimmedNamespace, err.Param())
		case reflect.Array, reflect.Slice, reflect.Map:
			return fmt.Sprintf("%s field must contain at least %s element(s)", trimmedNamespace, err.Param())
		default:
			return fmt.Sprintf("%s field must be at least %s", trimmedNamespace, err.Param())
		}
	}
	return fmt.Sprintf("%s field is not valid", trimmedNamespace)
}
