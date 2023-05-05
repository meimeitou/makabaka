package tpl

import (
	"bytes"
	"fmt"
	"testing"
)

func TestTpl(t *testing.T) {
	tpl, err := temp.Parse(`select * from {{ toType .DB "string" }}`)
	if err != nil {
		panic(err)
	}
	b := &bytes.Buffer{}
	err = tpl.Execute(b, map[string]interface{}{
		"DB": "test",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(b.String())

	tpl, err = temp.Parse(`select * from {{ toType .DB "string" }}`)
	if err != nil {
		panic(err)
	}
	b = &bytes.Buffer{}
	err = tpl.Execute(b, map[string]interface{}{
		"DB": "xxx",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(b.String())
}
