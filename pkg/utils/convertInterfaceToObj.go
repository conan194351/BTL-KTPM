package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func ConvertToStruct[T any](data interface{}) (T, error) {
	var result T
	if data == nil {
		return result, fmt.Errorf("input data is nil")
	}
	if reflect.TypeOf(data) == reflect.TypeOf(result) {
		var ok bool
		result, ok = data.(T)
		if ok {
			return result, nil
		}
	}
	var jsonBytes []byte
	var err error

	if str, ok := data.(string); ok {
		jsonBytes = []byte(str)
	} else {
		jsonBytes, err = json.Marshal(data)
		if err != nil {
			return result, fmt.Errorf("failed to marshal data: %w", err)
		}
	}
	err = json.Unmarshal(jsonBytes, &result)
	if err != nil {
		return result, fmt.Errorf("failed to unmarshal to target type: %w", err)
	}

	return result, nil
}
