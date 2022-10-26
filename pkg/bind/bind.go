package bind

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

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
	Validate = validator.New()
	Validate.RegisterValidation("newDate", newDate)
	Validate.RegisterValidation("is-awesome", ValidateMyVal)
}

func ValidateInput(data map[string]interface{}, rules map[string]interface{}) error {
	res := Validate.ValidateMap(data, rules)
	if len(res) > 0 {
		return fmt.Errorf("request field not valid: %v", res)
	}
	return nil
}
