package utils

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

var ugandaNINRegex = regexp.MustCompile(`^C[MF]\d{2}[A-Za-z0-9]{10}$`)
var dhis2UID = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9]{10}$`)
var yesNo = regexp.MustCompile(`^(Yes|No)$`)
var maleFemale = regexp.MustCompile(`^(Male|Female)$`)

// UgandaNINValidation is the custom validator function.
func UgandaNINValidation(fl validator.FieldLevel) bool {
	nin, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	// Allow empty string since NIN is optional.
	if nin == "" {
		return true
	}
	return ugandaNINRegex.MatchString(nin)
}

// Dhis2UIDValidation is the custom validator function.
func Dhis2UIDValidation(fl validator.FieldLevel) bool {
	uid, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	return dhis2UID.MatchString(uid)
}

// YesNoValidation is the custom validator function.
func YesNoValidation(fl validator.FieldLevel) bool {
	yesNoValue, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	return yesNo.MatchString(yesNoValue)
}

// MaleFemaleValidation ... is the custom validator function.
func MaleFemaleValidation(fl validator.FieldLevel) bool {
	maleFemaleValue, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	return maleFemale.MatchString(maleFemaleValue)
}
