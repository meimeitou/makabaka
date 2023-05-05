package tpl

import (
	"bytes"
	"fmt"
	"text/template"
)

var (
	temp *template.Template
)

func init() {
	temp = template.New("").Funcs(funcs)
}

func RegisterFunction(funcMap template.FuncMap) {
	temp = temp.Funcs(funcMap)
}

func Parse(data string, values map[string]interface{}) (string, error) {
	p, err := temp.Parse(data)
	if err != nil {
		return "", fmt.Errorf("template err: %v", err)
	}
	b := &bytes.Buffer{}
	err = p.Execute(b, values)
	if err != nil {
		return "", fmt.Errorf("template value error: %v", err)
	}
	return b.String(), nil
}
