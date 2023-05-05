package bind

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

var newDate validator.Func = func(fl validator.FieldLevel) bool {
	date, ok := fl.Field().Interface().(time.Time)
	if ok {
		today := time.Now()
		if today.After(date) {
			return false
		}
	}
	return true
}

func ValidateMyVal(fl validator.FieldLevel) bool {
	return fl.Field().String() == "awesome"
}

func init() {
	validate = validator.New()
	validate.RegisterValidation("newDate", newDate)
	validate.RegisterValidation("is-awesome", ValidateMyVal)
}

// custom validation
func RegisterValidation(tag string, fn validator.Func, callValidationEvenIfNull ...bool) {
	validate.RegisterValidation(tag, fn, callValidationEvenIfNull...)
}

func ValidateInput(data map[string]interface{}, rules map[string]interface{}) error {
	res := validate.ValidateMap(data, rules)
	if len(res) > 0 {
		return fmt.Errorf("request field not valid: %v", res)
	}
	return nil
}
