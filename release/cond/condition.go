/**
 * Apache 2.0
 * Copyright 2022 Beijing Volcano Engine Technology Co., Ltd.
 */

package cond

import (
	"github.com/volcengine/datatester-go-sdk/consts"
	"github.com/volcengine/datatester-go-sdk/log"
	"reflect"
	"strconv"
)

type LogicOperatorType = string

const (
	AND LogicOperatorType = "&&"
	OR  LogicOperatorType = "||"
)

type LogicFn func(x, y bool) bool

func LogicFunc(logic LogicOperatorType) LogicFn {
	switch logic {
	case AND:
		return func(x, y bool) bool {
			return x && y
		}
	case OR:
		return func(x, y bool) bool {
			return x || y
		}
	default:
		return nil
	}
}

type StringCompareMethod = string

const (
	DictCompare   StringCompareMethod = "dict"
	SemverCompare StringCompareMethod = "version"
)

type ConditionOpType string

const (
	GreaterThan        ConditionOpType = ">"
	GreaterThanOrEqual ConditionOpType = ">="
	LessThan           ConditionOpType = "<"
	LessThanOrEqual    ConditionOpType = "<="
	IN                 ConditionOpType = "in"
	NotIn              ConditionOpType = "ni"
	IsNull             ConditionOpType = "is_null"
	IsNotNull          ConditionOpType = "is_not_null"
)

type MatcherFn func(attrValue interface{}, condValue interface{}, method StringCompareMethod) bool

func MathFunc(op ConditionOpType) MatcherFn {
	switch op {
	case GreaterThan:
		return greaterThan
	case GreaterThanOrEqual:
		return greaterEqual
	case LessThan:
		return lessThan
	case LessThanOrEqual:
		return lessEqual
	case IN:
		return in
	case NotIn:
		return notIn
	case IsNull:
		return isNull
	case IsNotNull:
		return notNull
	default:
		return alwaysFalse
	}
}

type Condition struct {
	Key           string              `json:"key"`
	Op            ConditionOpType     `json:"op"`
	LogicOperator LogicOperatorType   `json:"logic_operator"`
	Value         interface{}         `json:"value"`
	Type          ConditionValueType  `json:"type"`
	Method        StringCompareMethod `json:"method"`
	PropertyType  string              `json:"property_type"`

	valueTyper ConditionValueTyper
}

func (c *Condition) Evaluate(attributes map[string]interface{}) bool {
	if len(c.Key) == 0 {
		return false
	}

	if c.Key == consts.ExperimentCohort && c.PropertyType == consts.ExperimentCohort {
		return c.JudgeExperimentCohort(attributes)
	}

	attrValue, ok := attributes[c.Key]

	if c.Op == IsNull || c.Op == IsNotNull {
		return MathFunc(c.Op)(attrValue, nil, c.Method)
	}

	if !ok || attrValue == nil {
		return false
	}

	if c.valueTyper == nil {
		c.generateValueTyper()
	}
	if !c.valueTyper.EvaluateKind(reflect.ValueOf(attrValue)) {
		return false
	}

	if c.Op == IN || c.Op == NotIn {
		condValue := reflect.ValueOf(c.Value)
		condKind := condValue.Kind()
		if condKind != reflect.Slice && condKind != reflect.Array {
			condValue = reflect.ValueOf([]interface{}{c.Value})
		}
		adaptedConditions := make([]interface{}, 0)
		var adaptedAttrs interface{}
		for i := 0; i < condValue.Len(); i++ {
			adaptedCond, adaptedAttr, err := c.valueTyper.AdaptValue(condValue.Index(i).Interface(), attrValue)
			if err != nil {
				log.WithFields(log.Fields{"err": err}).Error("cond value adapt occur an error")
				return false
			}
			adaptedConditions = append(adaptedConditions, adaptedCond)
			adaptedAttrs = adaptedAttr
		}
		return MathFunc(c.Op)(adaptedAttrs, adaptedConditions, c.Method)
	}
	return c.evaluateAttrValue(c.Value, attrValue)
}

func (c *Condition) evaluateAttrValue(condValue, attrValue interface{}) bool {
	adaptedCond, adaptedAttr, err := c.valueTyper.AdaptValue(condValue, attrValue)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("cond value adapt occur an error")
		return false
	}
	return MathFunc(c.Op)(adaptedAttr, adaptedCond, c.Method)
}

func (c *Condition) generateValueTyper() {
	switch c.Type {
	case INT, FLOAT, NUMBER:
		c.valueTyper = &NumberTyper{}
	case BOOL:
		c.valueTyper = &BoolTyper{}
	case STRING:
		c.valueTyper = &StringTyper{}
	}
}

func (c *Condition) JudgeExperimentCohort(attributes map[string]interface{}) bool {
	configValue := c.Value
	if configValue == nil {
		return true
	}
	result := true
	if c.Op == IN {
		result = false
	}
	experimentIds := configValue.([]interface{})
	for _, experimentId := range experimentIds {
		if attributes == nil {
			continue
		}
		experimentIdStr := strconv.FormatFloat(experimentId.(float64), 'f', 0, 64)
		value, ok := attributes[consts.ExperimentPrefix+experimentIdStr]
		if !ok {
			continue
		}
		if c.Op == IN {
			result = value.(bool) || result
		} else {
			result = (!value.(bool)) && result
		}
	}
	return result
}
