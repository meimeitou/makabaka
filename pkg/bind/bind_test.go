package bind

import (
	"fmt"
	"testing"
)

type myStruct struct {
	String string `validate:"is-awesome"`
}

func TestBind(t *testing.T) {
	s := myStruct{String: "awesome"}

	err := validate.Struct(s)
	if err != nil {
		fmt.Printf("Err(s):\n%+v\n", err)
	}
	s.String = "not awesome"
	err = validate.Struct(s)
	if err != nil {
		fmt.Printf("Err(s):\n%+v\n", err)
	}

	user := map[string]interface{}{"name": "Arshiya Kiani", "email": "s", "age": 13}

	rules := map[string]interface{}{"name": "required,min=8,max=32", "email": "gt=0", "age": "number,gt=14", "ss": ""}

	errs := validate.ValidateMap(user, rules)

	if len(errs) > 0 {
		fmt.Println(errs)
	}
}
