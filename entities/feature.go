/**
 * Apache 2.0
 * Copyright 2022 Beijing Volcano Engine Technology Co., Ltd.
 */

package entities

import (
	"github.com/volcengine/datatester-go-sdk/release"
)

type Feature struct {
	Id              string             `json:"id"`
	Name            string             `json:"name"`
	Key             string             `json:"variant_key"`
	Releases        []release.Release  `json:"releases"`
	Status          int                `json:"status"`
	WhiteList       map[string]string  `json:"white_list"`
	LaunchStartTime int64              `json:"launch_start_time"`
	VariantMap      map[string]Variant `json:"variants"`

	allowListMap map[string]Variant
	CohortIds    []string
}

func (f *Feature) GetAllowListMap() map[string]Variant {
	if f.allowListMap == nil {
		f.generateAllowListMap()
	}
	return f.allowListMap
}

func (f *Feature) generateAllowListMap() {
	if f.WhiteList == nil || len(f.WhiteList) == 0 {
		return
	}
	f.allowListMap = make(map[string]Variant, len(f.WhiteList))
	for uid, vid := range f.WhiteList {
		f.allowListMap[uid] = f.VariantMap[vid]
	}
}

func (f *Feature) GenerateCohortIds() []string {
	for _, r := range f.Releases {
		f.CohortIds = append(f.CohortIds, release.GetFilterCohortIds(r.Filters)...)
	}
	return f.CohortIds
}
