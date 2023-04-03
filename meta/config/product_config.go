/**
 * Apache 2.0
 * Copyright 2022 Beijing Volcano Engine Technology Co., Ltd.
 */

package config

import (
	"errors"
	"fmt"
	et "github.com/volcengine/datatester-go-sdk/entities"
)

type ProductConfig struct {
	AppId         string                   `json:"app_id"`
	HashStrategy  string                   `json:"hash_strategy"`
	ModifyTime    int                      `json:"modify_time"`
	LayerMap      map[string]et.Layer      `json:"layers"`
	ExperimentMap map[string]et.Experiment `json:"experiments"`
	FeatureMap    map[string]et.Feature    `json:"features"`

	// key:name, value:entity id
	ExperimentNameToIdMap map[string]string
	FeatureNameToIdMap    map[string]string
}

// Init entity map and white list map for cache
func (c *ProductConfig) Init() {
	c.generateExperimentNameToIdMap()
	c.generateFeatureNameToIdMap()
}

func (c *ProductConfig) generateExperimentNameToIdMap() {
	c.ExperimentNameToIdMap = make(map[string]string, len(c.ExperimentMap))
	for id, experiment := range c.ExperimentMap {
		c.ExperimentNameToIdMap[experiment.Name] = id
	}
}

func (c *ProductConfig) generateFeatureNameToIdMap() {
	c.FeatureNameToIdMap = make(map[string]string, len(c.FeatureMap))
	for id, feature := range c.FeatureMap {
		c.FeatureNameToIdMap[feature.Name] = id
	}
}

func (c *ProductConfig) GetExperimentFromId(experimentId string) (et.Experiment, error) {
	if c.ExperimentMap == nil {
		return et.Experiment{}, errors.New("no experiments in product config")
	}
	if experiment, ok := c.ExperimentMap[experimentId]; ok {
		return experiment, nil
	}
	return et.Experiment{}, fmt.Errorf("failed to find experiment[%s]", experimentId)
}

func (c *ProductConfig) GetFeatureFromId(featureId string) (et.Feature, error) {
	if c.FeatureMap == nil {
		return et.Feature{}, errors.New("no features in product config")
	}
	if feature, ok := c.FeatureMap[featureId]; ok {
		return feature, nil
	}
	return et.Feature{}, fmt.Errorf("failed to find feature[%s]", featureId)
}

func (c *ProductConfig) GetFeatureAllowList(featureId string) (map[string]et.Variant, error) {
	feature, err := c.GetFeatureFromId(featureId)
	if err != nil {
		return nil, err
	}
	return feature.GetAllowListMap(), nil
}
