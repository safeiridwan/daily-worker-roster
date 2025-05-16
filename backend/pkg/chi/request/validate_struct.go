package request

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

// init initializes the validator and registers the custom tag name function.
func init() {
	validate = validator.New()
	validate.RegisterTagNameFunc(getFieldName)
}

// ValidateStruct validates the given struct using the validator.
// It returns an error if the struct fails validation.
func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

// getFieldName returns the field name to be used for validation.
// It checks for custom JSON tags and nested JSON tags.
func getFieldName(fld reflect.StructField) string {
	name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
	if name == "-" {
		return ""
	}
	if name2 := strings.SplitN(fld.Tag.Get("json_nested"), ",", 2)[0]; name2 != "" {
		return name2
	}
	return name
}
