package utils

import (
	"errors"
	"reflect"
	"strconv"
)

func ToInt(val interface{}) (int, error) {
	if reflect.TypeOf(val).Kind() == reflect.Float64 {
		asFloat, ok := val.(float64)
		if !ok {
			return 0, errors.New("not an int")
		}
		if asFloat != float64(int(asFloat)) {
			return 0, errors.New("not an int")
		}
		return int(asFloat), nil
	}
	asInt, ok := val.(int)
	if ok {
		return asInt, nil
	}
	asString, ok := val.(string)
	if ok {
		asInt, err := strconv.Atoi(asString)
		if err != nil {
			return 0, errors.New("not an int")
		}
		return asInt, nil
	}
	return 0, errors.New("not an int")
}

func ToBool(val interface{}) (bool, error) {
	if reflect.TypeOf(val).Kind() == reflect.String {
		asString, ok := val.(string)
		if ok {
			asBool, err := strconv.ParseBool(asString)
			if err != nil {
				return false, errors.New("not a bool")
			}
			return asBool, nil
		}
	}
	asBool, ok := val.(bool)
	if !ok {
		return false, errors.New("not a bool")
	}
	return asBool, nil

}
