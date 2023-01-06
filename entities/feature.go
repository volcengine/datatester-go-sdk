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

	whiteListMap map[string]Variant
}

func (f *Feature) GetWhiteListMap() map[string]Variant {
	if f.whiteListMap == nil {
		f.generateWhiteListMap()
	}
	return f.whiteListMap
}

func (f *Feature) generateWhiteListMap() {
	if f.WhiteList == nil || len(f.WhiteList) == 0 {
		return
	}
	f.whiteListMap = make(map[string]Variant, len(f.WhiteList))
	for uid, vid := range f.WhiteList {
		f.whiteListMap[uid] = f.VariantMap[vid]
	}
}
