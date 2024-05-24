/**
 * Apache 2.0
 * Copyright 2022 Beijing Volcano Engine Technology Co., Ltd.
 */

package release

type UserGroupRelease struct {
	UserGroupId string   `json:"user_group_id"`
	Filters     []Filter `json:"filter"`
}

func (r *UserGroupRelease) EvaluateRelease(attributes map[string]interface{}) bool {
	return EvaluateFilters(r.Filters, attributes)
}
