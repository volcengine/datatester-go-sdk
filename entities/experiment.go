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
	Id                  string                     `json:"id"`
	Name                string                     `json:"name"`
	LayerID             string                     `json:"layer_id"`
	Status              int                        `json:"status"`
	Release             release.Release            `json:"release"`
	WhiteList           map[string]string          `json:"white_list"`
	FreezeStatus        int                        `json:"freeze_status"`
	LaunchStartTime     int64                      `json:"launch_start_time"`
	VersionFreezeStatus int                        `json:"version_freeze_status"`
	SaveParam           int                        `json:"save_param"`
	ExperimentMode      int                        `json:"experiment_mode"`
	FatherExperimentId  string                     `json:"father_experiment_id"`
	VariantMap          map[string]*Variant        `json:"variants"`
	FilterAllowList     int                        `json:"filter_allowlist"`
	AssociatedRelations []string                   `json:"associated_relations"`
	ManageSubType       string                     `json:"manage_sub_type"`
	UserGroupReleases   []release.UserGroupRelease `json:"user_group_releases"`
	FlightPriority      int64                      `json:"experiment_priority"`
	whiteListMap        map[string]*Variant
	CohortIds           []string
}

func (e *Experiment) GetWhiteListMap() map[string]*Variant {
	if e.whiteListMap == nil {
		e.generateWhiteListMap()
	}
	return e.whiteListMap
}

func (e *Experiment) generateWhiteListMap() {
	if e.WhiteList == nil || len(e.WhiteList) == 0 {
		return
	}
	e.whiteListMap = make(map[string]*Variant, len(e.WhiteList))
	for uid, vid := range e.WhiteList {
		e.whiteListMap[uid] = e.VariantMap[vid]
	}
}

func (e *Experiment) GenerateCohortIds() []string {
	e.CohortIds = append(e.CohortIds, release.GetFilterCohortIds(e.Release.Filters)...)
	for _, r := range e.UserGroupReleases {
		e.CohortIds = append(e.CohortIds, release.GetFilterCohortIds(r.Filters)...)
	}
	return e.CohortIds
}

func (e *Experiment) IsCodingExperiment() bool {
	return e.ExperimentMode == consts.CodingExperiment
}

func (e *Experiment) IsUserGroupExperiment() bool {
	return e.ManageSubType == consts.UserGroupExperimentSubType
}

func (e *Experiment) IsCodingCampaign() bool {
	return e.ExperimentMode == consts.CampaignCodingMode || e.ExperimentMode == consts.CampaignCodingChildMode
}
