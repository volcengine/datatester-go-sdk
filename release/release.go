/**
 * Apache 2.0
 * Copyright 2022 Beijing Volcano Engine Technology Co., Ltd.
 */

package release

type Release struct {
	Filters           []Filter            `json:"filter"`
	TrafficAllocation []TrafficAllocation `json:"traffic_allocation"`
}

func (r *Release) EvaluateRelease(attributes map[string]interface{}, hashIndex uint32) string {
	if !EvaluateFilters(r.Filters, attributes) {
		return ""
	}
	return r.evaluateTraffics(hashIndex)
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
