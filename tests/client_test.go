/**
 * Apache 2.0
 * Copyright 2022 Beijing Volcano Engine Technology Co., Ltd.
 */

package tests

import (
	"github.com/stretchr/testify/assert"
	"github.com/volcengine/datatester-go-sdk/client"
	"github.com/volcengine/datatester-go-sdk/config"
	abUserInfoHandler "github.com/volcengine/datatester-go-sdk/handler"
	"github.com/volcengine/datatester-go-sdk/tests/mock"
	"testing"
	"time"
)

func setUp() (*client.AbClient, map[string]interface{}) {
	return client.NewClient4Test(mock.MetaMock, abUserInfoHandler.NewDefaultUserAbInfoHandler()), make(map[string]interface{})
}

func TestGetExperimentVariantWithImpression(t *testing.T) {
	abClient, attributes := setUp()
	variant, _ := abClient.GetExperimentVariantWithImpression("77773", "decisionId",
		"trackId", attributes)
	assert.True(t, len(variant.Id) != 0)
	assert.True(t, mapContainsKey(variant.GetConfig(), "father"))
}

func TestVerifyFeatureEnabled(t *testing.T) {
	abClient, attributes := setUp()
	result, _ := abClient.VerifyFeatureEnabled("10202", "decisionId", attributes)
	assert.True(t, result)
	result, _ = abClient.VerifyFeatureEnabled("10203", "decisionId", attributes)
	assert.False(t, result)
}

func TestGetEnabledFeatureIds(t *testing.T) {
	abClient, attributes := setUp()
	result, _ := abClient.GetEnabledFeatureIds("decisionId", attributes)
	assert.True(t, result[0] == "10202")
}

func TestGetFeatureConfigsWithImpression(t *testing.T) {
	abClient, attributes := setUp()
	result, _ := abClient.GetFeatureConfigsWithImpression("10202", "decisionId",
		"trackId", attributes)
	assert.True(t, mapContainsKey(result, "feature"))
	result, _ = abClient.GetFeatureConfigsWithImpression("10202", "test_user",
		"trackId", attributes)
	assert.True(t, result["feature"]["vid"].(string) == "20101622")
	result, _ = abClient.GetFeatureConfigsWithImpression("10202", "test_user_02",
		"trackId", attributes)
	assert.True(t, result["feature"]["vid"].(string) == "20101623")
	result, _ = abClient.GetFeatureConfigsWithImpression("10202", "decisionId",
		"trackId", attributes)
	assert.True(t, result["feature"]["vid"].(string) == "20101623")
	attributes["age"] = 12347
	result, _ = abClient.GetFeatureConfigsWithImpression("10202", "decisionId",
		"trackId", attributes)
	assert.True(t, result["feature"]["vid"].(string) == "20101622")
	attributes["age"] = 12345
	result, _ = abClient.GetFeatureConfigsWithImpression("10202", "decisionId",
		"trackId", attributes)
	assert.True(t, result["feature"]["vid"].(string) == "20101622")
	attributes["age"] = 12344
	result, _ = abClient.GetFeatureConfigsWithImpression("10202", "decisionId",
		"trackId", attributes)
	assert.True(t, result["feature"]["vid"].(string) == "20101623")
}

func TestGetFeatureConfigs(t *testing.T) {
	abClient, attributes := setUp()
	result, _ := abClient.GetFeatureConfigs("10203", "decisionId", attributes)
	assert.True(t, len(result) == 0)
}

func TestGetExperimentConfigsWithImpression(t *testing.T) {
	abClient, attributes := setUp()
	result, _ := abClient.GetExperimentConfigsWithImpression("77775", "decisionId",
		"trackId", attributes)
	assert.True(t, mapContainsKey(result, "asso"))
}

func TestGetExperimentConfigs(t *testing.T) {
	abClient, attributes := setUp()
	result, _ := abClient.GetExperimentConfigs("77775", "decisionId", attributes)
	assert.True(t, mapContainsKey(result, "asso"))
}

func TestGetAllExperimentConfigs(t *testing.T) {
	abClient, attributes := setUp()
	result, _ := abClient.GetAllExperimentConfigs("decisionId", attributes)
	assert.True(t, mapContainsKey(result, "asso"))
	assert.True(t, mapContainsKey(result, "father"))
	assert.True(t, mapContainsKey(result, "child_01"))
	assert.True(t, mapContainsKey(result, "filter_param"))
}

func TestGetAllExperimentConfigsWithImpression(t *testing.T) {
	abClient, attributes := setUp()
	result, _ := abClient.GetAllExperimentConfigsWithImpression("decisionId", "trackId", attributes)
	assert.True(t, mapContainsKey(result, "asso"))
	assert.True(t, mapContainsKey(result, "father"))
	assert.True(t, mapContainsKey(result, "child_01"))
	assert.True(t, mapContainsKey(result, "filter_param"))
}

func TestGetExperimentVariantNameWithImpression(t *testing.T) {
	abClient, attributes := setUp()
	name, _ := abClient.GetExperimentVariantNameWithImpression("77775", "decisionId",
		"trackId", attributes)
	assert.True(t, name == "对照版本")
}

func TestGetExperimentVariantName(t *testing.T) {
	abClient, attributes := setUp()
	name, _ := abClient.GetExperimentVariantName("77773", "decisionId", attributes)
	assert.True(t, name == "实验版本1")
}

func TestActivate(t *testing.T) {
	abClient, attributes := setUp()
	result, _ := abClient.Activate(
		"asso", "decisionId", "trackId", nil, attributes)
	assert.True(t, result.(string) == "a" || result.(string) == "b")
	result, _ = abClient.Activate(
		"asso_invalid", "decisionId", "trackId", nil, attributes)
	assert.True(t, result == nil)
}

func TestActivateWithIdType(t *testing.T) {
	abClient, attributes := setUp()
	result, _ := abClient.ActivateWithIdType(
		"asso", "decisionId", "trackId", nil, attributes, "")
	assert.True(t, result.(string) == "a" || result.(string) == "b")
	result, _ = abClient.ActivateWithIdType(
		"asso_invalid", "decisionId", "trackId", nil, attributes, "")
	assert.True(t, result == nil)
}

func TestActivateWithVid(t *testing.T) {
	abClient, attributes := setUp()
	result, _ := abClient.ActivateWithVid("asso", "decisionId", "trackId", attributes)
	assert.True(t, result["val"].(string) == "a" || result["val"].(string) == "b")
	result, _ = abClient.ActivateWithVid("asso_invalid", "decisionId", "trackId", attributes)
	assert.True(t, result == nil)
}

func TestActivateWithoutImpression(t *testing.T) {
	abClient, attributes := setUp()
	result, _ := abClient.ActivateWithoutImpression("asso", "decisionId", attributes)
	assert.True(t, result["vid"] == "120314" || result["vid"] == "120315")
	result, _ = abClient.ActivateWithoutImpression("asso_invalid", "decisionId", attributes)
	assert.True(t, result == nil)
}

func TestGetAllFeatureConfigs(t *testing.T) {
	abClient, attributes := setUp()
	result, _ := abClient.GetAllFeatureConfigs("decisionId", attributes)
	assert.True(t, mapContainsKey(result, "feature"))
}

func mapContainsKey(targetMap map[string]map[string]interface{}, key string) bool {
	_, ok := targetMap[key]
	return ok
}

func TestNewClient(t *testing.T) {
	abClient := client.NewClient("3ac7ff2fad9d005f97857596a9078aff", config.WithMetaHost(config.MetaHostCN),
		config.WithTrackHost(config.TrackHostCN),
		config.WithFetchInterval(30*time.Second), config.WithWorkerNumOnce(5))
	assert.True(t, abClient != nil)
	val, _ := abClient.Activate(
		"go", "decisionId", "trackId", nil, map[string]interface{}{})
	if val != nil {
		assert.True(t, val.(string) == "v1" || val.(string) == "v2")
	}
}
