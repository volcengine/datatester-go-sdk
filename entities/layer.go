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
}
