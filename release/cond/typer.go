/**
 * Apache 2.0
 * Copyright 2022 Beijing Volcano Engine Technology Co., Ltd.
 */

package cond

import (
	"reflect"
	"strconv"
)

type ConditionValueType string

const (
	INT    ConditionValueType = "int"
	FLOAT  ConditionValueType = "float"
	NUMBER ConditionValueType = "number"
	STRING ConditionValueType = "string"
	BOOL   ConditionValueType = "boolean"
)

var (
	numberValueTyper = NumberTyper{}
	boolValueTyper   = BoolTyper{}
	stringValueTyper = StringTyper{}
)

type ConditionValueTyper interface {
	EvaluateKind(v reflect.Value) bool
	AdaptValue(condValue, attrValue interface{}) (adaptedCond interface{}, adaptedAttr interface{}, err error)
}

type NumberTyper struct {
}

func (n NumberTyper) EvaluateKind(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Invalid:
		return false
	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		return true
	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return true
	case reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}

func (n NumberTyper) AdaptValue(condValue, attrValue interface{}) (float64Cond, float64Attr interface{}, err error) {
	if condValueStr, ok := condValue.(string); ok {
		float64Cond, err = strconv.ParseFloat(condValueStr, 64)
	} else {
		float64Cond, err = ToFloat64(condValue)
	}
	float64Attr, err = ToFloat64(attrValue)
	return float64Cond, float64Attr, err
}

type BoolTyper struct {
}

func (b BoolTyper) EvaluateKind(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Invalid:
		return false
	case reflect.Bool:
		return true
	default:
		return false
	}
}

func (b BoolTyper) AdaptValue(condValue, attrValue interface{}) (boolCond interface{}, boolAttr interface{}, err error) {
	if condValueStr, ok := condValue.(string); ok {
		boolCond, err = strconv.ParseBool(condValueStr)
	} else {
		boolCond = condValue.(bool)
	}
	boolAttr = attrValue.(bool)
	return boolCond, boolAttr, err
}

type StringTyper struct {
}

func (s StringTyper) EvaluateKind(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Invalid:
		return false
	case reflect.String:
		return true
	default:
		return false
	}
}

func (s StringTyper) AdaptValue(condValue, attrValue interface{}) (stringCond interface{}, stringAttr interface{}, err error) {
	return condValue.(string), attrValue.(string), nil
}
