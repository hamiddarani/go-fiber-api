package validations

import (
	"log"
	"regexp"
	"sync"

	"github.com/go-playground/validator/v10"
)

var once sync.Once
var validate *validator.Validate

func NewValidator() *validator.Validate {
	once.Do(func() {
		validate = validator.New()
	})

	return validate
}

const pattern string = `^09(1[0-9]|2[0-2]|3[0-9]|9[0-9])[0-9]{7}$`

func IranianMobileNumberValidator(fld validator.FieldLevel) bool {
	value, ok := fld.Field().Interface().(string)
	if !ok {
		return false
	}

	res, err := regexp.MatchString(pattern, value)
	if err != nil {
		log.Print(err.Error())
	}

	return res

}
