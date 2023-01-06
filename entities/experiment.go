/**
 * Apache 2.0
 * Copyright 2022 Beijing Volcano Engine Technology Co., Ltd.
 */

package entities

import (
	"github.com/volcengine/datatester-go-sdk/consts"
	"github.com/volcengine/datatester-go-sdk/release"
)

type Experiment struct {
	Id                  string             `json:"id"`
	Name                string             `json:"name"`
	LayerID             string             `json:"layer_id"`
	Status              int                `json:"status"`
	Release             release.Release    `json:"release"`
	WhiteList           map[string]string  `json:"white_list"`
	FreezeStatus        int                `json:"freeze_status"`
	LaunchStartTime     int64              `json:"launch_start_time"`
	VersionFreezeStatus int                `json:"version_freeze_status"`
	SaveParam           int                `json:"save_param"`
	ExperimentMode      int                `json:"experiment_mode"`
	FatherExperimentId  string             `json:"father_experiment_id"`
	VariantMap          map[string]Variant `json:"variants"`
	FilterWhitelist     int                `json:"filter_allowlist"`
	AssociatedRelations []string           `json:"associated_relations"`
	whiteListMap        map[string]Variant
}

func (e *Experiment) GetWhiteListMap() map[string]Variant {
	if e.whiteListMap == nil {
		e.generateWhiteListMap()
	}
	return e.whiteListMap
}

func (e *Experiment) generateWhiteListMap() {
	if e.WhiteList == nil || len(e.WhiteList) == 0 {
		return
	}
	e.whiteListMap = make(map[string]Variant, len(e.WhiteList))
	for uid, vid := range e.WhiteList {
		e.whiteListMap[uid] = e.VariantMap[vid]
	}
}

func (e *Experiment) IsCodingExperiment() bool {
	return e.ExperimentMode == consts.CodingExperiment
}
