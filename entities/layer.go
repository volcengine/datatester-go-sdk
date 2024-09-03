/**
 * Apache 2.0
 * Copyright 2022 Beijing Volcano Engine Technology Co., Ltd.
 */

package entities

import (
	"github.com/volcengine/datatester-go-sdk/release"
)

type Layer struct {
	ID                string                      `json:"id"`
	Name              string                      `json:"name"`
	TrafficAllocation []release.TrafficAllocation `json:"traffic_allocation"`
	ExperimentIds     []string                    `json:"experiment_ids"`
	Domain            *Domain                     `json:"domain"`
	ParentDomains     []Domain                    `json:"parent_domains"`
}

// HitDomain 方法用于检查给定的 decisionID 是否命中当前层及其所有父域
func (l Layer) HitDomain(decisionID string) bool {
	// 层不属于互斥域，直接返回
	if l.Domain == nil || l.Domain.ID == "0" {
		return true
	}
	hit := l.Domain.Hit(decisionID)
	if !hit {
		return false
	}
	for _, d := range l.ParentDomains {
		hit = d.Hit(decisionID)
		if !hit {
			return false
		}
	}
	return true
}
