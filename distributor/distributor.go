/**
 * Apache 2.0
 * Copyright 2022 Beijing Volcano Engine Technology Co., Ltd.
 */

package distributor

import (
	"fmt"
	"strings"

	"github.com/volcengine/datatester-go-sdk/consts"
	"github.com/volcengine/datatester-go-sdk/distributor/bucketer"
	e "github.com/volcengine/datatester-go-sdk/entities"
	"github.com/volcengine/datatester-go-sdk/meta/config"
	"github.com/volcengine/datatester-go-sdk/release"
)

type Distributor interface {
	GetExperimentVariant(c *config.ProductConfig, experimentId, decisionId string,
		attributes map[string]interface{}, experiment2variant map[string]string) (*e.Variant, error)
	GetFeatureVariant(c *config.ProductConfig, featureId, decisionId string,
		attributes map[string]interface{}) (*e.Variant, error)
	GetAllFeatureVariants(c *config.ProductConfig, decisionId string,
		attributes map[string]interface{}) ([]*e.Variant, []string, error)
}

type VariantsDistributor struct {
	bucketService *bucketer.Mmh3BucketService
	emptyVariant  *e.Variant
}

func NewVariantsDistributor() *VariantsDistributor {
	return &VariantsDistributor{
		bucketService: bucketer.NewMmh3BucketService(),
		emptyVariant:  &e.Variant{},
	}
}

func (v *VariantsDistributor) GetExperimentVariant(c *config.ProductConfig, experimentId, decisionId string,
	attributes map[string]interface{}, experiment2variant map[string]string) (*e.Variant, error) {
	experiment, err := c.GetExperimentFromId(experimentId)
	if err != nil || len(experiment.Id) == 0 {
		return nil, err
	}
	variant, err := v.tabExperimentVariant(c, experiment, decisionId, attributes, experiment2variant, true)
	if err != nil || len(variant.Id) == 0 {
		return nil, err
	}
	return variant, nil
}

func (v *VariantsDistributor) GetFeatureVariant(c *config.ProductConfig, featureId, decisionId string,
	attributes map[string]interface{}) (*e.Variant, error) {
	feature, err := c.GetFeatureFromId(featureId)
	if err != nil || len(feature.Id) == 0 || feature.Status != consts.FeatureStatusEnable {
		return nil, err
	}
	variant, err := v.getFeatVariantFromAllowList(c, featureId, decisionId)
	if err == nil && len(variant.Id) != 0 {
		return &variant, nil
	}
	variant, err = v.tabFeatureVariant(feature, decisionId, attributes)
	if err != nil || len(variant.Id) == 0 {
		return nil, err
	}
	return &variant, nil
}

func (v *VariantsDistributor) GetAllFeatureVariants(c *config.ProductConfig, decisionId string,
	attributes map[string]interface{}) ([]*e.Variant, []string, error) {
	variants := make([]*e.Variant, 0)
	featureIds := make([]string, 0)
	for id := range c.FeatureMap {
		variant, err := v.GetFeatureVariant(c, id, decisionId, attributes)
		if err == nil && variant != nil {
			variants = append(variants, variant)
			featureIds = append(featureIds, id)
		}
	}
	return variants, featureIds, nil
}

func (v *VariantsDistributor) handleAssociationRelations(c *config.ProductConfig, experiment *e.Experiment,
	decisionId string, attributes map[string]interface{}, experiment2variant map[string]string) {
	if len(experiment.AssociatedRelations) == 0 {
		return
	}
	for _, associateExperimentId := range experiment.AssociatedRelations {
		if attributes != nil {
			key := consts.ExperimentPrefix + associateExperimentId
			if _, ok := attributes[key]; ok {
				continue
			}
		}
		if attributes == nil {
			attributes = make(map[string]interface{})
		}
		associateExperiment, _ := c.GetExperimentFromId(associateExperimentId)
		variant, err := v.tabExperimentVariant(c, associateExperiment, decisionId, attributes, experiment2variant,
			false)
		if err != nil || len(variant.Id) == 0 {
			attributes[consts.ExperimentPrefix+associateExperimentId] = false
			continue
		}
		attributes[consts.ExperimentPrefix+associateExperimentId] = true
	}
}

func (v *VariantsDistributor) handleAllowList(experiment *e.Experiment, decisionId string,
	attributes map[string]interface{}) *e.Variant {
	whiteListMap := experiment.GetWhiteListMap()
	if len(whiteListMap) == 0 {
		return v.emptyVariant
	}
	variant, ok := whiteListMap[decisionId]
	if !ok {
		return v.emptyVariant
	}

	if experiment.IsUserGroupExperiment() {
		if v.EvaluateUserGroup(experiment, variant.Id, decisionId, attributes, experiment.FilterAllowList == consts.NeedFilterAllowList) {
			return variant
		}
		return v.emptyVariant
	}

	if experiment.FilterAllowList == consts.NeedFilterAllowList {
		if release.EvaluateFilters(experiment.Release.Filters, attributes) {
			return variant
		}
		return v.emptyVariant
	}
	return variant
}

func (v *VariantsDistributor) handleCacheAndReturnVariant(variant *e.Variant, experiment2variant map[string]string,
	needCache bool, experimentId string) (*e.Variant, error) {
	if needCache {
		experiment2variant[experimentId] = variant.Id
	}
	return variant, nil
}

func (v *VariantsDistributor) handleFreezeStatus(experiment *e.Experiment,
	experiment2variant map[string]string) (*e.Variant, error, bool) {
	if experiment.FreezeStatus != consts.ExperimentFreezeStatus &&
		experiment.VersionFreezeStatus != consts.ExperimentVersionFreezeStatus {
		return v.emptyVariant, nil, true
	}
	variantId, vidExist := experiment2variant[experiment.Id]
	if vidExist {
		variant, variantExist := experiment.VariantMap[variantId]
		if variantExist {
			return variant, nil, false
		}
	}
	if experiment.FreezeStatus == consts.ExperimentFreezeStatus {
		return v.emptyVariant, fmt.Errorf("experiment[%s] is freeze", experiment.Id), false
	}
	return v.emptyVariant, nil, true
}

func (v *VariantsDistributor) handleFatherExperiment(c *config.ProductConfig, experiment *e.Experiment,
	decisionId string, attributes map[string]interface{},
	experiment2variant map[string]string, variant *e.Variant, needCache bool) (*e.Variant, error) {
	if len(experiment.FatherExperimentId) == 0 || len(variant.FatherVariantIds) == 0 {
		return v.handleCacheAndReturnVariant(variant, experiment2variant, needCache, experiment.Id)
	}
	fatherExperiment, ok := c.ExperimentMap[experiment.FatherExperimentId]
	if !ok {
		return v.emptyVariant, fmt.Errorf("no father experiment[%s] exist in config", experiment.FatherExperimentId)
	}
	fatherVariant, err := v.tabExperimentVariant(c, fatherExperiment, decisionId,
		attributes, experiment2variant, false)
	if err != nil || len(fatherVariant.Id) == 0 {
		return v.emptyVariant, nil
	}
	for _, fatherVid := range variant.FatherVariantIds {
		if fatherVid == fatherVariant.Id {
			return v.handleCacheAndReturnVariant(variant, experiment2variant, needCache, experiment.Id)
		}
	}
	return v.emptyVariant, nil
}

func (v *VariantsDistributor) EvaluateUserGroup(experiment *e.Experiment, vid string, decisionId string,
	attributes map[string]interface{}, needFilter bool) bool {
	var hitUserGroupIds []string
	for _, r := range experiment.UserGroupReleases {
		if !needFilter || r.EvaluateRelease(attributes) {
			hitUserGroupIds = append(hitUserGroupIds, r.UserGroupId)
		}
	}

	// choose user group id
	if len(hitUserGroupIds) == 0 {
		return false
	}
	selectedUserGroupId := hitUserGroupIds[0]
	if len(hitUserGroupIds) > 1 {
		index, err := v.bucketService.GetTrafficBucketIndexWithMod(
			strings.Join([]string{decisionId, experiment.Name, "user-group"}, ":"), int32(len(hitUserGroupIds)))
		if err != nil {
			return false
		}
		selectedUserGroupId = hitUserGroupIds[index]
	}

	// write user group id to attributes
	if _, exist := attributes[consts.UserGroupRelation]; !exist {
		attributes[consts.UserGroupRelation] = make(map[string]string)
	}
	userGroupRelation, ok := attributes[consts.UserGroupRelation].(map[string]string)
	if !ok {
		return false
	}
	userGroupRelation[vid] = selectedUserGroupId

	return true
}

func (v *VariantsDistributor) tabExperimentVariant(c *config.ProductConfig, experiment *e.Experiment,
	decisionId string, attributes map[string]interface{},
	experiment2variant map[string]string, needCache bool) (*e.Variant, error) {
	if !experiment.IsCodingExperiment() && !experiment.IsCodingCampaign() {
		return v.emptyVariant, fmt.Errorf("experiment[%s] is not coding experiment", experiment.Id)
	}

	// handle association experiments
	v.handleAssociationRelations(c, experiment, decisionId, attributes, experiment2variant)

	// allow list
	allowListVariant := v.handleAllowList(experiment, decisionId, attributes)
	if len(allowListVariant.Id) != 0 {
		return v.handleCacheAndReturnVariant(allowListVariant, experiment2variant, needCache, experiment.Id)
	}

	// freeze experiment and traffic changes will not affect exposed users.
	cacheVariant, freezeErr, shouldContinue := v.handleFreezeStatus(experiment, experiment2variant)
	if !shouldContinue {
		// 客群实验，进组不出组仍需满足客群规则
		if cacheVariant.Id != "" && experiment.IsUserGroupExperiment() && !v.EvaluateUserGroup(experiment, cacheVariant.Id,
			decisionId, attributes, true) {
			return v.emptyVariant, nil
		}
		return cacheVariant, freezeErr
	}

	// validating experiments only handle allow list
	if experiment.Status != consts.RUNNING {
		return v.emptyVariant, fmt.Errorf("experiment[%s] status is [%d]", experiment.Id, experiment.Status)
	}

	// layer hash -> experiment
	layer, ok := c.LayerMap[experiment.LayerID]
	if !ok {
		return v.emptyVariant, fmt.Errorf("no layer[%s] exist in config", experiment.LayerID)
	}
	// check 互斥域
	hitDomain := layer.HitDomain(decisionId)
	if !hitDomain {
		return v.emptyVariant, nil
	}
	eid, err := v.tabLayerExperimentId(layer, decisionId)
	if err != nil || len(eid) == 0 || eid != experiment.Id {
		return v.emptyVariant, err
	}
	vid := v.tabExperimentVariantId(experiment, decisionId, attributes)
	if len(vid) == 0 {
		return v.emptyVariant, nil
	}
	variant := experiment.VariantMap[vid]

	// handle father experiments
	return v.handleFatherExperiment(c, experiment, decisionId, attributes, experiment2variant, variant, needCache)
}

func (v *VariantsDistributor) tabLayerExperimentId(layer *e.Layer, decisionId string) (string, error) {
	exBucketIndex, err := v.bucketService.GetTrafficBucketIndex(
		strings.Join([]string{decisionId, layer.Name}, ":"))
	if err != nil {
		return "", err
	}
	var eid string
	for _, trafficAllocation := range layer.TrafficAllocation {
		eid = trafficAllocation.EvaluateTraffic(exBucketIndex)
		if len(eid) != 0 {
			break
		}
	}
	return eid, nil
}

func (v *VariantsDistributor) tabExperimentVariantId(experiment *e.Experiment, decisionId string,
	attributes map[string]interface{}) string {
	vBucketIndex, err := v.bucketService.GetTrafficBucketIndex(
		strings.Join([]string{decisionId, experiment.Name}, ":"))
	if err != nil {
		return ""
	}

	vid := experiment.Release.EvaluateRelease(attributes, vBucketIndex)
	if vid == "" {
		return ""
	}

	if experiment.IsUserGroupExperiment() && !v.EvaluateUserGroup(experiment, vid, decisionId, attributes, true) {
		return ""
	}

	return vid
}

func (v *VariantsDistributor) getFeatVariantFromAllowList(c *config.ProductConfig, featureId,
	decisionId string) (e.Variant, error) {
	whiteListMap, err := c.GetFeatureAllowList(featureId)
	if err != nil || len(whiteListMap) == 0 {
		return e.Variant{}, err
	}
	variant, ok := whiteListMap[decisionId]
	if !ok {
		return e.Variant{}, nil
	}
	return variant, nil
}

func (v *VariantsDistributor) tabFeatureVariant(feature *e.Feature, decisionId string,
	attributes map[string]interface{}) (e.Variant, error) {
	if feature.Releases == nil || len(feature.Releases) == 0 {
		return e.Variant{}, nil
	}
	vBucketIndex, err := v.bucketService.GetTrafficBucketIndex(strings.Join([]string{decisionId, feature.Name}, ":"))
	if err != nil {
		return e.Variant{}, err
	}
	for _, r := range feature.Releases {
		vid := r.EvaluateRelease(attributes, vBucketIndex)
		if len(vid) != 0 {
			return feature.VariantMap[vid], nil
		}
	}
	return e.Variant{}, nil
}

var _ Distributor = &VariantsDistributor{}
