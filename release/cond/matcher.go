/**
 * Apache 2.0
 * Copyright 2022 Beijing Volcano Engine Technology Co., Ltd.
 */

package cond

import (
	"github.com/Masterminds/semver/v3"
	"reflect"
)

func greaterThan(attrValue interface{}, condValue interface{}, method StringCompareMethod) bool {
	switch attrValue.(type) {
	case float64:
		attrValueFloat64, condValueFloat64 := attrValue.(float64), condValue.(float64)
		return attrValueFloat64 > condValueFloat64
	case string:
		attrValueString, condValueString := attrValue.(string), condValue.(string)
		if method == SemverCompare {
			attrVersion, err := semver.NewVersion(attrValueString)
			if err != nil {
				return false
			}
			condVersion, err := semver.NewVersion(condValueString)
			if err != nil {
				return false
			}
			return attrVersion.GreaterThan(condVersion)
		} else {
			return attrValueString > condValueString
		}
	}
	return false
}

func greaterEqual(attrValue interface{}, condValue interface{}, method StringCompareMethod) bool {
	switch attrValue.(type) {
	case float64:
		attrValueFloat64, condValueFloat64 := attrValue.(float64), condValue.(float64)
		return attrValueFloat64 >= condValueFloat64
	case string:
		attrValueString, condValueString := attrValue.(string), condValue.(string)
		if method == SemverCompare {
			attrVersion, err := semver.NewVersion(attrValueString)
			if err != nil {
				return false
			}
			condVersion, err := semver.NewVersion(condValueString)
			if err != nil {
				return false
			}
			return attrVersion.GreaterThan(condVersion) || attrVersion.Equal(condVersion)
		} else {
			return attrValueString >= condValueString
		}
	}
	return false
}

func lessThan(attrValue interface{}, condValue interface{}, method StringCompareMethod) bool {
	switch attrValue.(type) {
	case float64:
		attrValueFloat64, condValueFloat64 := attrValue.(float64), condValue.(float64)
		return attrValueFloat64 < condValueFloat64
	case string:
		attrValueString, condValueString := attrValue.(string), condValue.(string)
		if method == SemverCompare {
			attrVersion, err := semver.NewVersion(attrValueString)
			if err != nil {
				return false
			}
			condVersion, err := semver.NewVersion(condValueString)
			if err != nil {
				return false
			}
			return attrVersion.LessThan(condVersion)
		} else {
			return attrValueString < condValueString
		}
	}
	return false
}

func lessEqual(attrValue interface{}, condValue interface{}, method StringCompareMethod) bool {
	switch attrValue.(type) {
	case float64:
		attrValueFloat64, condValueFloat64 := attrValue.(float64), condValue.(float64)
		return attrValueFloat64 <= condValueFloat64
	case string:
		attrValueString, condValueString := attrValue.(string), condValue.(string)
		if method == SemverCompare {
			attrVersion, err := semver.NewVersion(attrValueString)
			if err != nil {
				return false
			}
			condVersion, err := semver.NewVersion(condValueString)
			if err != nil {
				return false
			}
			return attrVersion.LessThan(condVersion) || attrVersion.Equal(condVersion)
		} else {
			return attrValueString <= condValueString
		}
	}
	return false
}

func equal(attrValue interface{}, condValue interface{}, method StringCompareMethod) bool {
	switch attrValue.(type) {
	case float64:
		attrValueFloat64, condValueFloat64 := attrValue.(float64), condValue.(float64)
		return attrValueFloat64 == condValueFloat64
	case string:
		attrValueString, condValueString := attrValue.(string), condValue.(string)
		if method == SemverCompare {
			attrVersion, err := semver.NewVersion(attrValueString)
			if err != nil {
				return false
			}
			condVersion, err := semver.NewVersion(condValueString)
			if err != nil {
				return false
			}
			return attrVersion.Equal(condVersion)
		} else {
			return attrValueString == condValueString
		}
	case bool:
		attrValueBool, condValueBool := attrValue.(bool), condValue.(bool)
		return attrValueBool == condValueBool
	}
	return false
}

func notEqual(attrValue interface{}, condValue interface{}, method StringCompareMethod) bool {
	switch attrValue.(type) {
	case float64:
		attrValueFloat64, condValueFloat64 := attrValue.(float64), condValue.(float64)
		return attrValueFloat64 != condValueFloat64
	case string:
		attrValueString, condValueString := attrValue.(string), condValue.(string)
		if method == SemverCompare {
			attrVersion, err := semver.NewVersion(attrValueString)
			if err != nil {
				return false
			}
			condVersion, err := semver.NewVersion(condValueString)
			if err != nil {
				return false
			}
			return !attrVersion.Equal(condVersion)
		} else {
			return attrValueString != condValueString
		}
	case bool:
		attrValueBool, condValueBool := attrValue.(bool), condValue.(bool)
		return attrValueBool != condValueBool
	}
	return false
}

func in(attrValue interface{}, condValue interface{}, method StringCompareMethod) bool {
	condValues := reflect.ValueOf(condValue)
	for i := 0; i < condValues.Len(); i++ {
		if equal(attrValue, condValues.Index(i).Interface(), method) {
			return true
		}
	}
	return false
}

func notIn(attrValue interface{}, condValue interface{}, method StringCompareMethod) bool {
	condValues := reflect.ValueOf(condValue)
	for i := 0; i < condValues.Len(); i++ {
		if !notEqual(attrValue, condValues.Index(i).Interface(), method) {
			return false
		}
	}
	return true
}

func isNull(attrValue interface{}, condValue interface{}, method StringCompareMethod) bool {
	if attrValue == nil {
		return true
	}
	return false
}

func notNull(attrValue interface{}, condValue interface{}, method StringCompareMethod) bool {
	if attrValue != nil {
		return true
	}
	return false
}

func alwaysFalse(attrValue interface{}, condValue interface{}, method StringCompareMethod) bool {
	return false
}
