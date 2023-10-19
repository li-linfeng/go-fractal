package fractal

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

// 使用反射检查是否为切片
func IsSlice(v interface{}) bool {
	return reflect.TypeOf(v).Kind() == reflect.Slice
}

// 是否为map
func IsMap(v interface{}) bool {
	return reflect.TypeOf(v).Kind() == reflect.Map
}

func IsStruct(v interface{}) bool {
	return reflect.TypeOf(v).Kind() == reflect.Struct
}

func ConvertToInterfaceSlice(input interface{}) ([]interface{}, error) {
	v := reflect.ValueOf(input)
	// 判断输入是否为切片
	if v.Kind() != reflect.Slice {
		return nil, fmt.Errorf("input is not a slice: %T", input)
	}
	result := make([]interface{}, v.Len())

	for i := 0; i < v.Len(); i++ {
		result[i] = v.Index(i).Interface()
	}

	return result, nil
}

func ToJson(v interface{}) (string, error) {
	jsonData, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func BuildNestedMap(paths []string) map[string]interface{} {
	result := make(map[string]interface{})

	for _, path := range paths {
		parts := strings.Split(path, ".")
		currentMap := result

		for _, part := range parts {
			if _, exists := currentMap[part]; !exists {
				currentMap[part] = make(map[string]interface{})
			}
			currentMap = currentMap[part].(map[string]interface{})
		}
	}

	return result
}

// 是否为基本类型或基本类型切片, 或者nil
func isBasicOrBasicSlice(i interface{}) bool {
	if i == nil {
		return true
	}
	switch i.(type) {
	case int:
		return true
	case int8:
		return true
	case int16:
		return true
	case int32:
		return true
	case int64:
		return true
	case uint:
		return true
	case uint8:
		return true
	case uint16:
		return true
	case uint32:
		return true
	case uint64:
		return true
	case uintptr:
		return true
	case float32:
		return true
	case float64:
		return true
	case complex64:
		return true
	case complex128:
		return true
	case string:
		return true
	case bool:
		return true
	case []int:
		return true
	case []int8:
		return true
	case []int16:
		return true
	case []int32:
		return true
	case []int64:
		return true
	case []uint:
		return true
	case []uint8:
		return true
	case []uint16:
		return true
	case []uint32:
		return true
	case []uint64:
		return true
	case []uintptr:
		return true
	case []float32:
		return true
	case []float64:
		return true
	case []complex64:
		return true
	case []complex128:
		return true
	case []string:
		return true
	case []bool:
		return true
	default:
		return false
	}
}
