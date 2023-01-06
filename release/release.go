/**
 * Apache 2.0
 * Copyright 2022 Beijing Volcano Engine Technology Co., Ltd.
 */

package release

import (
	"github.com/volcengine/datatester-go-sdk/release/cond"
	"reflect"
)

type Release struct {
	Filters           []Filter            `json:"filter"`
	TrafficAllocation []TrafficAllocation `json:"traffic_allocation"`
}

func (r *Release) EvaluateRelease(attributes map[string]interface{}, hashIndex uint32) string {
	if !r.EvaluateFilters(attributes) {
		return ""
	}
	return r.evaluateTraffics(hashIndex)
}

func (r *Release) EvaluateFilters(attributes map[string]interface{}) bool {
	if r.Filters == nil || len(r.Filters) == 0 {
		return true
	}
	if reflect.ValueOf(attributes).Kind() != reflect.Map {
		return false
	}
	result := r.Filters[0].EvaluateFilter(attributes)
	for i := 0; i < len(r.Filters)-1; i++ {
		logicFn := cond.LogicFunc(r.Filters[i].LogicOperator)
		if logicFn == nil {
			return false
		}
		result = logicFn(result, r.Filters[i+1].EvaluateFilter(attributes))
	}
	return result
}

func (r *Release) evaluateTraffics(hashIndex uint32) string {
	if r.TrafficAllocation == nil || len(r.TrafficAllocation) == 0 {
		return ""
	}
	for _, t := range r.TrafficAllocation {
		if entityId := t.EvaluateTraffic(hashIndex); len(entityId) != 0 {
			return entityId
		}
	}
	return ""
}
