/**
 * Apache 2.0
 * Copyright 2022 Beijing Volcano Engine Technology Co., Ltd.
 */

package entities

type Variant struct {
	Name             string              `json:"name"`
	Id               string              `json:"id"`
	ExperimentId     string              `json:"entity_id"`
	Config           map[string]Variable `json:"config"`
	Type             uint32              `json:"type"`
	FatherVariantIds []string            `json:"father_variants"`
}

func (v *Variant) GetConfig() map[string]map[string]interface{} {
	if v.Config == nil || len(v.Config) == 0 {
		return nil
	}
	config := make(map[string]map[string]interface{})
	for key, value := range v.Config {
		config[key] = map[string]interface{}{"val": value.Value, "vid": v.Id}
	}
	return config
}
