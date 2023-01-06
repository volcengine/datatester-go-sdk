/**
 * Apache 2.0
 * Copyright 2022 Beijing Volcano Engine Technology Co., Ltd.
 */

package release

import "github.com/volcengine/datatester-go-sdk/release/cond"

type Filter struct {
	Id            string                 `json:"id"`
	Conditions    []cond.Condition       `json:"conditions"`
	LogicOperator cond.LogicOperatorType `json:"logic_operator"`
}

func (f *Filter) EvaluateFilter(attributes map[string]interface{}) bool {
	if f.Conditions == nil || len(f.Conditions) == 0 {
		return true
	}
	result := f.Conditions[0].Evaluate(attributes)
	for i := 0; i < len(f.Conditions)-1; i++ {
		logicFn := cond.LogicFunc(f.Conditions[i].LogicOperator)
		if logicFn == nil {
			return false
		}
		result = logicFn(result, f.Conditions[i+1].Evaluate(attributes))
	}
	return result
}
