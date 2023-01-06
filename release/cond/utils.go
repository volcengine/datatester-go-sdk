/**
 * Apache 2.0
 * Copyright 2022 Beijing Volcano Engine Technology Co., Ltd.
 */

package cond

import (
	"errors"
	"reflect"
)

// ToFloat64 attempts to convert the given value to a float64
func ToFloat64(value interface{}) (float64, error) {
	if value == nil {
		return 0, errors.New("cond value is nil")
	}
	var floatType = reflect.TypeOf(float64(0))
	v := reflect.ValueOf(value)
	v = reflect.Indirect(v)

	if v.Type().String() == "float64" || v.Type().ConvertibleTo(floatType) {
		floatValue := v.Convert(floatType).Float()
		return floatValue, nil
	}
	return 0, errors.New("cond value can not convert to number")
}
