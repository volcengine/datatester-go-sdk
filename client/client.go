/**
 * Apache 2.0
 * Copyright 2022 Beijing Volcano Engine Technology Co., Ltd.
 */

package client

import (
	"encoding/json"
	"fmt"
	"github.com/volcengine/datatester-go-sdk/cohort"
	"github.com/volcengine/datatester-go-sdk/config"
	"github.com/volcengine/datatester-go-sdk/distributor"
	"github.com/volcengine/datatester-go-sdk/entities"
	"github.com/volcengine/datatester-go-sdk/event"
	"github.com/volcengine/datatester-go-sdk/handler"
	"github.com/volcengine/datatester-go-sdk/log"
	"github.com/volcengine/datatester-go-sdk/meta/manager"
	"github.com/volcengine/datatester-go-sdk/utils"
)

type AbClient struct {
	token             string
	metaManager       *manager.DynamicMetaManager
	distributor       *distributor.VariantsDistributor
	dispatcher        event.Dispatcher
	userAbInfoHandler handler.UserAbInfoHandler
	cohortHandler     cohort.Client
}

func NewClient(token string, configs ...config.Func) *AbClient {
	return NewClientWithUserAbInfo(token, handler.NewDefaultUserAbInfoHandler(), configs...)
}

func NewClientWithUserAbInfo(token string, userAbInfoHandler handler.UserAbInfoHandler,
	configs ...config.Func) *AbClient {
	metaOptionFunc := initConfigAndTransferMetaOptionFunc(configs...)
	if len(token) == 0 {
		log.ErrorF("token[%s] invalid", token)
		return nil
	}
	m := manager.NewDynamicMetaManager(token, metaOptionFunc...)
	if m.GetConfig() == nil {
		log.ErrorF("meta url [%s] invalid", m.GetFetchUrl())
		return nil
	}
	appId := m.GetConfig().AppId
	if len(appId) == 0 {
		log.ErrorF("token[%s] invalid", token)
		return nil
	}

	d := event.NewDispatcher(appId)
	return &AbClient{
		token:             token,
		metaManager:       m,
		distributor:       distributor.NewVariantsDistributor(),
		dispatcher:        d,
		userAbInfoHandler: userAbInfoHandler,
	}
}

func (t *AbClient) SetCohortHandler(handler cohort.Client) {
	t.cohortHandler = handler
}

func initConfigAndTransferMetaOptionFunc(configs ...config.Func) []manager.MetaOptionFunc {
	metaOptionFuncArray := make([]manager.MetaOptionFunc, 0)
	log.InitDefaultLogger()
	for _, configFunc := range configs {
		metaOptionFunc, shouldSetMetaManager := configFunc()
		if shouldSetMetaManager {
			metaOptionFuncArray = append(metaOptionFuncArray, metaOptionFunc)
		}
	}
	return metaOptionFuncArray
}

// NewClient4Test for tests
func NewClient4Test(meta string, userAbInfoHandler handler.UserAbInfoHandler) *AbClient {
	log.InitDefaultLogger()
	m := manager.NewDynamicMetaManager4Test(meta)
	d := event.NewDispatcher4Test()
	return &AbClient{
		metaManager:       m,
		distributor:       distributor.NewVariantsDistributor(),
		dispatcher:        d,
		userAbInfoHandler: userAbInfoHandler,
	}
}

func (t *AbClient) GetExperimentVariantName(experimentId, decisionId string,
	attributes map[string]interface{}) (string, error) {
	if attributes == nil {
		attributes = make(map[string]interface{})
	}
	experiment2variant := t.initUserAbInfo(decisionId)
	t.initExpCohort(experimentId, decisionId, attributes)
	variant, err := t.distributor.GetExperimentVariant(t.metaManager.GetConfig(), experimentId,
		decisionId, attributes, experiment2variant)
	if err != nil || variant == nil {
		return "", err
	}
	t.updateUserAbInfo(decisionId, experiment2variant)
	return variant.Name, nil
}

func (t *AbClient) GetExperimentVariantNameWithImpression(experimentId, decisionId, trackId string,
	attributes map[string]interface{}) (string, error) {
	variant, err := t.GetExperimentVariantWithImpression(experimentId, decisionId, trackId, attributes)
	if err != nil || len(variant.Id) == 0 {
		return "", err
	}
	return variant.Name, nil
}

func (t *AbClient) GetExperimentVariantWithImpression(experimentId, decisionId, trackId string,
	attributes map[string]interface{}) (entities.Variant, error) {
	if attributes == nil {
		attributes = make(map[string]interface{})
	}
	experiment2variant := t.initUserAbInfo(decisionId)
	t.initExpCohort(experimentId, decisionId, attributes)
	variant, err := t.distributor.GetExperimentVariant(t.metaManager.GetConfig(), experimentId,
		decisionId, attributes, experiment2variant)
	if err != nil || variant == nil {
		return entities.Variant{}, err
	}
	if err = t.dispatcher.DispatchEvent(trackId, variant.Id, attributes); err != nil {
		return entities.Variant{}, err
	}
	t.updateUserAbInfo(decisionId, experiment2variant)
	return *variant, nil
}

func (t *AbClient) GetExperimentConfigs(experimentId, decisionId string,
	attributes map[string]interface{}) (map[string]map[string]interface{}, error) {
	if attributes == nil {
		attributes = make(map[string]interface{})
	}
	experiment2variant := t.initUserAbInfo(decisionId)
	t.initExpCohort(experimentId, decisionId, attributes)
	variant, err := t.distributor.GetExperimentVariant(t.metaManager.GetConfig(), experimentId, decisionId,
		attributes, experiment2variant)
	if err != nil || variant == nil {
		return nil, err
	}
	t.updateUserAbInfo(decisionId, experiment2variant)
	return variant.GetConfig(), nil
}

func (t *AbClient) GetExperimentConfigsWithImpression(experimentId, decisionId, trackId string,
	attributes map[string]interface{}) (map[string]map[string]interface{}, error) {
	variant, err := t.GetExperimentVariantWithImpression(experimentId, decisionId, trackId, attributes)
	if err != nil || len(variant.Id) == 0 {
		return nil, err
	}
	return variant.GetConfig(), nil
}

func (t *AbClient) GetAllExperimentConfigs(decisionId string,
	attributes map[string]interface{}) (map[string]map[string]interface{}, error) {
	if attributes == nil {
		attributes = make(map[string]interface{})
	}
	experiment2variant := t.initUserAbInfo(decisionId)
	t.initAllExpCohort(decisionId, attributes)
	variants := make([]*entities.Variant, 0)
	for _, experiment := range t.metaManager.GetConfig().ExperimentMap {
		variant, err := t.distributor.GetExperimentVariant(t.metaManager.GetConfig(), experiment.Id, decisionId,
			attributes, experiment2variant)
		if err != nil || variant == nil {
			continue
		}
		variants = append(variants, variant)
	}
	configs := make(map[string]map[string]interface{})
	vid2ExperimentIdMap := make(map[string]string)
	for _, variant := range variants {
		confMap := variant.GetConfig()
		if confMap == nil {
			continue
		}
		for k, v := range confMap {
			existConfig, conflict := configs[k]
			if !conflict || utils.IsHigherPriorityConfig2(existConfig, v,
				vid2ExperimentIdMap[existConfig["vid"].(string)], variant.ExperimentId, t.metaManager.GetConfig().ExperimentMap) {
				configs[k] = v
			}
		}
		vid2ExperimentIdMap[variant.Id] = variant.ExperimentId
	}
	t.updateUserAbInfo(decisionId, experiment2variant)
	return configs, nil
}

func (t *AbClient) getAllExperimentConfigs4Activate(variantKey, decisionId string,
	attributes map[string]interface{}) (map[string]map[string]interface{}, error) {
	experiment2variant := t.initUserAbInfo(decisionId)
	t.initAllExpCohort(decisionId, attributes)
	experiment2variantCopy := make(map[string]string)
	for k, v := range experiment2variant {
		experiment2variantCopy[k] = v
	}
	variants := make([]*entities.Variant, 0)
	for _, experiment := range t.metaManager.GetConfig().ExperimentMap {
		variant, err := t.distributor.GetExperimentVariant(t.metaManager.GetConfig(), experiment.Id, decisionId,
			attributes, experiment2variantCopy)
		if err != nil || variant == nil {
			continue
		}
		variants = append(variants, variant)
	}
	configs := make(map[string]map[string]interface{})
	vid2ExperimentIdMap := make(map[string]string)
	for _, variant := range variants {
		confMap := variant.GetConfig()
		if confMap == nil {
			continue
		}
		for k, v := range confMap {
			existConfig, conflict := configs[k]
			if !conflict || utils.IsHigherPriorityConfig2(existConfig, v,
				vid2ExperimentIdMap[existConfig["vid"].(string)], variant.ExperimentId, t.metaManager.GetConfig().ExperimentMap) {
				configs[k] = v
			}
		}
		vid2ExperimentIdMap[variant.Id] = variant.ExperimentId
	}
	if valueMap, ok := configs[variantKey]; ok {
		vid := valueMap["vid"].(string)
		experimentId, exist := vid2ExperimentIdMap[vid]
		if exist && len(experimentId) > 0 {
			experiment2variant[experimentId] = vid
		}
	}
	t.updateUserAbInfo(decisionId, experiment2variant)
	return configs, nil
}

func (t *AbClient) GetAllExperimentConfigsWithImpression(decisionId string, trackId string,
	attributes map[string]interface{}) (map[string]map[string]interface{}, error) {
	if attributes == nil {
		attributes = make(map[string]interface{})
	}
	allExperimentConfigs, err := t.GetAllExperimentConfigs(decisionId, attributes)
	if err == nil && len(allExperimentConfigs) != 0 {
		for _, valueMap := range allExperimentConfigs {
			if err = t.dispatcher.DispatchEvent(trackId, valueMap["vid"].(string), attributes); err != nil {
				return nil, err
			}
		}
		return allExperimentConfigs, nil
	}
	return nil, fmt.Errorf("no hit experiment or feature")
}

// GetExperimentConfigsWithImpressionWithVariantKeys 获取指定变体list的实验配置，并返回指定变体键的值
// 参数:
//   - decisionId: 分流ID
//   - trackId: 上报ID
//   - variantKeys: 变体键列表
//   - attributes: 属性映射
//
// 返回值:
//   - 变体对应的值映射，如果变体没有命中任何实验，则变体的值为nil
//   - 错误信息
func (t *AbClient) GetExperimentConfigsWithImpressionWithVariantKeys(decisionId string, trackId string,
	variantKeys []string, attributes map[string]interface{}) (map[string]interface{}, error) {
	if attributes == nil {
		attributes = make(map[string]interface{})
	}
	if len(variantKeys) == 0 {
		return nil, fmt.Errorf("variantKeys is empty")
	}
	allExperimentConfigs, err1 := t.GetAllExperimentConfigs(decisionId, attributes)
	if err1 != nil {
		return nil, err1
	}
	variantsResult := make(map[string]interface{})
	for _, variantKey := range variantKeys {
		if valueMap, ok := allExperimentConfigs[variantKey]; ok {
			variantsResult[variantKey] = valueMap["val"]
			_ = t.dispatcher.DispatchEvent(trackId, valueMap["vid"].(string), attributes)
			continue
		}
		variantsResult[variantKey] = nil
	}
	return variantsResult, nil
}

func (t *AbClient) Activate(variantKey, decisionId, trackId string, defaultValue interface{},
	attributes map[string]interface{}) (interface{}, error) {
	if attributes == nil {
		attributes = make(map[string]interface{})
	}
	experimentConfigs, err := t.getAllExperimentConfigs4Activate(variantKey, decisionId, attributes)
	if err == nil && len(experimentConfigs) != 0 {
		if value, err := t.activateConfig(variantKey, experimentConfigs, trackId, attributes); err == nil {
			return value, nil
		}
	}
	featureConfigs, err := t.GetAllFeatureConfigs(decisionId, attributes)
	if err == nil || len(featureConfigs) != 0 {
		if value, err := t.activateConfig(variantKey, featureConfigs, trackId, attributes); err == nil {
			return value, nil
		}
	}
	return defaultValue, fmt.Errorf("no hit experiment or feature")
}

func (t *AbClient) activateConfig(variantKey string, Configs map[string]map[string]interface{},
	trackId string, attributes map[string]interface{}) (interface{}, error) {
	if valueMap, ok := Configs[variantKey]; ok {
		if err := t.dispatcher.DispatchEvent(trackId, valueMap["vid"].(string), attributes); err != nil {
			return nil, err
		}
		return valueMap["val"], nil
	}
	return nil, fmt.Errorf("no value exist in config[%v]", variantKey)
}

func (t *AbClient) ActivateWithIdType(variantKey, decisionId, trackId string, defaultValue interface{},
	attributes map[string]interface{}, uuidType string) (interface{}, error) {
	if attributes == nil {
		attributes = make(map[string]interface{})
	}
	experimentConfigs, err := t.getAllExperimentConfigs4Activate(variantKey, decisionId, attributes)
	if err == nil && len(experimentConfigs) != 0 {
		if value, err := t.activateConfigWithIdType(variantKey, experimentConfigs, trackId, uuidType,
			attributes); err == nil {
			return value, nil
		}
	}
	featureConfigs, err := t.GetAllFeatureConfigs(decisionId, attributes)
	if err == nil || len(featureConfigs) != 0 {
		if value, err := t.activateConfigWithIdType(variantKey, featureConfigs, trackId, uuidType,
			attributes); err == nil {
			return value, nil
		}
	}
	return defaultValue, fmt.Errorf("no hit experiment or feature")
}

func (t *AbClient) activateConfigWithIdType(variantKey string, Configs map[string]map[string]interface{},
	trackId string, uuidType string, attributes map[string]interface{}) (interface{}, error) {
	if valueMap, ok := Configs[variantKey]; ok {
		if err := t.dispatcher.DispatchEventWithIdType(trackId, valueMap["vid"].(string), uuidType,
			attributes); err != nil {
			return nil, err
		}
		return valueMap["val"], nil
	}
	return nil, fmt.Errorf("no value exist in config[%v]", variantKey)
}

func (t *AbClient) ActivateWithVid(variantKey, decisionId, trackId string,
	attributes map[string]interface{}) (map[string]interface{}, error) {
	if attributes == nil {
		attributes = make(map[string]interface{})
	}
	experimentConfigs, err := t.getAllExperimentConfigs4Activate(variantKey, decisionId, attributes)
	if err == nil && len(experimentConfigs) != 0 {
		if value, err := t.activateConfigWithVid(variantKey, experimentConfigs, trackId, attributes); err == nil {
			return value, nil
		}
	}
	featureConfigs, err := t.GetAllFeatureConfigs(decisionId, attributes)
	if err == nil || len(featureConfigs) != 0 {
		if value, err := t.activateConfigWithVid(variantKey, featureConfigs, trackId, attributes); err == nil {
			return value, nil
		}
	}
	return nil, fmt.Errorf("no hit experiment or feature")
}

func (t *AbClient) activateConfigWithVid(variantKey string, Configs map[string]map[string]interface{},
	trackId string, attributes map[string]interface{}) (map[string]interface{}, error) {
	if valueMap, ok := Configs[variantKey]; ok {
		if err := t.dispatcher.DispatchEvent(trackId, valueMap["vid"].(string), attributes); err != nil {
			return nil, err
		}
		return valueMap, nil
	}
	return nil, fmt.Errorf("no value exist in config[%v]", variantKey)
}

func (t *AbClient) ActivateWithoutImpression(variantKey, decisionId string,
	attributes map[string]interface{}) (map[string]interface{}, error) {
	if attributes == nil {
		attributes = make(map[string]interface{})
	}
	experimentConfigs, err := t.getAllExperimentConfigs4Activate(variantKey, decisionId, attributes)
	if err == nil && len(experimentConfigs) != 0 {
		if value, err := t.activateConfigWithoutImpression(variantKey, experimentConfigs); err == nil {
			return value, nil
		}
	}
	featureConfigs, err := t.GetAllFeatureConfigs(decisionId, attributes)
	if err == nil || len(featureConfigs) != 0 {
		if value, err := t.activateConfigWithoutImpression(variantKey, featureConfigs); err == nil {
			return value, nil
		}
	}
	return nil, fmt.Errorf("no hit experiment or feature")
}

func (t *AbClient) activateConfigWithoutImpression(variantKey string,
	Configs map[string]map[string]interface{}) (map[string]interface{}, error) {
	if valueMap, ok := Configs[variantKey]; ok {
		return valueMap, nil
	}
	return nil, fmt.Errorf("no value exist in config[%v]", variantKey)
}

func (t *AbClient) VerifyFeatureEnabled(featureId, decisionId string, attributes map[string]interface{}) (bool, error) {
	if attributes == nil {
		attributes = make(map[string]interface{})
	}
	t.initFeatureCohort(featureId, decisionId, attributes)
	variant, err := t.distributor.GetFeatureVariant(t.metaManager.GetConfig(), featureId, decisionId, attributes)
	if err != nil || variant == nil {
		return false, err
	}
	return true, nil
}

func (t *AbClient) GetFeatureConfigs(featureId, decisionId string,
	attributes map[string]interface{}) (map[string]map[string]interface{}, error) {
	if attributes == nil {
		attributes = make(map[string]interface{})
	}
	t.initFeatureCohort(featureId, decisionId, attributes)
	variant, err := t.distributor.GetFeatureVariant(t.metaManager.GetConfig(), featureId, decisionId, attributes)
	if err != nil || variant == nil {
		return nil, err
	}
	return variant.GetConfig(), nil
}

func (t *AbClient) GetFeatureConfigsWithImpression(featureId, decisionId, trackId string,
	attributes map[string]interface{}) (map[string]map[string]interface{}, error) {
	if attributes == nil {
		attributes = make(map[string]interface{})
	}
	t.initFeatureCohort(featureId, decisionId, attributes)
	variant, err := t.distributor.GetFeatureVariant(t.metaManager.GetConfig(), featureId, decisionId, attributes)
	if err != nil || variant == nil {
		return nil, err
	}
	if err = t.dispatcher.DispatchEvent(trackId, variant.Id, attributes); err != nil {
		return nil, err
	}
	return variant.GetConfig(), nil
}

func (t *AbClient) GetAllFeatureConfigs(decisionId string,
	attributes map[string]interface{}) (map[string]map[string]interface{}, error) {
	if attributes == nil {
		attributes = make(map[string]interface{})
	}
	t.initAllFeatureCohort(decisionId, attributes)
	variants, _, err := t.distributor.GetAllFeatureVariants(t.metaManager.GetConfig(), decisionId, attributes)
	if err != nil || len(variants) == 0 {
		return nil, err
	}
	configs := make(map[string]map[string]interface{})
	for _, variant := range variants {
		confMap := variant.GetConfig()
		if confMap == nil {
			continue
		}
		for k, v := range confMap {
			configs[k] = v
		}
	}
	return configs, nil
}

func (t *AbClient) GetEnabledFeatureIds(decisionId string, attributes map[string]interface{}) ([]string, error) {
	if attributes == nil {
		attributes = make(map[string]interface{})
	}
	t.initAllFeatureCohort(decisionId, attributes)
	_, featureIds, err := t.distributor.GetAllFeatureVariants(t.metaManager.GetConfig(), decisionId, attributes)
	if err != nil || len(featureIds) == 0 {
		return nil, err
	}
	return featureIds, nil
}

func (t *AbClient) initUserAbInfo(decisionId string) map[string]string {
	ex2variantStr := t.userAbInfoHandler.Query(decisionId)
	result := make(map[string]string)
	if ex2variantStr != "" {
		_ = json.Unmarshal([]byte(ex2variantStr), &result)
	}
	return result
}

func (t *AbClient) updateUserAbInfo(decisionId string, experiment2variant map[string]string) bool {
	if !t.userAbInfoHandler.NeedPersistData() {
		return true
	}
	result := make(map[string]string)
	for exId, vId := range experiment2variant {
		if _, ok := t.metaManager.GetConfig().ExperimentMap[exId]; ok {
			result[exId] = vId
		}
	}
	ex2variantByte, _ := json.Marshal(result)
	return t.userAbInfoHandler.CreateOrUpdate(decisionId, string(ex2variantByte))
}

func initAttributeCohort(cohortIds []string, attributes map[string]interface{}) {
	cohortIdMap := make(map[string]bool)
	for _, id := range cohortIds {
		cohortIdMap[id] = true
	}
	attributes["cohort"] = cohortIdMap
}

func (t *AbClient) initExpCohort(experimentId string, decisionId string, attributes map[string]interface{}) {
	if t.cohortHandler == nil {
		return
	}

	// collect cohort ids of the specified experiment as well as its father
	var cohortIds []string
	targetExperimentId := experimentId
	for len(targetExperimentId) != 0 {
		experiment, err := t.metaManager.GetConfig().GetExperimentFromId(targetExperimentId)
		if err != nil {
			break
		}
		cohortIds = append(cohortIds, experiment.CohortIds...)
		targetExperimentId = experiment.FatherExperimentId
	}
	if len(cohortIds) == 0 {
		return
	}

	// call cohort handler, and init cohort into attributes
	hitCohortIds := t.cohortHandler.UseCohort(decisionId, cohortIds)
	initAttributeCohort(hitCohortIds, attributes)
}

func (t *AbClient) initAllExpCohort(decisionId string, attributes map[string]interface{}) {
	if t.cohortHandler == nil || len(t.metaManager.GetConfig().ExpCohortIds) == 0 {
		return
	}

	// call cdp, and init cohort into attributes
	hitCohortIds := t.cohortHandler.UseCohort(decisionId, t.metaManager.GetConfig().ExpCohortIds)
	initAttributeCohort(hitCohortIds, attributes)
}

func (t *AbClient) initFeatureCohort(featureId string, decisionId string, attributes map[string]interface{}) {
	if t.cohortHandler == nil {
		return
	}

	// collect cohort ids of the specified feature
	feature, err := t.metaManager.GetConfig().GetFeatureFromId(featureId)
	if err != nil {
		return
	}
	if len(feature.CohortIds) == 0 {
		return
	}

	// call cdp, and init cohort into attributes
	hitCohortIds := t.cohortHandler.UseCohort(decisionId, feature.CohortIds)
	initAttributeCohort(hitCohortIds, attributes)
}

func (t *AbClient) initAllFeatureCohort(decisionId string, attributes map[string]interface{}) {
	if t.cohortHandler == nil || len(t.metaManager.GetConfig().FeatureCohortIds) == 0 {
		return
	}

	// call cdp, and init cohort into attributes
	hitCohortIds := t.cohortHandler.UseCohort(decisionId, t.metaManager.GetConfig().FeatureCohortIds)
	initAttributeCohort(hitCohortIds, attributes)
}
