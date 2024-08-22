/**
 * Apache 2.0
 * Copyright 2022 Beijing Volcano Engine Technology Co., Ltd.
 */

package release

import (
	"github.com/volcengine/datatester-go-sdk/release/cond"
	"reflect"
	"strconv"
)

// GetFilterCohortIds extract all cohort ids contained in filter
func GetFilterCohortIds(filters []Filter) []string {
	var cohortIds []string
	for _, outerCondition := range filters {
		for _, innerCondition := range outerCondition.Conditions {
			if innerCondition.Key != cond.CohortKeyName {
				continue
			}
			values, ok := innerCondition.Value.([]interface{})
			if !ok {
				continue
			}
			for _, value := range values {
				if v, ok := value.(float64); ok {
					cohortIds = append(cohortIds, strconv.FormatFloat(v, 'f', 0, 64))
				}
			}
		}
	}
	return cohortIds
}

// EvaluateFilters evaluate
func EvaluateFilters(filters []Filter, attributes map[string]interface{}) bool {
	if filters == nil || len(filters) == 0 {
		return true
	}
	if reflect.ValueOf(attributes).Kind() != reflect.Map {
		return false
	}
	result := filters[0].EvaluateFilter(attributes)
	for i := 0; i < len(filters)-1; i++ {
		logicFn := cond.LogicFunc(filters[i].LogicOperator)
		if logicFn == nil {
			return false
		}
		result = logicFn(result, filters[i+1].EvaluateFilter(attributes))
	}
	return result
}
