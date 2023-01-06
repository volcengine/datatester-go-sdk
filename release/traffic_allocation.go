/**
 * Apache 2.0
 * Copyright 2022 Beijing Volcano Engine Technology Co., Ltd.
 */

package release

type TrafficAllocation struct {
	Begin    uint32 `json:"begin"`
	End      uint32 `json:"end"`
	EntityId string `json:"entity_id"`
}

func (t *TrafficAllocation) EvaluateTraffic(hashIndex uint32) string {
	if hashIndex < t.Begin || hashIndex >= t.End {
		return ""
	}
	return t.EntityId
}
