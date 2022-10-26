package exec

import (
	"errors"
	"text/template"
)

var (
	toType = template.FuncMap{
		"toType": func(data interface{}, target string) (interface{}, error) {
			return TryConvert(data, target)
		},
	}
)

func TryConvert(val interface{}, target string) (interface{}, error) {
	// logx.Info(reflect.ValueOf(val).Kind())
	switch target {
	case "string":
		v, ok := val.(string)
		if !ok {
			return nil, errors.New("类型错误")
		}
		return v, nil
	case "int":
		switch v := val.(type) {
		case int:
			return int64(v), nil
		case int32:
			return int64(v), nil
		case int64:
			return v, nil
		case float32:
			return int64(v), nil
		case float64:
			return int64(v), nil
		default:
			return nil, errors.New("类型错误")
		}
	case "float":
		switch v := val.(type) {
		case int:
			return float64(v), nil
		case int32:
			return float64(v), nil
		case int64:
			return float64(v), nil
		case float32:
			return float64(v), nil
		case float64:
			return v, nil
		default:
			return nil, errors.New("类型错误")
		}
	}
	return nil, errors.New("不支持的类型")
}
